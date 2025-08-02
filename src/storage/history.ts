import { promises as fs } from 'fs'
import { homedir } from 'os'
import { join } from 'path'
import type { Round } from '../types/index.js'

const CONFIG_DIR = join(homedir(), '.dictcli')
const HISTORY_FILE = join(CONFIG_DIR, 'history.jsonl')

export class HistoryStorage {
  async ensureConfigDir(): Promise<void> {
    try {
      await fs.mkdir(CONFIG_DIR, { recursive: true })
    } catch (error) {
      console.error('Failed to create config directory:', error)
    }
  }

  async append(round: Round): Promise<void> {
    try {
      await this.ensureConfigDir()
      const line = JSON.stringify(round) + '\n'
      await fs.appendFile(HISTORY_FILE, line)
    } catch (error) {
      console.error('Failed to append to history:', error)
    }
  }

  async getAll(): Promise<Round[]> {
    try {
      await this.ensureConfigDir()
      const data = await fs.readFile(HISTORY_FILE, 'utf-8')
      return data
        .split('\n')
        .filter((line) => line.trim())
        .map((line) => JSON.parse(line))
    } catch {
      return []
    }
  }

  async getRecent(count: number): Promise<Round[]> {
    const all = await this.getAll()
    return all.slice(-count)
  }

  async calculateStats(days: number = 7): Promise<{
    totalRounds: number
    averageScore: number
    averageWER: number
  }> {
    const cutoff = new Date()
    cutoff.setDate(cutoff.getDate() - days)

    const rounds = await this.getAll()
    const recentRounds = rounds.filter((r) => new Date(r.timestamp) > cutoff)

    if (recentRounds.length === 0) {
      return { totalRounds: 0, averageScore: 0, averageWER: 0 }
    }

    const totalScore = recentRounds.reduce((sum, r) => sum + r.score, 0)
    const totalWER = recentRounds.reduce((sum, r) => sum + r.wer, 0)

    return {
      totalRounds: recentRounds.length,
      averageScore: totalScore / recentRounds.length,
      averageWER: totalWER / recentRounds.length,
    }
  }
}

// Create singleton instance
const historyStorage = new HistoryStorage()

// Export convenience function for saving history
export async function saveHistory(round: Round): Promise<void> {
  return historyStorage.append(round)
}
