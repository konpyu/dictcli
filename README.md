# DictCLI

[![npm version](https://badge.fury.io/js/dictcli.svg)](https://badge.fury.io/js/dictcli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**DictCLI** は、日本人英語学習者向けのLLM駆動型ディクテーション練習ツールです。ターミナル上で動作し、OpenAI APIを使用して英文の生成、音声再生、日本語での採点フィードバックを提供します。

<p align="center">
  <img src="https://github.com/yourusername/dictcli/assets/demo.gif" alt="DictCLI Demo" width="600">
</p>

## ✨ 特徴

- 🤖 **LLM駆動**: OpenAI APIによる動的な問題生成
- 🎯 **レベル別学習**: CEFR A1〜C2の6段階から選択可能
- 🔊 **音声再生**: 6種類の音声（男性3・女性3）から選択
- 📝 **日本語フィードバック**: 間違いを日本語で詳しく解説
- ⚡ **高速な学習サイクル**: 1分未満でラウンドを完了
- 🎨 **美しいTUI**: Ink v4による洗練されたターミナルUI

## 📋 必要要件

- Node.js v20以上
- macOS（音声再生機能のため）
- OpenAI APIキー

## 🚀 インストール

### グローバルインストール（推奨）

```bash
npm install -g dictcli
```

### ローカルインストール

```bash
git clone https://github.com/yourusername/dictcli.git
cd dictcli
npm install
npm run build
npm link
```

インストール後、ターミナルで以下のように実行できます：

```bash
dictcli
```

アンインストールする場合：

```bash
npm unlink -g dictcli
```

## 🔧 セットアップ

### 1. OpenAI APIキーの設定

OpenAI APIキーを環境変数に設定してください：

```bash
export OPENAI_API_KEY=your-api-key-here
```

永続的に設定する場合は、シェルの設定ファイル（`.bashrc`、`.zshrc`など）に追加してください。

### 2. APIキーの取得方法

1. [OpenAI Platform](https://platform.openai.com/)にアクセス
2. アカウントを作成またはログイン
3. [API Keys](https://platform.openai.com/api-keys)ページでキーを生成
4. 生成されたキーをコピーして環境変数に設定

## 📖 使い方

### 基本的な使い方

```bash
dictcli
```

デフォルト設定で起動します（ビジネス英語、A1レベル、10単語）。

### オプション指定

```bash
dictcli --topic Technology --level CEFR_B1 --words 15 --voice SARA
```

### 利用可能なオプション

| オプション | 説明 | 選択肢 | デフォルト |
|-----------|------|--------|------------|
| `--topic` | 学習トピック | EverydayLife, Travel, Technology, Health, Entertainment, Business, Random | Business |
| `--level` | 難易度レベル | CEFR_A1, CEFR_A2, CEFR_B1, CEFR_B2, CEFR_C1, CEFR_C2 | CEFR_A1 |
| `--words` | 単語数 | 5-30 | 6 |
| `--voice` | 音声の種類 | ALEX, SARA, EVAN, NOVA, NICK, FAYE | ALEX |

### 音声の説明

| 名前 | 性別 | 特徴 |
|------|------|------|
| ALEX | 男性 | 親しみやすい標準的な声 |
| SARA | 女性 | クリアで聞き取りやすい声 |
| EVAN | 男性 | 落ち着いた声 |
| NOVA | 女性 | モダンな声 |
| NICK | 男性 | 深みのある声 |
| FAYE | 女性 | 優しい声 |

## 🎮 操作方法

### 学習画面での操作

- **英文入力**: 聞き取った英文をタイプして Enter
- **スラッシュコマンド**: `/` を入力してコマンドメニューを表示
  - `/replay` - 音声を再生
  - `/settings` - 設定画面を開く
  - `/giveup` - ヒント（空欄付き答え）を表示
  - `/quit` - アプリを終了

### 結果画面での操作

- **Enter** または **N** - 次の問題へ
- **R** - 音声を再生
- **S** - 設定画面を開く
- **Q** - アプリを終了

### 設定画面での操作

- **↑↓** - 項目を選択
- **←→** - 値を変更（Voice, Level, Topic）
- **-+** - 値を増減（Word Count）
- **Enter** - 保存して次の問題へ
- **Esc** - キャンセル

## 💡 学習のヒント

### 効果的な学習方法

1. **レベル選択**: 80%程度正解できるレベルから始める
2. **繰り返し練習**: 同じトピックで繰り返し練習して専門用語に慣れる
3. **速度調整**: `/settings`で再生速度を調整（0.8x〜1.2x）
4. **エラー分析**: 日本語の解説をしっかり読んで間違いのパターンを理解

### レベルの目安

- **A1**: 基本的な単語と簡単な文
- **A2**: 日常的な表現と基本的な文法
- **B1**: 仕事や旅行で使う実用的な表現
- **B2**: より複雑な文法と語彙
- **C1**: 流暢で自然な表現
- **C2**: ネイティブレベルの高度な表現

## 🐛 トラブルシューティング

### よくある問題

**Q: 音声が再生されない**
- A: macOSでのみ動作します。他のOSでは音声機能は利用できません。

**Q: "OPENAI_API_KEY environment variable is required" エラー**
- A: 環境変数が設定されているか確認してください：
  ```bash
  echo $OPENAI_API_KEY
  ```

**Q: APIエラーが頻発する**
- A: OpenAIのAPIクレジットが残っているか確認してください。

**Q: 文字化けする**
- A: ターミナルがUTF-8をサポートしているか確認してください。

### デバッグモード

問題が解決しない場合は、デバッグモードで詳細情報を確認できます：

```bash
DICTCLI_DEBUG=true dictcli
```

## 🤝 コントリビューション

バグ報告や機能提案は [GitHub Issues](https://github.com/yourusername/dictcli/issues) からお願いします。

プルリクエストも歓迎です！

## 📄 ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 🙏 謝辞

- [Ink](https://github.com/vadimdemedes/ink) - 美しいCLI UIフレームワーク
- [OpenAI](https://openai.com/) - 強力なAI API
- [Zustand](https://github.com/pmndrs/zustand) - シンプルな状態管理

---

<p align="center">
  Made with ❤️ for Japanese English Learners
</p>
