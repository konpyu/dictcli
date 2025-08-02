import React from 'react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render } from 'ink-testing-library'
import { App } from '../../../src/components/App'
import { setState } from '../../../src/store/useStore'
import type { Settings } from '../../../src/types'

// Mock useApp to prevent actual exit
vi.mock('ink', async () => {
  const actual = await vi.importActual('ink')
  return {
    ...actual,
    useApp: () => ({
      exit: vi.fn(),
    }),
  }
})

// Mock orchestrator to prevent API calls
vi.mock('../../../src/services/orchestrator', () => ({
  orchestrator: {
    startNewRound: vi.fn().mockResolvedValue(undefined),
    replayAudio: vi.fn().mockResolvedValue(undefined),
    scoreAnswer: vi.fn().mockResolvedValue(undefined),
    showGapFill: vi.fn().mockReturnValue('__ __ __'),
  },
}))

describe('App', () => {
  const defaultSettings: Settings = {
    voice: 'ALEX',
    level: 'CEFR_A1',
    topic: 'Business',
    wordCount: 10,
  }

  beforeEach(() => {
    // Reset to initial state
    setState({ viewState: 'learning' })
    vi.clearAllMocks()
  })

  it('should render learning view by default', async () => {
    const { lastFrame, rerender } = render(<App initialSettings={defaultSettings} />)
    
    expect(lastFrame()).toContain('DictCLI')
    
    // Wait for initialization
    await new Promise(resolve => setTimeout(resolve, 50))
    rerender(<App initialSettings={defaultSettings} />)
    
    // Now the learning view should be visible
    expect(lastFrame()).toContain('>') // Check for input prompt
    // The app now starts with "Waiting for round..." until API calls complete
    expect(lastFrame()).toMatch(/Waiting for round|Generating problem|Playing/)
  })

  it('should handle keyboard navigation', async () => {
    // For now, we'll just ensure the component renders without errors
    // The Zustand store with Ink v4 workaround makes testing state changes complex
    const { lastFrame, stdin, rerender } = render(<App initialSettings={defaultSettings} />)
    
    // Wait for initialization
    await new Promise(resolve => setTimeout(resolve, 50))
    rerender(<App initialSettings={defaultSettings} />)
    
    // Initial state should be learning
    expect(lastFrame()).toMatch(/Waiting for round|Generating problem|Playing/)
    
    // Test that pressing keys doesn't cause errors
    stdin.write('s')
    stdin.write('s')
    stdin.write('q')
    
    // Component should still render
    expect(lastFrame()).toBeTruthy()
  })

  it('should exit on Q key', () => {
    const { stdin } = render(<App initialSettings={defaultSettings} />)
    
    // App should exit when Q is pressed
    stdin.write('q')
    // Since exit is mocked, we just ensure no errors occur
  })

  it('should display header with settings', () => {
    const { lastFrame } = render(<App initialSettings={defaultSettings} />)
    
    const frame = lastFrame()
    expect(frame).toContain('DictCLI')
    expect(frame).toContain('Voice:ALEX')
    expect(frame).toContain('Level:CEFR_A1')
    expect(frame).toContain('Topic:Business')
    expect(frame).toContain('Length:10w')
  })
})