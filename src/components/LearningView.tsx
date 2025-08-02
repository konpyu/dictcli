import React, { useState, useEffect } from 'react'
import { Box, Text } from 'ink'
import TextInput from './common/TextInput.js'
import { useStore } from '../store/useStore.js'
import SlashCommandMenu from './SlashCommandMenu.js'
import { orchestrator } from '../services/orchestrator.js'
import Spinner from './common/Spinner.js'
import PulsingText from './common/PulsingText.js'

const LearningView: React.FC = () => {
  const [input, setInput] = useState('')
  const [showCommandMenu, setShowCommandMenu] = useState(false)
  const [selectedCommandIndex, setSelectedCommandIndex] = useState(0)

  const audioState = useStore((state) => state.audioState)
  const roundState = useStore((state) => state.roundState)
  const settings = useStore((state) => state.settings)
  const setViewState = useStore((state) => state.setViewState)
  const isScoring = useStore((state) => state.roundState.isScoring)

  // Slash commands
  const slashCommands = [
    {
      command: '/replay',
      description: 'Replay the audio',
      action: async () => {
        setShowCommandMenu(false)
        setInput('')
        try {
          await orchestrator.replayAudio()
        } catch (error) {
          console.error('Failed to replay audio:', error)
        }
      },
    },
    {
      command: '/settings',
      description: 'Open settings panel',
      action: () => {
        setShowCommandMenu(false)
        setInput('')
        setViewState('settings')
      },
    },
    {
      command: '/giveup',
      description: 'Show answer with gaps',
      action: () => {
        setShowCommandMenu(false)
        const gapFill = orchestrator.showGapFill()
        setInput(gapFill)
      },
    },
    {
      command: '/quit',
      description: 'Exit the application',
      action: () => {
        process.exit(0)
      },
    },
  ]

  // Filter commands based on input
  const filteredCommands = slashCommands.filter((cmd) => cmd.command.startsWith(input))

  // Start new round when component mounts
  useEffect(() => {
    if (process.env.DICTCLI_DEBUG === 'true') {
      console.log('LearningView: Starting new round with settings:', settings)
    }
    // Always start a new round when LearningView is mounted
    orchestrator.startNewRound(settings).catch((error) => {
      console.error('Failed to start new round:', error)
    })
  }, [settings]) // Re-run when settings change

  const handleInputChange = (value: string) => {
    setInput(value)

    // Show command menu when "/" is typed
    if (value === '/') {
      setShowCommandMenu(true)
      setSelectedCommandIndex(0)
    } else if (value === '') {
      setShowCommandMenu(false)
    } else if (value.startsWith('/')) {
      // Keep menu open while typing command
      setShowCommandMenu(true)
      // Reset selection when filtered list changes
      setSelectedCommandIndex(0)
    } else {
      setShowCommandMenu(false)
    }
  }

  const handleSubmit = async () => {
    // Prevent submit during scoring
    if (isScoring) return

    if (showCommandMenu && filteredCommands.length > 0) {
      // Execute selected command
      if (filteredCommands[selectedCommandIndex]) {
        await filteredCommands[selectedCommandIndex].action()
      }
    } else if (!input.startsWith('/') && input.trim() !== '') {
      // Submit answer for scoring
      try {
        await orchestrator.scoreAnswer(input)
        setInput('')
      } catch (error) {
        console.error('Failed to score answer:', error)
      }
    }
  }

  const handleKeyPress = (
    _input: string,
    key: { upArrow?: boolean; downArrow?: boolean; escape?: boolean },
  ) => {
    if (showCommandMenu) {
      if (key.upArrow) {
        setSelectedCommandIndex(Math.max(0, selectedCommandIndex - 1))
      } else if (key.downArrow) {
        setSelectedCommandIndex(Math.min(filteredCommands.length - 1, selectedCommandIndex + 1))
      } else if (key.escape) {
        setShowCommandMenu(false)
        setInput('')
      }
    }
  }

  const isPlaying = audioState.status === 'playing'
  const playbackSpeed = audioState.speed || 1.0

  return (
    <Box flexDirection="column" width="100%">
      <Box marginBottom={1}>
        {roundState.isGenerating ? (
          <Box>
            <Spinner color="cyan" />
            <Text> </Text>
            <PulsingText text="Generating problem..." baseColor="cyan" />
          </Box>
        ) : isScoring ? (
          <Box>
            <Spinner color="magenta" />
            <Text> </Text>
            <PulsingText text="Scoring your answer..." baseColor="magenta" />
          </Box>
        ) : isPlaying ? (
          <Text color="yellow">(🔊 Playing… ⏩{playbackSpeed}×)</Text>
        ) : roundState.currentRound ? (
          <Text color="green">(Ready for input - Type your answer or use /commands)</Text>
        ) : (
          <Text color="gray">(Waiting for round...)</Text>
        )}
      </Box>

      <Box flexDirection="row" alignItems="center">
        <Text> {'>'} </Text>
        <Box marginLeft={1}>
          <TextInput
            value={input}
            onChange={handleInputChange}
            onSubmit={handleSubmit}
            onKeypress={handleKeyPress}
            isDisabled={isScoring}
          />
        </Box>
      </Box>

      {showCommandMenu && (
        <SlashCommandMenu commands={filteredCommands} selectedIndex={selectedCommandIndex} />
      )}
    </Box>
  )
}

export default LearningView
