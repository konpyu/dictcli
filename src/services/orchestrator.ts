import type { Settings, Round } from '../types/index.js'
import { ProblemGenerator } from './openai/generator.js'
import { TTSService } from './openai/tts.js'
import { ScorerService } from './openai/scorer.js'
import { AudioPlayerService } from './audio/player.js'
import { store } from '../store/useStore.js'
import { saveHistory } from '../storage/history.js'
import { v4 as uuidv4 } from 'uuid'

export class Orchestrator {
  private generator: ProblemGenerator
  private tts: TTSService
  private scorer: ScorerService
  private audioPlayer: AudioPlayerService
  private isRunning = false

  constructor() {
    this.generator = new ProblemGenerator()
    this.tts = new TTSService()
    this.scorer = new ScorerService()
    this.audioPlayer = new AudioPlayerService()
  }

  async startNewRound(settings: Settings): Promise<void> {
    if (this.isRunning) return
    this.isRunning = true

    try {
      const { addDebugLog, setIsGenerating, updateAudioState } = store.getState()

      // Check if we have a pre-generated round
      const preGeneratedRound = store.getState().roundState.preGeneratedRound
      let round: Round

      if (preGeneratedRound) {
        addDebugLog('state', 'Using pre-generated round', {
          sentence: preGeneratedRound.sentence,
          preGeneratedAt: preGeneratedRound.timestamp,
        })
        round = preGeneratedRound
        store.getState().setPreGeneratedRound(null)
      } else {
        // Generate new problem
        addDebugLog('api', 'Generating new problem...', {
          level: settings.level,
          topic: settings.topic,
          wordCount: settings.wordCount,
          voice: settings.voice,
        })
        setIsGenerating(true)

        const problem = await this.generator.generateProblem(
          settings.level,
          settings.topic,
          settings.wordCount,
        )

        addDebugLog('api', 'Problem generated', { problem })

        // Generate audio
        addDebugLog('api', 'Generating audio...', { voice: settings.voice, speed: settings.speed })
        const audioPath = await this.tts.generateSpeech(problem, settings.voice, settings.speed)
        addDebugLog('api', 'Audio generated', { audioPath })

        round = {
          id: uuidv4(),
          sentence: problem,
          userInput: '',
          score: 0,
          wer: 0,
          errors: [],
          alternatives: [],
          timestamp: new Date(),
        }
        setIsGenerating(false)
      }

      // Set current round
      store.getState().setCurrentRound(round)

      // Play audio
      const audioSpeed = store.getState().audioState.speed
      addDebugLog('audio', 'Playing audio', { speed: audioSpeed })
      updateAudioState({ status: 'playing' })

      // Get audio path from cache
      const audioPath = await this.tts.getCachedPath(round.sentence, settings.voice, audioSpeed)
      if (audioPath) {
        await this.audioPlayer.play(audioPath, { speed: audioSpeed })
        updateAudioState({ status: 'idle' })
      } else {
        throw new Error('Audio file not found in cache')
      }

      // Start pre-generating next round
      this.preGenerateNextRound(settings)
    } catch (error) {
      store.getState().addDebugLog('error', 'Failed to start new round', { error })
      store.getState().setIsGenerating(false)
      store.getState().updateAudioState({ status: 'error' })
      throw error
    } finally {
      this.isRunning = false
    }
  }

  async replayAudio(): Promise<void> {
    const { currentRound } = store.getState().roundState
    const { settings, audioState } = store.getState()

    if (!currentRound) return

    try {
      store.getState().addDebugLog('audio', 'Replaying audio')
      store.getState().updateAudioState({ status: 'playing' })

      const audioPath = await this.tts.getCachedPath(
        currentRound.sentence,
        settings.voice,
        audioState.speed,
      )
      if (audioPath) {
        await this.audioPlayer.play(audioPath, { speed: audioState.speed })
        store.getState().updateAudioState({ status: 'idle' })
      }
    } catch (error) {
      store.getState().addDebugLog('error', 'Failed to replay audio', { error })
      store.getState().updateAudioState({ status: 'error' })
    }
  }

  async scoreAnswer(userInput: string): Promise<void> {
    const { currentRound, isScoring } = store.getState().roundState
    if (!currentRound) return

    // Prevent multiple scoring requests
    if (isScoring) {
      store.getState().addDebugLog('state', 'Scoring already in progress, ignoring request')
      return
    }

    try {
      store.getState().setIsScoring(true)
      store.getState().addDebugLog('api', 'Scoring answer...', { userInput })

      const result = await this.scorer.scoreAnswer(currentRound.sentence, userInput)
      store.getState().addDebugLog('api', 'Scoring complete', { result })

      // Update round with results
      const scoredRound: Round = {
        ...currentRound,
        userInput,
        score: result.score,
        wer: result.wer,
        errors: result.errors,
        alternatives: result.alternatives,
      }

      // Update store
      store.getState().setCurrentRound(scoredRound)
      store.getState().addRoundToHistory(scoredRound)

      // Save to history
      await saveHistory(scoredRound)

      // Change view to result
      store.getState().setViewState('result')
    } catch (error) {
      store.getState().addDebugLog('error', 'Failed to score answer', { error })
      throw error
    } finally {
      store.getState().setIsScoring(false)
    }
  }

  async stopAudio(): Promise<void> {
    try {
      await this.audioPlayer.stop()
      store.getState().updateAudioState({ status: 'stopped' })
    } catch (error) {
      store.getState().addDebugLog('error', 'Failed to stop audio', { error })
    }
  }

  private async preGenerateNextRound(settings: Settings): Promise<void> {
    try {
      store.getState().addDebugLog('state', 'Pre-generating next round...', {
        level: settings.level,
        topic: settings.topic,
        wordCount: settings.wordCount,
      })

      // Generate problem in background
      const problem = await this.generator.generateProblem(
        settings.level,
        settings.topic,
        settings.wordCount,
      )

      // Generate and cache audio
      await this.tts.generateSpeech(problem, settings.voice, settings.speed)

      const preGeneratedRound: Round = {
        id: uuidv4(),
        sentence: problem,
        userInput: '',
        score: 0,
        wer: 0,
        errors: [],
        alternatives: [],
        timestamp: new Date(),
      }

      store.getState().setPreGeneratedRound(preGeneratedRound)
      store.getState().addDebugLog('state', 'Pre-generation complete', {
        sentence: problem,
        level: settings.level,
        topic: settings.topic,
        wordCount: settings.wordCount,
      })
    } catch (error) {
      store.getState().addDebugLog('error', 'Pre-generation failed', { error })
      // Don't throw - pre-generation is optional
    }
  }

  showGapFill(): string {
    const { currentRound } = store.getState().roundState
    if (!currentRound) return ''

    // Simple gap-fill: show first and last word, hide middle words
    const words = currentRound.sentence.split(' ')
    if (words.length <= 2) return currentRound.sentence

    return words
      .map((word: string, index: number) => {
        if (index === 0 || index === words.length - 1) {
          return word
        }
        return '_'.repeat(word.length)
      })
      .join(' ')
  }
}

// Singleton instance
export const orchestrator = new Orchestrator()
