import { describe, it, expect, beforeEach } from 'vitest'
import { store, getState, setState } from '../../../src/store/useStore'

describe('useStore', () => {
  beforeEach(() => {
    // Reset store to initial state
    setState({ viewState: 'learning' })
  })

  it('should have initial view state as learning', () => {
    expect(getState().viewState).toBe('learning')
  })

  it('should update view state', () => {
    store.getState().setViewState('settings')
    expect(getState().viewState).toBe('settings')
  })

  it('should handle all view states', () => {
    const states: Array<'learning' | 'result' | 'settings'> = ['learning', 'result', 'settings']
    
    states.forEach(state => {
      store.getState().setViewState(state)
      expect(getState().viewState).toBe(state)
    })
  })
})