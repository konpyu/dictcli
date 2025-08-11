import { promises as fs } from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'
import type { Topic } from '../types/index.js'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

export interface Scene {
  id: number
  description: string
  descriptionEn: string
}

export class SceneLoader {
  private scenesCache: Map<Topic, Scene[]> = new Map()

  private getSceneFileName(topic: Topic): string {
    const topicFileMap: Record<Topic, string> = {
      EverydayLife: 'everyday_life',
      Travel: 'travel',
      Technology: 'business', // Use business scenes for Technology since technology.txt doesn't exist
      Health: 'health',
      Entertainment: 'entertainment',
      Business: 'business',
      Random: 'everyday_life', // fallback for Random
    }

    const baseName = topicFileMap[topic] || 'everyday_life'
    return `${baseName}_en.txt`
  }

  private async loadSceneFile(topic: Topic): Promise<string[]> {
    const fileName = this.getSceneFileName(topic)
    const filePath = path.join(__dirname, '..', '..', 'scenes', fileName)

    try {
      const content = await fs.readFile(filePath, 'utf-8')
      return content
        .split('\n')
        .filter((line) => line.trim())
        .map((line) => {
          // Remove line numbers and arrow if present (e.g., "1→" or just content)
          const match = line.match(/^\d+→(.+)$/)
          return match ? match[1].trim() : line.trim()
        })
    } catch (error) {
      console.error(`Failed to load scene file ${fileName}:`, error)
      return []
    }
  }

  async loadScenes(topic: Topic): Promise<Scene[]> {
    // Check cache first
    if (this.scenesCache.has(topic)) {
      return this.scenesCache.get(topic)!
    }

    // Handle Random topic by selecting a random actual topic
    const actualTopic = topic === 'Random' ? this.getRandomTopic() : topic

    const enScenes = await this.loadSceneFile(actualTopic)

    const scenes: Scene[] = []

    for (let i = 0; i < enScenes.length; i++) {
      scenes.push({
        id: i + 1,
        description: enScenes[i], // Use English text for both fields
        descriptionEn: enScenes[i],
      })
    }

    // Cache the result
    this.scenesCache.set(topic, scenes)
    return scenes
  }

  getRandomScene(scenes: Scene[]): Scene {
    if (scenes.length === 0) {
      // Fallback scene if no scenes are loaded
      return {
        id: 0,
        description: 'Daily conversation',
        descriptionEn: 'Daily conversation',
      }
    }
    return scenes[Math.floor(Math.random() * scenes.length)]
  }

  private getRandomTopic(): Topic {
    const topics: Topic[] = [
      'EverydayLife',
      'Travel',
      'Technology',
      'Health',
      'Entertainment',
      'Business',
    ]
    return topics[Math.floor(Math.random() * topics.length)]
  }

  async getRandomSceneForTopic(topic: Topic): Promise<Scene> {
    const scenes = await this.loadScenes(topic)
    return this.getRandomScene(scenes)
  }
}

export const sceneLoader = new SceneLoader()
