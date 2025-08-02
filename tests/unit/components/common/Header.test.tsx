import React from 'react'
import { describe, it, expect } from 'vitest'
import { render } from 'ink-testing-library'
import { Header } from '../../../../src/components/common/Header'
import type { Settings } from '../../../../src/types'

describe('Header', () => {
  const defaultSettings: Settings = {
    voice: 'ALEX',
    level: 'CEFR_A1',
    topic: 'Business',
    wordCount: 10,
    speed: 1,
  }

  it('should render with default settings', () => {
    const { lastFrame } = render(<Header settings={defaultSettings} />)
    
    expect(lastFrame()).toContain('DictCLI')
    expect(lastFrame()).toContain('Voice:ALEX')
    expect(lastFrame()).toContain('Level:CEFR_A1')
    expect(lastFrame()).toContain('Topic:Business')
    expect(lastFrame()).toContain('Length:10w')
  })

  it('should update when settings change', () => {
    const updatedSettings: Settings = {
      voice: 'SARA',
      level: 'CEFR_B2',
      topic: 'Tech',
      wordCount: 20,
      speed: 1,
    }
    
    const { lastFrame } = render(<Header settings={updatedSettings} />)
    
    expect(lastFrame()).toContain('Voice:SARA')
    expect(lastFrame()).toContain('Level:CEFR_B2')
    expect(lastFrame()).toContain('Topic:Tech')
    expect(lastFrame()).toContain('Length:20w')
  })
})