import OpenAI from 'openai'

export class OpenAIClient {
  private client: OpenAI

  constructor() {
    const apiKey = process.env.OPENAI_API_KEY
    if (!apiKey) {
      throw new Error('OPENAI_API_KEY environment variable is required')
    }

    this.client = new OpenAI({
      apiKey,
    })
  }

  getClient(): OpenAI {
    return this.client
  }
}

let instance: OpenAIClient | null = null

export const getOpenAIClient = (): OpenAI => {
  if (!instance) {
    instance = new OpenAIClient()
  }
  return instance.getClient()
}
