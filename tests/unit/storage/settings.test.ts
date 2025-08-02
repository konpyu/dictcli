import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { promises as fs } from 'fs'
import { SettingsStorage } from '../../../src/storage/settings.js'
import type { Settings } from '../../../src/types/index.js'

vi.mock('fs', () => ({
  promises: {
    mkdir: vi.fn(),
    readFile: vi.fn(),
    writeFile: vi.fn(),
  },
}))

describe('SettingsStorage', () => {
  let storage: SettingsStorage

  beforeEach(() => {
    storage = new SettingsStorage()
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('load', () => {
    it('should return default settings when file does not exist', async () => {
      vi.mocked(fs.readFile).mockRejectedValue(new Error('File not found'))

      const settings = await storage.load()

      expect(settings).toEqual({
        voice: 'ALEX',
        level: 'CEFR_A1',
        topic: 'EverydayLife',
        wordCount: 10,
        speed: 1,
      })
    })

    it('should load settings from file', async () => {
      const savedSettings: Settings = {
        voice: 'SARA',
        level: 'CEFR_B1',
        topic: 'Technology',
        wordCount: 15,
        speed: 1,
      }
      vi.mocked(fs.readFile).mockResolvedValue(JSON.stringify(savedSettings))

      const settings = await storage.load()

      expect(settings).toEqual(savedSettings)
    })
  })

  describe('save', () => {
    it('should save settings to file', async () => {
      const settings: Settings = {
        voice: 'NOVA',
        level: 'CEFR_B2',
        topic: 'Travel',
        wordCount: 20,
        speed: 1,
      }

      await storage.save(settings)

      expect(fs.writeFile).toHaveBeenCalledWith(
        expect.stringContaining('settings.json'),
        JSON.stringify(settings, null, 2),
      )
    })
  })

  describe('update', () => {
    it('should merge partial settings with existing ones', async () => {
      const existingSettings: Settings = {
        voice: 'ALEX',
        level: 'CEFR_A1',
        topic: 'EverydayLife',
        wordCount: 10,
        speed: 1,
      }
      vi.mocked(fs.readFile).mockResolvedValue(JSON.stringify(existingSettings))

      const updated = await storage.update({ voice: 'SARA', wordCount: 15 })

      expect(updated).toEqual({
        voice: 'SARA',
        level: 'CEFR_A1',
        topic: 'EverydayLife',
        wordCount: 15,
        speed: 1,
      })
    })
  })
})