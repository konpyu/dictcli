import { getOpenAIClient } from './client.js'
import type { Level, Topic } from '../../types/index.js'
import { sceneLoader } from '../sceneLoader.js'

const LEVEL_DESCRIPTIONS: Record<Level, string> = {
  CEFR_A1: 'very simple sentences with basic vocabulary, present tense only',
  CEFR_A2: 'simple sentences with common vocabulary, present and past tense',
  CEFR_B1: 'intermediate sentences with everyday expressions and phrasal verbs',
  CEFR_B2: 'complex sentences with advanced vocabulary and various tenses',
  CEFR_C1: 'sophisticated sentences with idioms and nuanced expressions',
  CEFR_C2: 'native-level sentences with complex structures and rare vocabulary',
}

const CONTEXT_VARIATIONS: Record<string, string[]> = {
  Business: [
    'during a meeting',
    'in an email',
    'at a conference',
    'in a presentation',
    'negotiating a deal',
  ],
  Tech: [
    'debugging code',
    'explaining features',
    'in documentation',
    'tech support',
    'discussing AI',
  ],
  Travel: [
    'at the airport',
    'in a hotel',
    'asking for directions',
    'ordering food',
    'planning an itinerary',
  ],
  Daily: ['at home', 'shopping', 'with friends', 'on the phone', 'running errands'],
  Technology: [
    'debugging code',
    'explaining features',
    'in documentation',
    'tech support',
    'discussing AI',
  ],
  Health: [
    'at the doctor',
    'discussing symptoms',
    'fitness advice',
    'mental wellness',
    'nutrition planning',
  ],
}

const SENTENCE_STRUCTURES = [
  'statement',
  'question',
  'exclamation',
  'conditional (if...then)',
  'comparative',
  'passive voice',
  'reported speech',
  'complex with multiple clauses',
]

const TONES = [
  'formal',
  'casual',
  'enthusiastic',
  'skeptical',
  'humorous',
  'serious',
  'encouraging',
  'cautious',
]

export class ProblemGenerator {
  async generateProblem(level: Level, topic: Topic, wordCount: number): Promise<string> {
    const client = getOpenAIClient()

    // ランダムに要素を選択
    const REAL_TOPICS: Topic[] = [
      'EverydayLife',
      'Travel',
      'Technology',
      'Health',
      'Entertainment',
      'Business',
    ]
    const effectiveTopic: Topic =
      topic === 'Random' ? REAL_TOPICS[Math.floor(Math.random() * REAL_TOPICS.length)] : topic

    // Load a random scene for the topic
    const scene = await sceneLoader.getRandomSceneForTopic(effectiveTopic)

    const contexts = CONTEXT_VARIATIONS[effectiveTopic] || ['general situation']
    const context = contexts[Math.floor(Math.random() * contexts.length)]
    const structure = SENTENCE_STRUCTURES[Math.floor(Math.random() * SENTENCE_STRUCTURES.length)]
    const tone = TONES[Math.floor(Math.random() * TONES.length)]

    // Build scene context for the prompt
    const sceneContext = scene.descriptionEn
      ? `\nScene context: "${scene.descriptionEn}" - Create a sentence that relates to or could occur in this specific scene/situation.`
      : ''

    const prompt = `Generate a REALISTIC and NATURAL English sentence that people would actually say in real-life situations.

Requirements:
- Level: ${level} (${LEVEL_DESCRIPTIONS[level]})
- Topic: ${effectiveTopic}
- Context: ${context}${sceneContext}
- Sentence structure: ${structure}
- Tone: ${tone}
- Word count: exactly ${wordCount} words

CRITICAL RULES:
1. The sentence MUST be something people would actually say in real conversation
2. Ensure grammatical correctness - comparatives need comparison targets, proper verb agreement, etc.
3. Avoid artificial or textbook-style sentences like "This app is easier today" or "The pen is on the table"
4. Create authentic dialogue or thoughts that fit the scene naturally
5. For short sentences (1-3 words), use common expressions like "No problem!", "Got it!", "Coming through!"
6. Include appropriate punctuation and contractions when natural

Examples of GOOD sentences:
- "Could you send me the report by Friday?" (Business, 8 words)
- "I'm running late!" (Daily, 3 words)
- "The meeting's been pushed back an hour." (Business, 7 words)

Examples of BAD sentences (DO NOT generate these):
- "This app is easier today." (grammatically incorrect, unnatural)
- "The cat sits on the mat." (textbook-style, nobody says this)
- "I see a red car." (too artificial)

Return only the sentence, nothing else.`

    // Add debug logging for prompt details
    if (process.env.DICTCLI_DEBUG === 'true') {
      console.log('Generating problem with:', {
        level,
        topic,
        wordCount,
        context,
        structure,
        tone,
        scene: scene.descriptionEn,
        levelDescription: LEVEL_DESCRIPTIONS[level],
      })
    }

    try {
      const response = await client.chat.completions.create({
        model: 'gpt-5-mini',
        messages: [
          {
            role: 'system',
            content: 'You are an English teacher creating dictation exercises.',
          },
          {
            role: 'user',
            content: prompt,
          },
        ],
        reasoning_effort: 'minimal',
        verbosity: 'low',
      })

      const sentence = response.choices[0]?.message?.content?.trim()
      if (!sentence) {
        throw new Error('Failed to generate sentence')
      }

      return sentence
    } catch (error) {
      console.error('Error generating problem:', error)
      throw new Error('Failed to generate problem')
    }
  }
}
