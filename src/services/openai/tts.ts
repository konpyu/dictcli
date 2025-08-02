import { promises as fs } from 'fs'
import { join } from 'path'
import { tmpdir } from 'os'
import { createHash } from 'crypto'
import { getOpenAIClient } from './client.js'
import { VOICE_MAPPING, type VoiceDisplayName } from '../../types/index.js'

const CACHE_DIR = join(tmpdir(), 'dictcli-audio-cache')
const CACHE_TTL = 15 * 60 * 1000 // 15 minutes
// const MAX_CACHE_SIZE = 100 * 1024 * 1024 // 100MB (unused for now)

export class TTSService {
  async ensureCacheDir(): Promise<void> {
    try {
      await fs.mkdir(CACHE_DIR, { recursive: true })
    } catch (error) {
      console.error('Failed to create cache directory:', error)
    }
  }

  private getCacheKey(text: string, voice: VoiceDisplayName, speed: number): string {
    const data = `${text}-${voice}-${speed}`
    return createHash('md5').update(data).digest('hex')
  }

  private getCachePath(key: string, format: string = 'mp3'): string {
    return join(CACHE_DIR, `${key}.${format}`)
  }

  async cleanupOldCache(): Promise<void> {
    try {
      const files = await fs.readdir(CACHE_DIR)
      const now = Date.now()

      for (const file of files) {
        const path = join(CACHE_DIR, file)
        const stats = await fs.stat(path)
        if (now - stats.mtime.getTime() > CACHE_TTL) {
          await fs.unlink(path)
        }
      }
    } catch (error) {
      console.error('Cache cleanup failed:', error)
    }
  }

  async generateSpeech(
    text: string,
    voice: VoiceDisplayName,
    speed: number = 1.0,
  ): Promise<string> {
    await this.ensureCacheDir()

    const cacheKey = this.getCacheKey(text, voice, speed)
    const cachePath = this.getCachePath(cacheKey)

    // Check cache
    try {
      await fs.access(cachePath)
      const stats = await fs.stat(cachePath)
      if (Date.now() - stats.mtime.getTime() < CACHE_TTL) {
        return cachePath
      }
    } catch {
      // Cache miss, continue to generate
    }

    // Generate new audio
    const client = getOpenAIClient()
    const openAIVoice = VOICE_MAPPING[voice]

    try {
      const response = await client.audio.speech.create({
        model: 'tts-1-hd', // HD model for better quality
        voice: openAIVoice,
        input: text,
        speed,
        response_format: 'mp3', // Can change to 'flac' or 'wav' for higher quality
      })

      const buffer = Buffer.from(await response.arrayBuffer())
      await fs.writeFile(cachePath, buffer)

      // Cleanup old cache periodically
      if (Math.random() < 0.1) {
        this.cleanupOldCache().catch(() => {})
      }

      return cachePath
    } catch (error) {
      console.error('Error generating audio:', error)
      throw new Error('Failed to generate audio')
    }
  }

  async getCachedPath(
    text: string,
    voice: VoiceDisplayName,
    speed: number = 1.0,
  ): Promise<string | null> {
    const cacheKey = this.getCacheKey(text, voice, speed)
    const cachePath = this.getCachePath(cacheKey)

    try {
      await fs.access(cachePath)
      const stats = await fs.stat(cachePath)
      if (Date.now() - stats.mtime.getTime() < CACHE_TTL) {
        return cachePath
      }
    } catch {
      // Cache miss
    }

    return null
  }
}
