import { localeService } from '../services/locale.js'

type TranslationKey =
  | 'welcome'
  | 'selectTopic'
  | 'selectLevel'
  | 'listening'
  | 'typeWhatYouHeard'
  | 'score'
  | 'nextQuestion'
  | 'quit'
  | 'settings'
  | 'replay'
  | 'giveUp'
  | 'errorOccurred'
  | 'apiKeyMissing'

const translations: Record<'ja' | 'en', Record<TranslationKey, string>> = {
  ja: {
    welcome: 'ディクテーション学習へようこそ',
    selectTopic: 'トピックを選択',
    selectLevel: 'レベルを選択',
    listening: '聞いています...',
    typeWhatYouHeard: '聞いた内容を入力してください',
    score: 'スコア',
    nextQuestion: '次の問題',
    quit: '終了',
    settings: '設定',
    replay: 'もう一度聞く',
    giveUp: 'ギブアップ',
    errorOccurred: 'エラーが発生しました',
    apiKeyMissing: 'APIキーが設定されていません',
  },
  en: {
    welcome: 'Welcome to Dictation Learning',
    selectTopic: 'Select Topic',
    selectLevel: 'Select Level',
    listening: 'Listening...',
    typeWhatYouHeard: 'Type what you heard',
    score: 'Score',
    nextQuestion: 'Next Question',
    quit: 'Quit',
    settings: 'Settings',
    replay: 'Replay',
    giveUp: 'Give Up',
    errorOccurred: 'An error occurred',
    apiKeyMissing: 'API key is not configured',
  },
}

export function t(key: TranslationKey): string {
  const lang = localeService.getUILanguage()
  return translations[lang][key] || translations.en[key]
}

export function getLanguage(): 'ja' | 'en' {
  return localeService.getUILanguage()
}
