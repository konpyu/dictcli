import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ProblemGenerator } from '../../../../src/services/openai/generator.js'
import * as clientModule from '../../../../src/services/openai/client.js'

vi.mock('../../../../src/services/openai/client.js', () => ({
  getOpenAIClient: vi.fn(),
}))

describe('ProblemGenerator', () => {
  let generator: ProblemGenerator
  let mockClient: any

  beforeEach(() => {
    generator = new ProblemGenerator()
    mockClient = {
      chat: {
        completions: {
          create: vi.fn(),
        },
      },
    }
    vi.mocked(clientModule.getOpenAIClient).mockReturnValue(mockClient)
  })

  it('should generate a sentence with correct parameters', async () => {
    const mockSentence = 'The meeting starts at nine.'
    mockClient.chat.completions.create.mockResolvedValue({
      choices: [
        {
          message: {
            content: mockSentence,
          },
        },
      ],
    })

    const result = await generator.generateProblem('CEFR_A1', 'EverydayLife', 5)

    expect(result).toBe(mockSentence)
    expect(mockClient.chat.completions.create).toHaveBeenCalledWith({
      model: 'gpt-4o-mini',
      messages: expect.any(Array),
      temperature: 0.9,
      max_tokens: 100,
      presence_penalty: 0.6,
      frequency_penalty: 0.6,
    })
  })

  it('should throw error when generation fails', async () => {
    mockClient.chat.completions.create.mockRejectedValue(new Error('API Error'))

    await expect(generator.generateProblem('CEFR_A1', 'EverydayLife', 5)).rejects.toThrow(
      'Failed to generate problem',
    )
  })
})