import { getOpenAIClient } from './client.js'
import type { Round } from '../../types/index.js'

interface ScoringResult {
  score: number
  wer: number
  errors: Array<{
    expected: string
    actual: string
    explanation: string
  }>
  alternatives: string[]
}

export class ScorerService {
  async scoreAnswer(reference: string, userInput: string): Promise<ScoringResult> {
    const client = getOpenAIClient()

    const prompt = `You are an English teacher evaluating a dictation exercise for a Japanese learner.

Reference sentence: "${reference}"
Student's answer: "${userInput}"

Analyze the student's answer and provide:
1. Score out of 100 (based on accuracy)
2. Word Error Rate (WER) as a decimal (0.0 to 1.0)
3. List of errors with Japanese explanations
4. 2-3 alternative correct expressions in English

IMPORTANT SCORING RULES:
- Do NOT penalize missing or extra periods (.) or exclamation marks (!) at the end of sentences
- Both "Hello." and "Hello" are correct, as are "Great!" and "Great"
- Focus on word accuracy and spelling, not end-of-sentence punctuation
- Other punctuation errors within sentences (commas, question marks, etc.) should still be considered

Response format (JSON):
{
  "score": 85,
  "wer": 0.15,
  "errors": [
    {
      "expected": "word1",
      "actual": "word2",
      "explanation": "日本語での説明"
    }
  ],
  "alternatives": ["Alternative expression 1", "Alternative expression 2"]
}`

    try {
      const response = await client.chat.completions.create({
        model: 'gpt-5-mini',
        messages: [
          {
            role: 'system',
            content:
              'You are an English teacher providing feedback in Japanese for Japanese learners.',
          },
          {
            role: 'user',
            content: prompt,
          },
        ],
        response_format: { type: 'json_object' },
        reasoning_effort: 'minimal',
        verbosity: 'low',
      })

      const content = response.choices[0]?.message?.content
      if (!content) {
        throw new Error('No response from OpenAI')
      }

      const result = JSON.parse(content) as ScoringResult
      return result
    } catch (error) {
      console.error('Error scoring answer:', error)
      // Fallback to simple comparison
      return this.simpleFallbackScoring(reference, userInput)
    }
  }

  private simpleFallbackScoring(reference: string, userInput: string): ScoringResult {
    const refWords = reference.toLowerCase().split(/\s+/)
    const userWords = userInput.toLowerCase().split(/\s+/)

    const errors = []
    // let correctCount = 0

    for (let i = 0; i < Math.max(refWords.length, userWords.length); i++) {
      const expected = refWords[i] || ''
      const actual = userWords[i] || ''

      if (expected === actual) {
        // correctCount++
      } else if (expected || actual) {
        errors.push({
          expected,
          actual,
          explanation: '単語が一致しません',
        })
      }
    }

    const wer = errors.length / refWords.length
    const score = Math.round((1 - wer) * 100)

    return {
      score,
      wer,
      errors,
      alternatives: [],
    }
  }

  createRound(sentence: string, userInput: string, scoringResult: ScoringResult): Round {
    return {
      id: Date.now().toString(),
      sentence,
      userInput,
      score: scoringResult.score,
      wer: scoringResult.wer,
      errors: scoringResult.errors,
      alternatives: scoringResult.alternatives,
      timestamp: new Date(),
    }
  }
}
