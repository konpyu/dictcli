import { promises as fs } from 'fs'
import { homedir } from 'os'
import { join } from 'path'
import type { Settings } from '../types/index.js'

const CONFIG_DIR = join(homedir(), '.dictcli')
const SETTINGS_FILE = join(CONFIG_DIR, 'settings.json')

const DEFAULT_SETTINGS: Settings = {
  voice: 'ALEX',
  level: 'CEFR_A1',
  topic: 'EverydayLife',
  wordCount: 10,
  speed: 1.0,
}

export class SettingsStorage {
  async ensureConfigDir(): Promise<void> {
    try {
      await fs.mkdir(CONFIG_DIR, { recursive: true })
    } catch (error) {
      console.error('Failed to create config directory:', error)
    }
  }

  async load(): Promise<Settings> {
    try {
      await this.ensureConfigDir()
      const data = await fs.readFile(SETTINGS_FILE, 'utf-8')
      return { ...DEFAULT_SETTINGS, ...JSON.parse(data) }
    } catch {
      return DEFAULT_SETTINGS
    }
  }

  async save(settings: Settings): Promise<void> {
    try {
      await this.ensureConfigDir()
      await fs.writeFile(SETTINGS_FILE, JSON.stringify(settings, null, 2))
    } catch (error) {
      console.error('Failed to save settings:', error)
    }
  }

  async update(partial: Partial<Settings>): Promise<Settings> {
    const current = await this.load()
    const updated = { ...current, ...partial }
    await this.save(updated)
    return updated
  }
}

// Convenience functions for easier use
const storage = new SettingsStorage()

export const loadSettings = async (): Promise<Settings> => {
  return storage.load()
}

export const getDefaultSettings = (): Settings => {
  return DEFAULT_SETTINGS
}

export const saveSettings = async (settings: Settings): Promise<void> => {
  await storage.save(settings)
}
