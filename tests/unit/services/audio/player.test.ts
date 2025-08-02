import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { AudioPlayerService } from '../../../../src/services/audio/player.js'
import { promises as fs } from 'fs'
import { spawn } from 'child_process'
import playSound from 'play-sound'
import { platform } from 'os'

// Mock dependencies
vi.mock('fs', () => ({
  promises: {
    access: vi.fn(),
    stat: vi.fn(),
  },
}))

vi.mock('child_process', () => ({
  spawn: vi.fn(),
}))

vi.mock('play-sound', () => ({
  default: vi.fn(),
}))

vi.mock('os', () => ({
  platform: vi.fn(),
}))

describe('AudioPlayerService', () => {
  let audioPlayer: AudioPlayerService
  let mockProcess: any
  let mockPlaySoundInstance: any

  beforeEach(() => {
    // Reset all mocks
    vi.clearAllMocks()

    // Mock file system
    vi.mocked(fs.access).mockResolvedValue(undefined)
    vi.mocked(fs.stat).mockResolvedValue({ isFile: () => true } as any)

    // Mock process
    mockProcess = {
      kill: vi.fn(),
      on: vi.fn(),
    }
    vi.mocked(spawn).mockReturnValue(mockProcess as any)

    // Mock play-sound instance
    mockPlaySoundInstance = {
      play: vi.fn(),
    }
    vi.mocked(playSound).mockReturnValue(mockPlaySoundInstance)

    // Create new instance for each test
    audioPlayer = new AudioPlayerService()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('constructor', () => {
    it('should initialize with idle status', () => {
      expect(audioPlayer.getStatus()).toBe('idle')
      expect(audioPlayer.isPlaying()).toBe(false)
    })

    it('should configure play-sound with afplay for macOS', () => {
      vi.mocked(platform).mockReturnValue('darwin')
      new AudioPlayerService()
      expect(playSound).toHaveBeenCalledWith({
        players: ['afplay'],
      })
    })

    it('should configure play-sound without specific player for non-macOS', () => {
      vi.mocked(platform).mockReturnValue('linux')
      new AudioPlayerService()
      expect(playSound).toHaveBeenCalledWith({
        players: undefined,
      })
    })
  })

  describe('play', () => {
    const testFilePath = '/test/audio.mp3'

    it('should play audio file successfully with default options', async () => {
      vi.mocked(platform).mockReturnValue('linux')
      mockPlaySoundInstance.play.mockImplementation((file, options, callback) => {
        callback(null)
      })

      await audioPlayer.play(testFilePath)

      expect(fs.access).toHaveBeenCalledWith(testFilePath)
      expect(fs.stat).toHaveBeenCalledWith(testFilePath)
      expect(mockPlaySoundInstance.play).toHaveBeenCalledWith(
        testFilePath,
        {},
        expect.any(Function),
      )
      expect(audioPlayer.getStatus()).toBe('idle')
    })

    it('should use afplay with speed control on macOS', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation((event, callback) => {
        if (event === 'close') {
          callback(0) // Success
        }
      })

      await audioPlayer.play(testFilePath, { speed: 0.9 })

      expect(spawn).toHaveBeenCalledWith('afplay', ['-r', '0.9', testFilePath])
      expect(audioPlayer.getStatus()).toBe('idle')
    })

    it('should use afplay with volume control on macOS', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation((event, callback) => {
        if (event === 'close') {
          setTimeout(() => callback(0), 0)
        }
      })

      await audioPlayer.play(testFilePath, { volume: 0.5 })

      expect(spawn).toHaveBeenCalledWith('afplay', ['-v', '0.5', testFilePath])
    })

    it('should use afplay with both speed and volume control on macOS', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation((event, callback) => {
        if (event === 'close') {
          callback(0)
        }
      })

      await audioPlayer.play(testFilePath, { speed: 1.1, volume: 0.8 })

      expect(spawn).toHaveBeenCalledWith('afplay', [
        '-v',
        '0.8',
        '-r',
        '1.1',
        testFilePath,
      ])
    })

    it('should clamp speed to valid range (0.5-2.0)', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation((event, callback) => {
        if (event === 'close') {
          callback(0)
        }
      })

      // Test speed too low
      await audioPlayer.play(testFilePath, { speed: 0.3 })
      expect(spawn).toHaveBeenCalledWith('afplay', ['-r', '0.5', testFilePath])

      // Test speed too high
      await audioPlayer.play(testFilePath, { speed: 3.0 })
      expect(spawn).toHaveBeenCalledWith('afplay', ['-r', '2', testFilePath])
    })

    it('should clamp volume to valid range (0-1)', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation((event, callback) => {
        if (event === 'close') {
          setTimeout(() => callback(0), 0)
        }
      })

      // Test volume too low
      await audioPlayer.play(testFilePath, { volume: -0.5 })
      expect(spawn).toHaveBeenCalledWith('afplay', ['-v', '0', testFilePath])

      // Test volume too high
      await audioPlayer.play(testFilePath, { volume: 1.5 })
      expect(spawn).toHaveBeenCalledWith('afplay', ['-v', '1', testFilePath])
    })

    it('should throw error if file does not exist', async () => {
      vi.mocked(fs.access).mockRejectedValue(new Error('File not found'))

      await expect(audioPlayer.play(testFilePath)).rejects.toThrow(
        'Audio file not found or not accessible',
      )
      expect(audioPlayer.getStatus()).toBe('error')
    })

    it('should throw error if path is not a file', async () => {
      vi.mocked(fs.stat).mockResolvedValue({ isFile: () => false } as any)

      await expect(audioPlayer.play(testFilePath)).rejects.toThrow(
        'Audio file not found or not accessible',
      )
    })

    it('should throw error if afplay fails', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation((event, callback) => {
        if (event === 'close') {
          callback(1) // Error exit code
        }
      })

      await expect(audioPlayer.play(testFilePath, { speed: 0.9 })).rejects.toThrow(
        'afplay exited with code 1',
      )
      expect(audioPlayer.getStatus()).toBe('error')
    })

    it('should throw error if afplay process errors', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation((event, callback) => {
        if (event === 'error') {
          callback(new Error('Process error'))
        }
      })

      await expect(audioPlayer.play(testFilePath, { speed: 0.9 })).rejects.toThrow(
        'Process error',
      )
      expect(audioPlayer.getStatus()).toBe('error')
    })

    it('should throw error if play-sound fails', async () => {
      vi.mocked(platform).mockReturnValue('linux')
      mockPlaySoundInstance.play.mockImplementation((file, options, callback) => {
        callback(new Error('Playback failed'))
      })

      await expect(audioPlayer.play(testFilePath)).rejects.toThrow('Playback failed')
      expect(audioPlayer.getStatus()).toBe('error')
    })

    it('should call stop before starting new playback', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      
      // Spy on the stop method
      const stopSpy = vi.spyOn(audioPlayer, 'stop')
      
      // Mock spawn to complete immediately
      mockProcess.on.mockImplementation((event, callback) => {
        if (event === 'close') {
          setTimeout(() => callback(0), 0)
        }
      })

      await audioPlayer.play(testFilePath, { speed: 0.9 })

      // Stop should have been called once during the play method
      expect(stopSpy).toHaveBeenCalledTimes(1)
      
      stopSpy.mockRestore()
    })
  })

  describe('stop', () => {
    it('should stop current playback', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation((event, callback) => {
        // Don't call callback to simulate ongoing playback
      })

      // Start playback but don't await
      audioPlayer.play('/test/audio.mp3', { speed: 0.9 })
      
      // Wait for async operation to start
      await new Promise(resolve => setTimeout(resolve, 10))
      expect(audioPlayer.getStatus()).toBe('playing')

      await audioPlayer.stop()

      expect(mockProcess.kill).toHaveBeenCalledWith('SIGTERM')
      expect(audioPlayer.getStatus()).toBe('stopped')
    })

    it('should handle stop when no playback is active', async () => {
      await audioPlayer.stop()
      expect(audioPlayer.getStatus()).toBe('stopped')
    })

    it('should handle error during stop', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      
      // Start playback first
      mockProcess.on.mockImplementation(() => {})
      audioPlayer.play('/test/audio.mp3', { speed: 0.9 })
      
      // Wait for async operation to start
      await new Promise(resolve => setTimeout(resolve, 10))
      expect(audioPlayer.getStatus()).toBe('playing')
      
      // Mock kill to throw error
      mockProcess.kill.mockImplementation(() => {
        throw new Error('Kill failed')
      })

      // Stop should handle error gracefully
      await audioPlayer.stop()
      expect(audioPlayer.getStatus()).toBe('error')
    })
  })

  describe('isPlaying', () => {
    it('should return true when status is playing', async () => {
      vi.mocked(platform).mockReturnValue('darwin')
      mockProcess.on.mockImplementation(() => {}) // Don't complete

      // Start playback but don't await
      audioPlayer.play('/test/audio.mp3', { speed: 0.9 })
      
      // Wait for async operation to start
      await new Promise(resolve => setTimeout(resolve, 10))
      expect(audioPlayer.isPlaying()).toBe(true)
    })

    it('should return false when status is not playing', () => {
      expect(audioPlayer.isPlaying()).toBe(false)
    })
  })

  describe('getStatus', () => {
    it('should return current status', () => {
      expect(audioPlayer.getStatus()).toBe('idle')
    })
  })

  describe('error handling', () => {
    it('should create proper AudioPlayerError with code and details', async () => {
      vi.mocked(fs.access).mockRejectedValue(new Error('Permission denied'))

      try {
        await audioPlayer.play('/test/audio.mp3')
        expect.fail('Should have thrown an error')
      } catch (error: any) {
        expect(error.code).toBe('FILE_NOT_FOUND')
        expect(error.details).toBe('Permission denied')
        expect(error.message).toContain('Audio file not found or not accessible')
      }
    })

    it('should handle non-Error exceptions gracefully', async () => {
      vi.mocked(fs.access).mockRejectedValue('String error')

      try {
        await audioPlayer.play('/test/audio.mp3')
        expect.fail('Should have thrown an error')
      } catch (error: any) {
        expect(error.code).toBe('FILE_NOT_FOUND')
        expect(error.details).toBeUndefined()
      }
    })
  })
})