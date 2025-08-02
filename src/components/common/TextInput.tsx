import React, { useState, useEffect } from 'react'
import { Text, useInput, type Key } from 'ink'

interface TextInputProps {
  value: string
  onChange: (value: string) => void
  onSubmit?: () => void
  onKeypress?: (
    input: string,
    key: {
      upArrow?: boolean
      downArrow?: boolean
      escape?: boolean
      return?: boolean
      backspace?: boolean
      delete?: boolean
      leftArrow?: boolean
      rightArrow?: boolean
      ctrl?: boolean
      meta?: boolean
    },
  ) => void
  focus?: boolean
  isDisabled?: boolean
}

const TextInput: React.FC<TextInputProps> = ({
  value = '',
  onChange,
  onSubmit = () => {},
  onKeypress = () => {},
  focus = true,
  isDisabled = false,
}) => {
  const [cursor, setCursor] = useState(value.length)

  useEffect(() => {
    setCursor(value.length)
  }, [value])

  useInput((input: string, key: Key) => {
    if (!focus || isDisabled) return

    // Pass all keypress events to parent
    onKeypress(input, key)

    if (key.return) {
      onSubmit()
    } else if (key.backspace || key.delete) {
      if (cursor > 0) {
        const newValue = value.slice(0, cursor - 1) + value.slice(cursor)
        onChange(newValue)
        setCursor(cursor - 1)
      }
    } else if (key.leftArrow) {
      setCursor(Math.max(0, cursor - 1))
    } else if (key.rightArrow) {
      setCursor(Math.min(value.length, cursor + 1))
    } else if (!key.ctrl && !key.meta && input) {
      const newValue = value.slice(0, cursor) + input + value.slice(cursor)
      onChange(newValue)
      setCursor(cursor + input.length)
    }
  })

  // Display with cursor
  const displayValue = value || ''
  const beforeCursor = displayValue.slice(0, cursor)
  const afterCursor = displayValue.slice(cursor)

  return (
    <Text dimColor={isDisabled}>
      {beforeCursor}
      <Text inverse={!isDisabled} dimColor={isDisabled}>
        ▌
      </Text>
      {afterCursor}
    </Text>
  )
}

export default TextInput
