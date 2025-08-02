import { promises as fs } from 'fs'
import { join } from 'path'

export async function getCacheSize(cacheDir: string): Promise<number> {
  try {
    const files = await fs.readdir(cacheDir)
    let totalSize = 0

    for (const file of files) {
      const path = join(cacheDir, file)
      const stats = await fs.stat(path)
      totalSize += stats.size
    }

    return totalSize
  } catch {
    return 0
  }
}

export async function pruneCache(cacheDir: string, maxSize: number): Promise<void> {
  try {
    const files = await fs.readdir(cacheDir)
    const fileStats = await Promise.all(
      files.map(async (file) => {
        const path = join(cacheDir, file)
        const stats = await fs.stat(path)
        return { path, size: stats.size, mtime: stats.mtime.getTime() }
      }),
    )

    // Sort by modification time (oldest first)
    fileStats.sort((a, b) => a.mtime - b.mtime)

    let totalSize = fileStats.reduce((sum, file) => sum + file.size, 0)

    // Remove oldest files until under maxSize
    for (const file of fileStats) {
      if (totalSize <= maxSize) break
      await fs.unlink(file.path)
      totalSize -= file.size
    }
  } catch (error) {
    console.error('Cache pruning failed:', error)
  }
}
