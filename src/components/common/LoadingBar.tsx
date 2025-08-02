import React, { useState, useEffect } from 'react'
import { Text, Box } from 'ink'

interface LoadingBarProps {
  label?: string
  color?: string
  width?: number
}

const LoadingBar: React.FC<LoadingBarProps> = ({ label = '', color = 'magenta', width = 20 }) => {
  const [position, setPosition] = useState(0)
  const [direction, setDirection] = useState(1)

  useEffect(() => {
    const timer = setInterval(() => {
      setPosition((prev) => {
        const next = prev + direction
        if (next >= width - 3 || next <= 0) {
          setDirection(-direction)
          return prev + direction
        }
        return next
      })
    }, 100)

    return () => clearInterval(timer)
  }, [direction, width])

  const bar = Array(width).fill('░')
  for (let i = 0; i < 3; i++) {
    if (position + i < width) {
      bar[position + i] = '█'
    }
  }

  return (
    <Box flexDirection="column">
      {label && <Text color={color}>{label}</Text>}
      <Text color={color}>[{bar.join('')}]</Text>
    </Box>
  )
}

export default LoadingBar
