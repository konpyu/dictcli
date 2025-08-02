import React from 'react'
import { Box, Text } from 'ink'
import type { Settings } from '../../types/index.js'

interface HeaderProps {
  settings: Settings
}

export const Header: React.FC<HeaderProps> = ({ settings }) => {
  const { voice, level, topic, wordCount, speed } = settings

  return (
    <Box borderStyle="single" borderColor="blue" paddingX={1}>
      <Text color="blue" bold>
        DictCLI
      </Text>
      <Text> ─ </Text>
      <Text>
        Voice:{voice} | Level:{level} | Topic:{topic} | Length:{wordCount}w | Speed:
        {speed.toFixed(1)}×
      </Text>
    </Box>
  )
}
