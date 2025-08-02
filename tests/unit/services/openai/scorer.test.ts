import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ScorerService } from '../../../../src/services/openai/scorer.js'
import * as clientModule from '../../../../src/services/openai/client.js'

vi.mock('../../../../src/services/openai/client.js', () => ({
  getOpenAIClient: vi.fn(),
}))

describe('ScorerService', () => {
  let scorer: ScorerService
  let mockClient: any

  beforeEach(() => {
    scorer = new ScorerService()
    mockClient = {
      chat: {
        completions: {
          create: vi.fn(),
        },
      },
    }
    vi.mocked(clientModule.getOpenAIClient).mockReturnValue(mockClient)
  })

  it('should score user input correctly', async () => {
    const mockResponse = {
      score: 85,
      wer: 0.15,
      errors: [
        {
          expected: 'going',
          actual: 'go',
          explanation: '動詞の時制が違います',
        },
      ],
      alternatives: ['I will go to the beach', "I'm heading to the beach"],
    }

    mockClient.chat.completions.create.mockResolvedValue({
      choices: [
        {
          message: {
            content: JSON.stringify(mockResponse),
          },
        },
      ],
    })

    const result = await scorer.scoreAnswer('I am going to the beach', 'I am go to the beach')

    expect(result).toEqual(mockResponse)
  })

  it('should use fallback scoring when API fails', async () => {
    mockClient.chat.completions.create.mockRejectedValue(new Error('API Error'))

    const result = await scorer.scoreAnswer('Hello world', 'Hello word')

    expect(result.score).toBe(50)
    expect(result.wer).toBe(0.5)
    expect(result.errors).toHaveLength(1)
    expect(result.errors[0]).toEqual({
      expected: 'world',
      actual: 'word',
      explanation: '単語が一致しません',
    })
  })

  it('should create round from scoring result', () => {
    const scoringResult = {
      score: 90,
      wer: 0.1,
      errors: [],
      alternatives: [],
    }

    const round = scorer.createRound('Test sentence', 'Test sentence', scoringResult)

    expect(round.sentence).toBe('Test sentence')
    expect(round.userInput).toBe('Test sentence')
    expect(round.score).toBe(90)
    expect(round.wer).toBe(0.1)
    expect(round.id).toBeTruthy()
    expect(round.timestamp).toBeInstanceOf(Date)
  })
})