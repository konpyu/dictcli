import { getOpenAIClient } from './client.js'
import type { Round } from '../../types/index.js'
import { localeService } from '../locale.js'

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
    const feedbackLanguage = localeService.getFullLanguageName()

    const prompt = `You are an English teacher evaluating a dictation exercise.

Reference sentence: "${reference}"
Student's answer: "${userInput}"

Analyze the student's answer and provide:
1. Score out of 100 (based on accuracy)
2. Word Error Rate (WER) as a decimal (0.0 to 1.0)
3. List of errors with explanations in ${feedbackLanguage}
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
      "explanation": "Explanation in ${feedbackLanguage}"
    }
  ],
  "alternatives": ["Alternative expression 1", "Alternative expression 2"]
}`

    const systemMessage = `You are an English teacher providing feedback in ${feedbackLanguage} for language learners. Always explain errors and provide guidance in ${feedbackLanguage}.`

    try {
      const response = await client.chat.completions.create({
        model: 'gpt-5-mini',
        messages: [
          {
            role: 'system',
            content: systemMessage,
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
    const feedbackLanguage = localeService.getFullLanguageName()

    const errors = []
    // let correctCount = 0

    // Simple multilingual error messages
    const errorMessages: Record<string, string> = {
      Japanese: '単語が一致しません',
      Chinese: '单词不匹配',
      Korean: '단어가 일치하지 않습니다',
      Spanish: 'Las palabras no coinciden',
      French: 'Les mots ne correspondent pas',
      German: 'Wörter stimmen nicht überein',
      Italian: 'Le parole non corrispondono',
      Portuguese: 'As palavras não correspondem',
      Russian: 'Слова не совпадают',
      English: 'Words do not match',
    }

    for (let i = 0; i < Math.max(refWords.length, userWords.length); i++) {
      const expected = refWords[i] || ''
      const actual = userWords[i] || ''

      if (expected === actual) {
        // correctCount++
      } else if (expected || actual) {
        errors.push({
          expected,
          actual,
          explanation: errorMessages[feedbackLanguage] || errorMessages.English,
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
