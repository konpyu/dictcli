import React, { useState, useEffect } from 'react'
import { Text } from 'ink'

interface SpinnerProps {
  label?: string
  color?: string
}

const Spinner: React.FC<SpinnerProps> = ({ label = '', color = 'magenta' }) => {
  const frames = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏']
  const [frameIndex, setFrameIndex] = useState(0)

  useEffect(() => {
    const timer = setInterval(() => {
      setFrameIndex((prev) => (prev + 1) % frames.length)
    }, 80)

    return () => clearInterval(timer)
  }, [frames.length])

  return (
    <Text color={color}>
      {frames[frameIndex]} {label}
    </Text>
  )
}

export default Spinner
