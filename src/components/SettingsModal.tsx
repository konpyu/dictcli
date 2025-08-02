import React, { useState } from 'react'
import { Box, Text } from 'ink'
import { useInput, type Key } from 'ink'
import { saveSettings } from '../storage/settings.js'
import { useStore, store } from '../store/useStore.js'
import type { Settings, VoiceDisplayName, Level, Topic } from '../types/index.js'

interface SettingsModalProps {
  onClose: () => void
  onSave: () => void
}

const voices: VoiceDisplayName[] = ['ALEX', 'SARA', 'EVAN', 'NOVA', 'NICK', 'FAYE']
const levels: Level[] = ['CEFR_A1', 'CEFR_A2', 'CEFR_B1', 'CEFR_B2', 'CEFR_C1', 'CEFR_C2']
const topics: Topic[] = [
  'Random',
  'EverydayLife',
  'Travel',
  'Technology',
  'Health',
  'Entertainment',
  'Business',
]

const SettingsModal: React.FC<SettingsModalProps> = ({ onClose, onSave }) => {
  const storeSettings = useStore((state) => state.settings)
  const updateStoreSettings = useStore((state) => state.updateSettings)
  const [settings, setSettings] = useState<Settings>(storeSettings)
  const [selectedField, setSelectedField] = useState<
    'voice' | 'level' | 'topic' | 'wordCount' | 'speed'
  >('voice')

  // Voice icons
  const getVoiceIcon = (voice: VoiceDisplayName) => {
    return ['ALEX', 'EVAN', 'NICK'].includes(voice) ? '♂' : '♀'
  }

  useInput((input: string, key: Key) => {
    if (key.escape) {
      onClose()
    } else if (key.return) {
      // Update both store and file
      updateStoreSettings(settings)
      saveSettings(settings).catch((error) => {
        console.error('Failed to save settings:', error)
      })

      // Update audio state with new speed
      store.getState().updateAudioState({ speed: settings.speed })

      // Clear pre-generated round when settings change
      store.getState().setPreGeneratedRound(null)

      onSave()
    } else if (key.upArrow) {
      const fields: Array<'voice' | 'level' | 'topic' | 'wordCount' | 'speed'> = [
        'voice',
        'level',
        'topic',
        'wordCount',
        'speed',
      ]
      const currentIndex = fields.indexOf(selectedField)
      if (currentIndex > 0) {
        setSelectedField(fields[currentIndex - 1])
      }
    } else if (key.downArrow) {
      const fields: Array<'voice' | 'level' | 'topic' | 'wordCount' | 'speed'> = [
        'voice',
        'level',
        'topic',
        'wordCount',
        'speed',
      ]
      const currentIndex = fields.indexOf(selectedField)
      if (currentIndex < fields.length - 1) {
        setSelectedField(fields[currentIndex + 1])
      }
    } else if (key.leftArrow) {
      switch (selectedField) {
        case 'voice': {
          const voiceIndex = voices.indexOf(settings.voice)
          if (voiceIndex > 0) {
            setSettings({ ...settings, voice: voices[voiceIndex - 1] })
          }
          break
        }
        case 'level': {
          const levelIndex = levels.indexOf(settings.level)
          if (levelIndex > 0) {
            setSettings({ ...settings, level: levels[levelIndex - 1] })
          }
          break
        }
        case 'topic': {
          const topicIndex = topics.indexOf(settings.topic)
          if (topicIndex > 0) {
            setSettings({ ...settings, topic: topics[topicIndex - 1] })
          }
          break
        }
        case 'wordCount':
          if (settings.wordCount > 5) {
            setSettings({ ...settings, wordCount: settings.wordCount - 1 })
          }
          break
        case 'speed':
          if (settings.speed > 0.5) {
            const newSpeed = Math.round((settings.speed - 0.1) * 10) / 10
            setSettings({ ...settings, speed: Math.max(0.5, newSpeed) })
          }
          break
      }
    } else if (key.rightArrow) {
      switch (selectedField) {
        case 'voice': {
          const voiceIndex = voices.indexOf(settings.voice)
          if (voiceIndex < voices.length - 1) {
            setSettings({ ...settings, voice: voices[voiceIndex + 1] })
          }
          break
        }
        case 'level': {
          const levelIndex = levels.indexOf(settings.level)
          if (levelIndex < levels.length - 1) {
            setSettings({ ...settings, level: levels[levelIndex + 1] })
          }
          break
        }
        case 'topic': {
          const topicIndex = topics.indexOf(settings.topic)
          if (topicIndex < topics.length - 1) {
            setSettings({ ...settings, topic: topics[topicIndex + 1] })
          }
          break
        }
        case 'wordCount':
          if (settings.wordCount < 30) {
            setSettings({ ...settings, wordCount: settings.wordCount + 1 })
          }
          break
        case 'speed':
          if (settings.speed < 2.0) {
            const newSpeed = Math.round((settings.speed + 0.1) * 10) / 10
            setSettings({ ...settings, speed: Math.min(2.0, newSpeed) })
          }
          break
      }
    } else if (input === '-') {
      if (selectedField === 'wordCount' && settings.wordCount > 5) {
        setSettings({ ...settings, wordCount: settings.wordCount - 1 })
      } else if (selectedField === 'speed' && settings.speed > 0.5) {
        const newSpeed = Math.round((settings.speed - 0.1) * 10) / 10
        setSettings({ ...settings, speed: Math.max(0.5, newSpeed) })
      }
    } else if (input === '+' || input === '=') {
      if (selectedField === 'wordCount' && settings.wordCount < 30) {
        setSettings({ ...settings, wordCount: settings.wordCount + 1 })
      } else if (selectedField === 'speed' && settings.speed < 2.0) {
        const newSpeed = Math.round((settings.speed + 0.1) * 10) / 10
        setSettings({ ...settings, speed: Math.min(2.0, newSpeed) })
      }
    }
  })

  return (
    <Box flexDirection="column" borderStyle="single" borderColor="blue" paddingX={2} paddingY={1}>
      <Box marginBottom={1}>
        <Text bold color="blue">
          Settings
        </Text>
      </Box>

      <Box marginBottom={1}>
        <Text>{'─'.repeat(40)}</Text>
      </Box>

      {/* Voice */}
      <Box marginBottom={1}>
        <Text color={selectedField === 'voice' ? 'blue' : 'white'}>
          Voice : <Text bold>{settings.voice}</Text> (←/→) {getVoiceIcon(settings.voice)} voices
        </Text>
      </Box>

      {/* Level */}
      <Box marginBottom={1}>
        <Text color={selectedField === 'level' ? 'blue' : 'white'}>
          Level : <Text bold>{settings.level}</Text> (←/→) {levels.join(' ')}
        </Text>
      </Box>

      {/* Topic */}
      <Box marginBottom={1}>
        <Text color={selectedField === 'topic' ? 'blue' : 'white'}>
          Topic : <Text bold>{settings.topic}</Text> (←/→) 選択可能
        </Text>
      </Box>

      {/* Word Count */}
      <Box marginBottom={1}>
        <Text color={selectedField === 'wordCount' ? 'blue' : 'white'}>
          Length : <Text bold>{settings.wordCount} words</Text> (−/＋) 5–30
        </Text>
      </Box>

      {/* Speed */}
      <Box marginBottom={1}>
        <Text color={selectedField === 'speed' ? 'blue' : 'white'}>
          Speed : <Text bold>{settings.speed.toFixed(1)}×</Text> (←/→) 0.5–2.0
        </Text>
      </Box>

      <Box marginTop={1}>
        <Text>{'─'.repeat(40)}</Text>
      </Box>

      <Box marginTop={1}>
        <Text>
          [<Text color="green">Enter</Text>] Save & Next Round [<Text color="red">Esc</Text>] Cancel
        </Text>
      </Box>
    </Box>
  )
}

export default SettingsModal
