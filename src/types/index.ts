export type VoiceDisplayName = 'ALEX' | 'SARA' | 'EVAN' | 'NOVA' | 'NICK' | 'FAYE'

export type OpenAIVoice = 'alloy' | 'shimmer' | 'echo' | 'nova' | 'onyx' | 'fable'

export const VOICE_MAPPING: Record<VoiceDisplayName, OpenAIVoice> = {
  ALEX: 'echo', // male voice
  SARA: 'shimmer', // female voice
  EVAN: 'onyx', // deep male voice
  NOVA: 'nova', // female voice
  NICK: 'fable', // British male voice
  FAYE: 'alloy', // neutral/female voice
}

export type Level = 'CEFR_A1' | 'CEFR_A2' | 'CEFR_B1' | 'CEFR_B2' | 'CEFR_C1' | 'CEFR_C2'

export type Topic =
  | 'EverydayLife'
  | 'Travel'
  | 'Technology'
  | 'Health'
  | 'Entertainment'
  | 'Business'
  | 'Random'

export interface Settings {
  voice: VoiceDisplayName
  level: Level
  topic: Topic
  wordCount: number
  speed: number
}

export interface SlashCommand {
  command: '/replay' | '/settings' | '/quit' | '/giveup'
  description: string
  action: () => void
}

export interface Round {
  id: string
  sentence: string
  userInput: string
  score: number
  wer: number
  errors: Array<{
    expected: string
    actual: string
    explanation: string
  }>
  alternatives: string[]
  timestamp: Date
}

// Audio Player Types
export interface AudioPlayerOptions {
  volume?: number
  speed?: number
}

export interface AudioPlayerError extends Error {
  code: string
  details?: string
}

export type AudioPlayerStatus = 'idle' | 'playing' | 'paused' | 'stopped' | 'error'

export interface AudioPlayer {
  play(filePath: string, options?: AudioPlayerOptions): Promise<void>
  stop(): Promise<void>
  isPlaying(): boolean
  getStatus(): AudioPlayerStatus
}
