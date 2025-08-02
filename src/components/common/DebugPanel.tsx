import React from 'react'
import { Box, Text } from 'ink'
import { useStore } from '../../store/useStore.js'

const DebugPanel: React.FC = () => {
  const debugState = useStore((state) => state.debugState)

  if (!debugState.isDebugMode) {
    return null
  }

  return (
    <Box
      flexDirection="column"
      borderStyle="single"
      borderColor="yellow"
      paddingX={1}
      marginBottom={1}
    >
      <Text color="yellow" bold>
        DEBUG MODE
      </Text>
      <Box flexDirection="column" marginTop={1}>
        {debugState.logs.slice(-5).map((log, index) => (
          <Box key={index} gap={1}>
            <Text color="gray">{log.timestamp.toISOString().substr(11, 8)}</Text>
            <Text
              color={
                log.type === 'error'
                  ? 'red'
                  : log.type === 'api'
                    ? 'cyan'
                    : log.type === 'audio'
                      ? 'magenta'
                      : 'green'
              }
            >
              [{log.type.toUpperCase()}]
            </Text>
            <Text>{log.message}</Text>
            {log.details !== undefined && log.details !== null && (
              <Text color="gray" dimColor>
                {' '}
                {(() => {
                  try {
                    const detailStr =
                      typeof log.details === 'string' ? log.details : JSON.stringify(log.details)
                    return detailStr.length > 50 ? detailStr.slice(0, 50) + '...' : detailStr
                  } catch {
                    return '[Invalid details]'
                  }
                })()}
              </Text>
            )}
          </Box>
        ))}
        {debugState.logs.length === 0 && <Text color="gray">No debug logs yet...</Text>}
      </Box>
    </Box>
  )
}

export default DebugPanel
