import React from 'react'
import { create } from 'zustand'
import { subscribeWithSelector } from 'zustand/middleware'
import type { Settings, Round, AudioPlayerStatus } from '../types/index.js'

type ViewState = 'learning' | 'result' | 'settings'

interface AudioState {
  status: AudioPlayerStatus
  currentFile: string | null
  speed: number
}

interface RoundState {
  currentRound: Round | null
  roundHistory: Round[]
  isGenerating: boolean
  isScoring: boolean
  preGeneratedRound: Round | null
}

interface DebugLog {
  timestamp: Date
  type: 'api' | 'audio' | 'state' | 'error'
  message: string
  details?: unknown
}

interface DebugState {
  logs: DebugLog[]
  isDebugMode: boolean
}

interface AppState {
  // View State
  viewState: ViewState
  setViewState: (viewState: ViewState) => void

  // Settings
  settings: Settings
  updateSettings: (settings: Partial<Settings>) => void

  // Audio State
  audioState: AudioState
  updateAudioState: (state: Partial<AudioState>) => void

  // Round State
  roundState: RoundState
  setCurrentRound: (round: Round) => void
  addRoundToHistory: (round: Round) => void
  setIsGenerating: (isGenerating: boolean) => void
  setIsScoring: (isScoring: boolean) => void
  setPreGeneratedRound: (round: Round | null) => void

  // User Input
  userInput: string
  setUserInput: (input: string) => void

  // Debug State
  debugState: DebugState
  addDebugLog: (
    type: 'api' | 'audio' | 'state' | 'error',
    message: string,
    details?: unknown,
  ) => void
  clearDebugLogs: () => void

  // Command State
  showSlashCommandMenu: boolean
  setShowSlashCommandMenu: (show: boolean) => void
}

const DEFAULT_SETTINGS: Settings = {
  voice: 'ALEX',
  level: 'CEFR_A1',
  topic: 'Business',
  wordCount: 10,
  speed: 1.0,
}

const useStoreImpl = create<AppState>()(
  subscribeWithSelector((set, get) => ({
    // View State
    viewState: 'learning',
    setViewState: (viewState) => {
      set({ viewState })
      get().addDebugLog('state', `View changed to: ${viewState}`)
    },

    // Settings
    settings: DEFAULT_SETTINGS,
    updateSettings: (newSettings) =>
      set((state) => ({
        settings: { ...state.settings, ...newSettings },
      })),

    // Audio State
    audioState: {
      status: 'idle',
      currentFile: null,
      speed: 1.0,
    },
    updateAudioState: (newState) =>
      set((state) => ({
        audioState: { ...state.audioState, ...newState },
      })),

    // Round State
    roundState: {
      currentRound: null,
      roundHistory: [],
      isGenerating: false,
      isScoring: false,
      preGeneratedRound: null,
    },
    setCurrentRound: (round) =>
      set((state) => ({
        roundState: { ...state.roundState, currentRound: round },
      })),
    addRoundToHistory: (round) =>
      set((state) => ({
        roundState: {
          ...state.roundState,
          roundHistory: [...state.roundState.roundHistory, round],
        },
      })),
    setIsGenerating: (isGenerating) =>
      set((state) => ({
        roundState: { ...state.roundState, isGenerating },
      })),
    setIsScoring: (isScoring) =>
      set((state) => ({
        roundState: { ...state.roundState, isScoring },
      })),
    setPreGeneratedRound: (round) =>
      set((state) => ({
        roundState: { ...state.roundState, preGeneratedRound: round },
      })),

    // User Input
    userInput: '',
    setUserInput: (input) => set({ userInput: input }),

    // Debug State
    debugState: {
      logs: [],
      isDebugMode: process.env.DICTCLI_DEBUG === 'true',
    },
    addDebugLog: (type, message, details) =>
      set((state) => ({
        debugState: {
          ...state.debugState,
          logs: [
            ...state.debugState.logs.slice(-49), // Keep last 50 logs
            { timestamp: new Date(), type, message, details },
          ],
        },
      })),
    clearDebugLogs: () =>
      set((state) => ({
        debugState: { ...state.debugState, logs: [] },
      })),

    // Command State
    showSlashCommandMenu: false,
    setShowSlashCommandMenu: (show) => set({ showSlashCommandMenu: show }),
  })),
)

// Ink v4 workaround: Manual subscription pattern
export const useStore = <T>(selector: (state: AppState) => T): T => {
  const [, forceUpdate] = React.useReducer((x) => x + 1, 0)
  const stateRef = React.useRef<T>(selector(useStoreImpl.getState()))

  React.useEffect(() => {
    const unsubscribe = useStoreImpl.subscribe(selector, (newState) => {
      stateRef.current = newState
      forceUpdate()
    })
    return unsubscribe
  }, [selector])

  return stateRef.current
}

// Export the store instance for direct access if needed
export const store = useStoreImpl

// Export getState for non-hook access
export const getState = useStoreImpl.getState

// Export setState for updates outside of components
export const setState = useStoreImpl.setState
