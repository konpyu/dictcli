import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { promises as fs } from 'fs'
import { HistoryStorage } from '../../../src/storage/history.js'
import type { Round } from '../../../src/types/index.js'

vi.mock('fs', () => ({
  promises: {
    mkdir: vi.fn(),
    readFile: vi.fn(),
    appendFile: vi.fn(),
  },
}))

describe('HistoryStorage', () => {
  let storage: HistoryStorage

  beforeEach(() => {
    storage = new HistoryStorage()
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  const createMockRound = (overrides?: Partial<Round>): Round => ({
    id: '123',
    sentence: 'Test sentence',
    userInput: 'Test input',
    score: 85,
    wer: 0.15,
    errors: [],
    alternatives: [],
    timestamp: new Date('2025-01-15'),
    ...overrides,
  })

  describe('append', () => {
    it('should append round to history file', async () => {
      const round = createMockRound()

      await storage.append(round)

      expect(fs.appendFile).toHaveBeenCalledWith(
        expect.stringContaining('history.jsonl'),
        JSON.stringify(round) + '\n',
      )
    })
  })

  describe('getAll', () => {
    it('should return empty array when file does not exist', async () => {
      vi.mocked(fs.readFile).mockRejectedValue(new Error('File not found'))

      const rounds = await storage.getAll()

      expect(rounds).toEqual([])
    })

    it('should parse all rounds from file', async () => {
      const round1 = createMockRound({ id: '1' })
      const round2 = createMockRound({ id: '2' })
      const jsonl = `${JSON.stringify(round1)}\n${JSON.stringify(round2)}\n`
      vi.mocked(fs.readFile).mockResolvedValue(jsonl)

      const rounds = await storage.getAll()

      expect(rounds).toHaveLength(2)
      expect(rounds[0].id).toBe('1')
      expect(rounds[1].id).toBe('2')
    })
  })

  describe('getRecent', () => {
    it('should return specified number of recent rounds', async () => {
      const rounds = Array.from({ length: 10 }, (_, i) => createMockRound({ id: String(i) }))
      const jsonl = rounds.map((r) => JSON.stringify(r)).join('\n')
      vi.mocked(fs.readFile).mockResolvedValue(jsonl)

      const recent = await storage.getRecent(3)

      expect(recent).toHaveLength(3)
      expect(recent[0].id).toBe('7')
      expect(recent[1].id).toBe('8')
      expect(recent[2].id).toBe('9')
    })
  })

  describe('calculateStats', () => {
    it('should calculate statistics for recent rounds', async () => {
      const now = new Date()
      const rounds = [
        createMockRound({ score: 80, wer: 0.2, timestamp: new Date(now.getTime() - 1000 * 60 * 60 * 24) }),
        createMockRound({ score: 90, wer: 0.1, timestamp: new Date(now.getTime() - 1000 * 60 * 60) }),
        createMockRound({ score: 85, wer: 0.15, timestamp: new Date(now.getTime() - 1000 * 60 * 60 * 24 * 10) }),
      ]
      const jsonl = rounds.map((r) => JSON.stringify(r)).join('\n')
      vi.mocked(fs.readFile).mockResolvedValue(jsonl)

      const stats = await storage.calculateStats(7)

      expect(stats.totalRounds).toBe(2)
      expect(stats.averageScore).toBe(85)
      expect(stats.averageWER).toBeCloseTo(0.15)
    })

    it('should return zeros when no recent rounds', async () => {
      vi.mocked(fs.readFile).mockResolvedValue('')

      const stats = await storage.calculateStats(7)

      expect(stats).toEqual({
        totalRounds: 0,
        averageScore: 0,
        averageWER: 0,
      })
    })
  })
})