import React from 'react'
import { describe, it, expect } from 'vitest'
import { render } from 'ink-testing-library'
import { StatusBar } from '../../../../src/components/common/StatusBar'

describe('StatusBar', () => {
  it('should render empty when no props provided', () => {
    const { lastFrame } = render(<StatusBar />)
    expect(lastFrame()).toBe('')
  })

  it('should show playing status', () => {
    const { lastFrame } = render(
      <StatusBar isPlaying={true} playbackSpeed={1.0} />
    )
    expect(lastFrame()).toContain('🔊 Playing…')
    expect(lastFrame()).toContain('⏩1×')
  })

  it('should show different playback speeds', () => {
    const { lastFrame } = render(
      <StatusBar isPlaying={true} playbackSpeed={1.2} />
    )
    expect(lastFrame()).toContain('⏩1.2×')
  })

  it('should show custom status when not playing', () => {
    const { lastFrame } = render(
      <StatusBar status="Ready to start" />
    )
    expect(lastFrame()).toContain('Ready to start')
  })

  it('should prioritize playing status over custom status', () => {
    const { lastFrame } = render(
      <StatusBar status="Custom status" isPlaying={true} playbackSpeed={1.0} />
    )
    expect(lastFrame()).toContain('🔊 Playing…')
    expect(lastFrame()).not.toContain('Custom status')
  })
})