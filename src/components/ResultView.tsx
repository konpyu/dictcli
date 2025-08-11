import React from 'react'
import { Box, Text } from 'ink'
import { useInput, type Key } from 'ink'
import { useStore } from '../store/useStore.js'
import { orchestrator } from '../services/orchestrator.js'

interface ResultViewProps {
  onNext?: () => void
  onReplay?: () => void
  onSettings?: () => void
  onQuit?: () => void
}

const ResultView: React.FC<ResultViewProps> = ({
  onNext = () => console.log('Next round'),
  onReplay: _onReplay = () => console.log('Replay audio'),
  onSettings = () => console.log('Open settings'),
  onQuit = () => process.exit(0),
}) => {
  const currentRound = useStore((state) => state.roundState.currentRound)

  // Use actual data from current round, with fallback
  // Return nothing if no current round to prevent display overlap during transitions
  if (!currentRound) {
    return null
  }

  const result = currentRound

  // Handle keyboard shortcuts
  useInput(async (input: string, key: Key) => {
    if (key.return || input === 'n' || input === 'N') {
      // Just transition to learning view, don't start new round here
      onNext()
    } else if (input === 'r' || input === 'R') {
      // Replay audio
      try {
        await orchestrator.replayAudio()
      } catch (error) {
        console.error('Failed to replay audio:', error)
      }
    } else if (input === 's' || input === 'S') {
      onSettings()
    } else if (input === 'q' || input === 'Q') {
      onQuit()
    }
  })

  return (
    <Box flexDirection="column" width="100%">
      {/* Score Header */}
      <Box marginBottom={1}>
        <Text color="green" bold>
          Score: {result.score}% WER {result.wer.toFixed(2)}
        </Text>
      </Box>

      {/* Divider */}
      <Box marginBottom={1}>
        <Text>{'─'.repeat(70)}</Text>
      </Box>

      {/* Errors */}
      {result.errors.length > 0 ? (
        result.errors.map((error, index) => (
          <Box key={index} marginBottom={1}>
            <Text>
              Error: <Text color="red">{error.actual}</Text> →{' '}
              <Text color="green">{error.expected}</Text>{' '}
              <Text color="gray">({error.explanation})</Text>
            </Text>
          </Box>
        ))
      ) : result.score === 100 ? (
        <Box marginBottom={1}>
          <Text color="green">Perfect! All correct.</Text>
        </Box>
      ) : null}

      {/* Correct Sentence */}
      <Box marginBottom={1}>
        <Text>
          Correct answer:{' '}
          <Text color="green" bold>
            {result.sentence}
          </Text>
        </Text>
      </Box>

      {/* User Input */}
      {result.userInput && (
        <Box marginBottom={1}>
          <Text>
            Your answer:{' '}
            <Text color={result.score === 100 ? 'green' : 'yellow'}>{result.userInput}</Text>
          </Text>
        </Box>
      )}

      {/* Alternatives */}
      {result.alternatives.length > 0 && (
        <>
          {result.alternatives.map((alt, index) => (
            <Box key={index}>
              <Text>
                Alternative {index + 1}: <Text color="cyan">{alt}</Text>
              </Text>
            </Box>
          ))}
        </>
      )}

      {/* Divider */}
      <Box marginTop={1} marginBottom={1}>
        <Text>{'─'.repeat(70)}</Text>
      </Box>

      {/* Instructions */}
      <Box>
        <Text color="gray">
          Press{' '}
          <Text color="white" bold>
            Enter
          </Text>{' '}
          for next round (or <Text color="white">N</Text>/R/S/Q for other actions)
        </Text>
      </Box>
    </Box>
  )
}

export default ResultView
