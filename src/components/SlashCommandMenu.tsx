import React from 'react'
import { Box, Text } from 'ink'

interface Command {
  command: string
  description: string
  action: () => void
}

interface SlashCommandMenuProps {
  commands: Command[]
  selectedIndex: number
}

const SlashCommandMenu: React.FC<SlashCommandMenuProps> = ({ commands, selectedIndex }) => {
  if (commands.length === 0) {
    return null
  }

  return (
    <Box
      flexDirection="column"
      marginTop={1}
      marginLeft={3}
      borderStyle="single"
      borderColor="blue"
      paddingX={1}
    >
      {commands.map((cmd, index) => (
        <Box key={cmd.command} paddingY={0}>
          <Text
            color={index === selectedIndex ? 'blue' : 'white'}
            backgroundColor={index === selectedIndex ? 'white' : undefined}
          >
            <Text bold>{cmd.command.padEnd(12)}</Text>
            <Text> {cmd.description}</Text>
          </Text>
        </Box>
      ))}
    </Box>
  )
}

export default SlashCommandMenu
