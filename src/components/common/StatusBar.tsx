import React from 'react'
import { Box, Text } from 'ink'

interface StatusBarProps {
  status?: string
  playbackSpeed?: number
  isPlaying?: boolean
}

export const StatusBar: React.FC<StatusBarProps> = ({
  status = '',
  playbackSpeed = 1.0,
  isPlaying = false,
}) => {
  return (
    <Box paddingX={1} paddingY={0}>
      {isPlaying && <Text>(🔊 Playing… ⏩{playbackSpeed}×)</Text>}
      {status && !isPlaying && <Text>{status}</Text>}
    </Box>
  )
}
