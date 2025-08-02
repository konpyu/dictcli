import React, { useState, useEffect } from 'react'
import { Text } from 'ink'

interface PulsingTextProps {
  text: string
  baseColor?: string
}

const PulsingText: React.FC<PulsingTextProps> = ({ text, baseColor = 'magenta' }) => {
  const [brightness, setBrightness] = useState(0)
  const [increasing, setIncreasing] = useState(true)

  useEffect(() => {
    const timer = setInterval(() => {
      setBrightness((prev) => {
        if (increasing && prev >= 2) {
          setIncreasing(false)
          return prev - 1
        } else if (!increasing && prev <= 0) {
          setIncreasing(true)
          return prev + 1
        }
        return increasing ? prev + 1 : prev - 1
      })
    }, 300)

    return () => clearInterval(timer)
  }, [increasing])

  const colors = {
    magenta: ['magenta', 'magentaBright', 'white'],
    cyan: ['cyan', 'cyanBright', 'white'],
    yellow: ['yellow', 'yellowBright', 'white'],
  }

  const colorArray = colors[baseColor as keyof typeof colors] || colors.magenta
  const currentColor = colorArray[brightness]

  return <Text color={currentColor}>{text}</Text>
}

export default PulsingText
