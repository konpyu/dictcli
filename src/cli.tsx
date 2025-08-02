#!/usr/bin/env node
import { Command } from 'commander'
import { render } from 'ink'
import React from 'react'
import { App } from './components/App.js'
import type { Settings, Level, Topic, VoiceDisplayName } from './types/index.js'
import { loadSettings } from './storage/settings.js'
import { readFileSync } from 'fs'
import { fileURLToPath } from 'url'
import { dirname, join } from 'path'

// Get package.json for version
const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)
const packageJson = JSON.parse(readFileSync(join(__dirname, '..', 'package.json'), 'utf-8'))

const program = new Command()

const VALID_TOPICS: Topic[] = [
  'EverydayLife',
  'Travel',
  'Technology',
  'Health',
  'Entertainment',
  'Business',
  'Random',
]
const VALID_LEVELS: Level[] = ['CEFR_A1', 'CEFR_A2', 'CEFR_B1', 'CEFR_B2', 'CEFR_C1', 'CEFR_C2']
const VALID_VOICES: VoiceDisplayName[] = ['ALEX', 'SARA', 'EVAN', 'NOVA', 'NICK', 'FAYE']

const validateEnv = () => {
  if (!process.env.OPENAI_API_KEY) {
    console.error('Error: OPENAI_API_KEY environment variable is required')
    console.error('Please set it by running: export OPENAI_API_KEY=your-api-key')
    process.exit(1)
  }
}

const validateOptions = (options: Record<string, unknown>, defaultSettings: Settings): Settings => {
  // Validate topic
  const topicInput = String(options.topic)
  const topic = VALID_TOPICS.find((t) => t.toLowerCase() === topicInput.toLowerCase()) as
    | Topic
    | undefined
  if (!topic) {
    console.error(`Error: Invalid topic "${topic}". Valid topics are: ${VALID_TOPICS.join(', ')}`)
    process.exit(1)
  }

  // Validate level
  const level = options.level as Level
  if (!VALID_LEVELS.includes(level)) {
    console.error(`Error: Invalid level "${level}". Valid levels are: ${VALID_LEVELS.join(', ')}`)
    process.exit(1)
  }

  // Validate voice
  const voice = options.voice as VoiceDisplayName
  if (!VALID_VOICES.includes(voice)) {
    console.error(`Error: Invalid voice "${voice}". Valid voices are: ${VALID_VOICES.join(', ')}`)
    process.exit(1)
  }

  // Validate word count
  const wordCount = parseInt(options.words as string, 10)
  if (isNaN(wordCount) || wordCount < 5 || wordCount > 30) {
    console.error('Error: Word count must be a number between 5 and 30')
    process.exit(1)
  }

  // Validate speed
  const speed = parseFloat(options.speed as string)
  if (isNaN(speed) || speed < 0.5 || speed > 2.0) {
    console.error('Error: Speed must be a number between 0.5 and 2.0')
    process.exit(1)
  }

  return {
    topic: topic as Topic,
    level: level || defaultSettings.level,
    wordCount: wordCount || defaultSettings.wordCount,
    voice: voice || defaultSettings.voice,
    speed: speed || defaultSettings.speed,
  }
}

program
  .name('dictcli')
  .description('LLM-First Dictation TUI App for Japanese English learners')
  .version(packageJson.version)
  .option('--topic <topic>', `learning topic (${VALID_TOPICS.join(', ')})`)
  .option('--level <level>', `CEFR level (${VALID_LEVELS.join(', ')})`)
  .option('--words <number>', 'word count (5-30)')
  .option('--voice <voice>', `voice name (${VALID_VOICES.join(', ')})`)
  .option('--speed <number>', 'playback speed (0.5-2.0)')
  .parse()

const options = program.opts()

if (process.env.DICTCLI_DEBUG === 'true') {
  console.log('CLI Options received:', options)
}

validateEnv()

const main = async () => {
  // Load saved settings or use defaults
  const savedSettings = await loadSettings()

  // Merge CLI options with saved settings (CLI options take precedence)
  const mergedOptions = {
    topic: options.topic || savedSettings.topic,
    level: options.level || savedSettings.level,
    words: options.words || String(savedSettings.wordCount),
    voice: options.voice || savedSettings.voice,
    speed: options.speed || String(savedSettings.speed),
  }

  const initialSettings = validateOptions(mergedOptions, savedSettings)

  if (process.env.DICTCLI_DEBUG === 'true') {
    console.log('Saved settings:', savedSettings)
    console.log('CLI options:', options)
    console.log('Initial settings:', initialSettings)
  }

  const app = render(<App initialSettings={initialSettings} />)

  const handleExit = () => {
    app.unmount()
    process.exit(0)
  }

  process.on('SIGINT', handleExit)
  process.on('SIGTERM', handleExit)
}

main().catch((error) => {
  console.error('Failed to start app:', error)
  process.exit(1)
})
