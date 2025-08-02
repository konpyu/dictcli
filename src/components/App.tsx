import React from 'react'
import { Box, useInput, useApp, type Key } from 'ink'
import { Header } from './common/Header.js'
import { useStore } from '../store/useStore.js'
import LearningView from './LearningView.js'
import ResultView from './ResultView.js'
import SettingsModal from './SettingsModal.js'
import DebugPanel from './common/DebugPanel.js'
import ErrorBoundary from './common/ErrorBoundary.js'
import type { Settings } from '../types/index.js'

interface AppProps {
  initialSettings: Settings
}

export const App: React.FC<AppProps> = ({ initialSettings }) => {
  const { exit } = useApp()
  const viewState = useStore((state) => state.viewState)
  const setViewState = useStore((state) => state.setViewState)
  const settings = useStore((state) => state.settings)
  const updateSettings = useStore((state) => state.updateSettings)
  const updateAudioState = useStore((state) => state.updateAudioState)
  const [isInitialized, setIsInitialized] = React.useState(false)

  // Initialize store settings on first render
  React.useEffect(() => {
    if (process.env.DICTCLI_DEBUG === 'true') {
      console.log('App: Updating settings with initialSettings:', initialSettings)
    }
    updateSettings(initialSettings)
    updateAudioState({ speed: initialSettings.speed })
    setIsInitialized(true)
  }, [initialSettings, updateSettings])

  // Handle keyboard input (only when modal is not shown)
  useInput((input: string, key: Key) => {
    // Don't handle input when settings modal is shown
    if (viewState === 'settings') return

    // Don't handle input when in learning view (handled by LearningView component)
    if (viewState === 'learning') return

    // Exit on 'q' or Ctrl+C (only in result view)
    if (input === 'q' || (key.ctrl && input === 'c')) {
      exit()
      return
    }

    // For result view, handle single key shortcuts
    if (viewState === 'result') {
      return // Handled in ResultView component
    }
  })

  return (
    <ErrorBoundary>
      <Box flexDirection="column" width="100%" height="100%">
        <DebugPanel />
        {viewState !== 'settings' && <Header settings={settings} />}

        {viewState === 'learning' && isInitialized && (
          <Box flexDirection="column" paddingX={1} paddingY={1}>
            <Box borderStyle="single" borderColor="blue" paddingX={2} paddingY={1} marginTop={1}>
              <LearningView />
            </Box>
          </Box>
        )}

        {viewState === 'result' && (
          <Box flexDirection="column" paddingX={1} paddingY={1}>
            <Box borderStyle="single" borderColor="green" paddingX={2} paddingY={1} marginTop={1}>
              <ResultView
                onNext={() => {
                  // Simply transition to learning view
                  // The LearningView will handle starting a new round
                  setViewState('learning')
                }}
                onReplay={() => console.log('Replay audio')}
                onSettings={() => setViewState('settings')}
                onQuit={() => exit()}
              />
            </Box>
          </Box>
        )}

        {viewState === 'settings' && (
          <Box flexDirection="column" alignItems="center" justifyContent="center" height="100%">
            <SettingsModal
              onClose={() => setViewState('learning')}
              onSave={() => {
                setViewState('learning')
              }}
            />
          </Box>
        )}
      </Box>
    </ErrorBoundary>
  )
}
