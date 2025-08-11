import { execSync } from 'child_process'

export type SystemLocale = {
  language: string // e.g., 'ja', 'en', 'zh'
  country?: string // e.g., 'JP', 'US', 'CN'
  full: string // e.g., 'ja_JP.UTF-8'
}

export class LocaleService {
  private static instance: LocaleService
  private locale: SystemLocale | null = null

  static getInstance(): LocaleService {
    if (!LocaleService.instance) {
      LocaleService.instance = new LocaleService()
    }
    return LocaleService.instance
  }

  getSystemLocale(): SystemLocale {
    if (this.locale) {
      return this.locale
    }

    // Try environment variable first (most reliable)
    const lang = process.env.LANG || process.env.LC_ALL || process.env.LC_MESSAGES || ''

    if (lang) {
      this.locale = this.parseLocaleString(lang)
      return this.locale
    }

    // Fallback to macOS defaults
    try {
      const appleLanguages = execSync('defaults read -g AppleLanguages', { encoding: 'utf8' })
      const match = appleLanguages.match(/"([a-z]{2})(-([A-Z]{2}))?"/i)
      if (match) {
        this.locale = {
          language: match[1].toLowerCase(),
          country: match[3],
          full: `${match[1].toLowerCase()}_${match[3] || 'XX'}.UTF-8`,
        }
        return this.locale
      }
    } catch {
      // Ignore errors and use fallback
    }

    // Default to English
    this.locale = {
      language: 'en',
      country: 'US',
      full: 'en_US.UTF-8',
    }
    return this.locale
  }

  private parseLocaleString(locale: string): SystemLocale {
    // Parse strings like 'ja_JP.UTF-8' or 'en_US.UTF-8'
    const match = locale.match(/^([a-z]{2})(?:_([A-Z]{2}))?/i)
    if (match) {
      return {
        language: match[1].toLowerCase(),
        country: match[2],
        full: locale,
      }
    }

    // Fallback
    return {
      language: 'en',
      country: 'US',
      full: locale,
    }
  }

  isJapanese(): boolean {
    const locale = this.getSystemLocale()
    return locale.language === 'ja'
  }

  getUILanguage(): 'ja' | 'en' {
    // For now, support Japanese and English
    // Can be extended to support more languages
    return this.isJapanese() ? 'ja' : 'en'
  }

  // Get the full language name for OpenAI prompts
  getFullLanguageName(): string {
    const locale = this.getSystemLocale()
    const languageMap: Record<string, string> = {
      ja: 'Japanese',
      en: 'English',
      zh: 'Chinese',
      ko: 'Korean',
      es: 'Spanish',
      fr: 'French',
      de: 'German',
      it: 'Italian',
      pt: 'Portuguese',
      ru: 'Russian',
      ar: 'Arabic',
      hi: 'Hindi',
      th: 'Thai',
      vi: 'Vietnamese',
      nl: 'Dutch',
      pl: 'Polish',
      tr: 'Turkish',
      id: 'Indonesian',
      sv: 'Swedish',
      da: 'Danish',
      no: 'Norwegian',
      fi: 'Finnish',
    }
    return languageMap[locale.language] || 'English'
  }
}

export const localeService = LocaleService.getInstance()
