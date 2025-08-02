// Example usage of the AudioPlayerService
// This file demonstrates how to use the audio player in the DictCLI application

import { audioPlayer } from './player.js'

async function exampleUsage() {
  try {
    // Basic playback
    await audioPlayer.play('/path/to/audio.mp3')
    console.log('Audio playback completed')

    // Playback with speed control (0.8x speed for slower learning)
    await audioPlayer.play('/path/to/audio.mp3', { speed: 0.8 })
    console.log('Slow playback completed')

    // Playback with speed control (1.2x speed for faster playback)
    await audioPlayer.play('/path/to/audio.mp3', { speed: 1.2 })
    console.log('Fast playback completed')

    // Playback with volume control
    await audioPlayer.play('/path/to/audio.mp3', { volume: 0.5 })
    console.log('Quiet playback completed')

    // Playback with both speed and volume control
    await audioPlayer.play('/path/to/audio.mp3', { speed: 0.9, volume: 0.8 })
    console.log('Custom playback completed')

    // Check player status
    console.log('Player status:', audioPlayer.getStatus())
    console.log('Is playing:', audioPlayer.isPlaying())

    // Stop playback if needed
    await audioPlayer.stop()
    console.log('Playback stopped')
  } catch (error) {
    console.error('Audio playback error:', error)
    if (error instanceof Error) {
      const err = error as Error & { code?: string; details?: string }
      if (err.code === 'FILE_NOT_FOUND') {
        console.error('Audio file not found or not accessible')
      } else if (err.code === 'PLAYBACK_FAILED') {
        console.error('Failed to play audio:', err.details)
      }
    }
  }
}

// Integration with TTSService example
import { TTSService } from '../openai/tts.js'

async function playGeneratedAudio() {
  const tts = new TTSService()

  try {
    // Generate audio from text
    const audioPath = await tts.generateSpeech(
      'Hello, this is a test sentence for dictation practice.',
      'ALEX',
      1.0,
    )

    // Play the generated audio with speed control for learning
    await audioPlayer.play(audioPath, { speed: 0.8 })
    console.log('Generated audio playback completed')
  } catch (error) {
    console.error('Error in TTS + Audio pipeline:', error)
  }
}

export { exampleUsage, playGeneratedAudio }
