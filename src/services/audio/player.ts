import playSound from 'play-sound'
import { spawn, ChildProcess } from 'child_process'
import { platform } from 'os'
import { promises as fs } from 'fs'
import {
  AudioPlayer,
  AudioPlayerOptions,
  AudioPlayerError,
  AudioPlayerStatus,
} from '../../types/index.js'

export class AudioPlayerService implements AudioPlayer {
  private status: AudioPlayerStatus = 'idle'
  private currentProcess: ChildProcess | null = null
  private playerInstance: ReturnType<typeof playSound> | null = null

  constructor() {
    this.playerInstance = playSound({
      // Configure play-sound for macOS with afplay
      players: platform() === 'darwin' ? ['afplay'] : undefined,
    })
  }

  /**
   * Play audio file with optional speed and volume control
   * @param filePath - Path to the audio file
   * @param options - Playback options including speed (0.5-2.0) and volume
   */
  async play(filePath: string, options: AudioPlayerOptions = {}): Promise<void> {
    try {
      // Validate file exists
      await this.validateAudioFile(filePath)

      // Validate speed range
      const speed = this.validateSpeed(options.speed)

      // Stop any currently playing audio
      await this.stop()

      this.status = 'playing'

      // Use afplay directly for macOS to support speed control
      if (platform() === 'darwin' && (speed !== 1.0 || options.volume !== undefined)) {
        await this.playWithAfplay(filePath, speed, options.volume)
      } else {
        await this.playWithPlaySound(filePath, options.volume)
      }
    } catch (error: unknown) {
      this.status = 'error'

      // If it's already an AudioPlayerError with a code, re-throw it
      if (error && typeof error === 'object' && 'code' in error) {
        throw error
      }

      // Otherwise, wrap it in a generic PLAYBACK_FAILED error
      throw this.createAudioError(
        'PLAYBACK_FAILED',
        `Failed to play audio: ${error instanceof Error ? error.message : 'Unknown error'}`,
        error instanceof Error ? error.message : undefined,
      )
    }
  }

  /**
   * Stop currently playing audio
   */
  async stop(): Promise<void> {
    try {
      if (this.currentProcess) {
        this.currentProcess.kill('SIGTERM')
        this.currentProcess = null
      }
      this.status = 'stopped'
    } catch (error) {
      console.error('Error stopping audio:', error)
      this.status = 'error'
    }
  }

  /**
   * Check if audio is currently playing
   */
  isPlaying(): boolean {
    return this.status === 'playing'
  }

  /**
   * Get current player status
   */
  getStatus(): AudioPlayerStatus {
    return this.status
  }

  /**
   * Play audio using afplay with speed control (macOS only)
   */
  private async playWithAfplay(filePath: string, speed: number, volume?: number): Promise<void> {
    return new Promise((resolve, reject) => {
      const args = [filePath]

      // Add speed control using -r flag (playback rate)
      if (speed !== 1.0) {
        args.unshift('-r', speed.toString())
      }

      // Add volume control if specified
      if (volume !== undefined) {
        const volumeLevel = Math.max(0, Math.min(1, volume))
        args.unshift('-v', volumeLevel.toString())
      }

      this.currentProcess = spawn('afplay', args)

      this.currentProcess.on('close', (code) => {
        this.currentProcess = null
        if (code === 0) {
          this.status = 'idle'
          resolve()
        } else {
          this.status = 'error'
          reject(new Error(`afplay exited with code ${code}`))
        }
      })

      this.currentProcess.on('error', (error) => {
        this.currentProcess = null
        this.status = 'error'
        reject(error)
      })
    })
  }

  /**
   * Play audio using play-sound library (fallback for other platforms)
   */
  private async playWithPlaySound(filePath: string, volume?: number): Promise<void> {
    return new Promise((resolve, reject) => {
      if (!this.playerInstance) {
        reject(new Error('Player instance not initialized'))
        return
      }

      const options: Record<string, unknown> = {}
      if (volume !== undefined) {
        options.volume = Math.max(0, Math.min(1, volume))
      }

      this.playerInstance.play(filePath, options, (err) => {
        this.status = 'idle'
        if (err) {
          this.status = 'error'
          reject(err)
        } else {
          resolve()
        }
      })
    })
  }

  /**
   * Validate that the audio file exists and is accessible
   */
  private async validateAudioFile(filePath: string): Promise<void> {
    try {
      await fs.access(filePath)
      const stats = await fs.stat(filePath)
      if (!stats.isFile()) {
        throw new Error('Path is not a file')
      }
    } catch (error) {
      throw this.createAudioError(
        'FILE_NOT_FOUND',
        `Audio file not found or not accessible: ${filePath}`,
        error instanceof Error ? error.message : undefined,
      )
    }
  }

  /**
   * Validate and normalize speed value
   */
  private validateSpeed(speed?: number): number {
    if (speed === undefined) {
      return 1.0
    }

    const normalizedSpeed = Math.max(0.5, Math.min(2.0, speed))

    if (speed !== normalizedSpeed) {
      console.warn(`Speed ${speed} is outside valid range (0.5-2.0), clamped to ${normalizedSpeed}`)
    }

    return normalizedSpeed
  }

  /**
   * Create a standardized audio error
   */
  private createAudioError(code: string, message: string, details?: string): AudioPlayerError {
    const error = new Error(message) as AudioPlayerError
    error.code = code
    error.details = details
    return error
  }
}

// Export a singleton instance for easy use
export const audioPlayer = new AudioPlayerService()
