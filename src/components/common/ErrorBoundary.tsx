import React from 'react'
import { Box, Text } from 'ink'

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

interface ErrorBoundaryProps {
  children: React.ReactNode
}

class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo)
  }

  render() {
    if (this.state.hasError) {
      return (
        <Box flexDirection="column" paddingX={2} paddingY={1}>
          <Text color="red" bold>
            ⚠️ エラーが発生しました
          </Text>
          <Box marginTop={1}>
            <Text>{this.state.error?.message || 'Unknown error'}</Text>
          </Box>
          <Box marginTop={1}>
            <Text color="gray">アプリケーションを再起動してください: Ctrl+C → npm run dev</Text>
          </Box>
        </Box>
      )
    }

    return this.props.children
  }
}

export default ErrorBoundary
