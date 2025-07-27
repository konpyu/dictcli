package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Width(80)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226"))
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var content string

	switch m.state {
	case StateWelcome:
		content = m.viewWelcome()
	case StateGenerating:
		content = m.viewGenerating()
	case StatePlaying:
		content = m.viewPlaying()
	case StateListening:
		content = m.viewListening()
	case StateGrading:
		content = m.viewGrading()
	case StateShowingResult:
		content = m.viewShowingResult()
	case StateSettings:
		content = m.viewSettings()
	case StateHelp:
		content = m.viewHelp()
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m Model) header() string {
	header := fmt.Sprintf("DictCLI - Voice:%s | Level:TOEIC%d | Topic:%s | Length:%dw",
		m.cfg.Voice, m.cfg.Level, m.cfg.Topic, m.cfg.Words)
	return headerStyle.Render(header)
}

func (m Model) viewWelcome() string {
	welcome := boxStyle.Render(
		titleStyle.Render("✻ Welcome to DictCLI!") + "\n\n" +
			"LLM-powered English dictation practice\n" +
			"for Japanese learners\n\n" +
			mutedStyle.Render("Press any key to start..."),
	)

	tips := mutedStyle.Render(`
Tips for getting started:

• Listen carefully to the audio
• Type what you hear
• Get instant feedback in Japanese
• Press 'Ctrl+R' to replay audio anytime
• Press 'Ctrl+H' or '?' for help`)

	return welcome + tips
}

func (m Model) viewGenerating() string {
	var content strings.Builder
	content.WriteString(m.header())
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("%s 文章を生成中...", m.spinner.View()))
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("OpenAI GPT-4o-miniで新しい練習文を作成しています"))
	
	if m.err != nil {
		content.WriteString("\n\n")
		content.WriteString(errorStyle.Render(fmt.Sprintf("エラー: %v", m.err)))
		content.WriteString("\n")
		content.WriteString(mutedStyle.Render("[Enter] 再試行  [Ctrl+Q] 終了"))
	}
	
	return boxStyle.Width(70).Render(content.String())
}

func (m Model) viewPlaying() string {
	var content strings.Builder
	content.WriteString(m.header())
	content.WriteString("\n\n")
	
	status := fmt.Sprintf("%s 音声再生中... (%.1fx速度)", m.spinner.View(), m.cfg.Speed)
	if m.currentSession != nil && m.currentSession.ReplayCount > 1 {
		status += fmt.Sprintf(" [%d回目]", m.currentSession.ReplayCount)
	}
	content.WriteString(status)
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("OpenAI TTS-1で音声を生成・再生しています"))
	
	if m.err != nil {
		content.WriteString("\n\n")
		content.WriteString(errorStyle.Render(fmt.Sprintf("エラー: %v", m.err)))
		content.WriteString("\n")
		content.WriteString(mutedStyle.Render("[Enter] 再試行  [Ctrl+Q] 終了"))
	}
	
	return boxStyle.Width(70).Render(content.String())
}

func (m Model) viewListening() string {
	var content strings.Builder
	content.WriteString(m.header())
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("(🔊 [Ctrl+R]eplay  [Ctrl+S]ettings  [Ctrl+H]elp  [Ctrl+Q]uit)"))
	content.WriteString("\n\n")

	if m.err != nil {
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		content.WriteString("\n\n")
	}

	content.WriteString(m.textInput.View())

	return boxStyle.Width(80).Render(content.String())
}

func (m Model) viewGrading() string {
	var content strings.Builder
	content.WriteString(m.header())
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("%s 採点中...", m.spinner.View()))
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("OpenAI GPT-4o-miniで詳細な採点と解説を生成しています"))
	
	if m.err != nil {
		content.WriteString("\n\n")
		content.WriteString(errorStyle.Render(fmt.Sprintf("エラー: %v", m.err)))
		content.WriteString("\n")
		content.WriteString(mutedStyle.Render("[Enter] 再試行  [Ctrl+Q] 終了"))
	}
	
	return boxStyle.Width(70).Render(content.String())
}

func (m Model) viewShowingResult() string {
	var content strings.Builder
	content.WriteString(m.header())
	content.WriteString("\n")

	if m.currentSession != nil && m.currentSession.Grade != nil {
		grade := m.currentSession.Grade
		scoreColor := "196"
		scoreEmoji := "😞"
		if grade.Score >= 80 {
			scoreColor = "46"
			scoreEmoji = "🎉"
		} else if grade.Score >= 60 {
			scoreColor = "226"
			scoreEmoji = "👍"
		}

		// Success animation for high scores
		if grade.Score >= 90 {
			content.WriteString(successStyle.Render("🌟 素晴らしい結果です！ 🌟"))
			content.WriteString("\n")
		} else if grade.Score >= 80 {
			content.WriteString(successStyle.Render("✨ よくできました！ ✨"))
			content.WriteString("\n")
		}

		content.WriteString(
			fmt.Sprintf("%s スコア: %s  WER: %.2f",
				scoreEmoji,
				lipgloss.NewStyle().Foreground(lipgloss.Color(scoreColor)).Bold(true).Render(fmt.Sprintf("%d%%", grade.Score)),
				grade.WER,
			),
		)
		content.WriteString("  ")
		content.WriteString(mutedStyle.Render("[N]ext  [R]eplay  [S]ettings  [Q]uit"))
		content.WriteString("\n\n")

		if m.err != nil {
			content.WriteString(errorStyle.Render(fmt.Sprintf("⚠️ エラー: %v", m.err)))
			content.WriteString("\n")
		}

		if len(grade.Mistakes) > 0 {
			content.WriteString("📝 間違い: ")
			for i, mistake := range grade.Mistakes {
				if i > 0 {
					content.WriteString(", ")
				}
				content.WriteString(fmt.Sprintf("%s → %s", 
					errorStyle.Render(mistake.Expected),
					highlightStyle.Render(mistake.Actual),
				))
			}
			content.WriteString("\n")
		} else {
			content.WriteString(successStyle.Render("✓ 完璧です！間違いはありませんでした"))
			content.WriteString("\n")
		}

		content.WriteString(fmt.Sprintf("📝 あなたの入力: %s\n", m.currentSession.UserInput))
		content.WriteString(fmt.Sprintf("✅ 正解: %s\n", successStyle.Render(m.currentSession.Sentence)))

		if grade.JapaneseExplanation != "" {
			content.WriteString("\n💡 解説:\n")
			content.WriteString(grade.JapaneseExplanation)
			content.WriteString("\n")
		}

		if len(grade.AlternativeExpressions) > 0 {
			content.WriteString("\n🔄 別の表現:\n")
			for i, alt := range grade.AlternativeExpressions {
				content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, alt))
			}
		}

		// Add session stats
		if m.currentSession.DurationSecs > 0 {
			content.WriteString("\n")
			content.WriteString(mutedStyle.Render(fmt.Sprintf("⏱️ 入力時間: %.1f秒", m.currentSession.DurationSecs)))
		}
	} else {
		content.WriteString(errorStyle.Render("❌ 採点結果が利用できません"))
	}

	return boxStyle.Width(85).Render(content.String())
}

func (m Model) viewSettings() string {
	var content strings.Builder
	content.WriteString(titleStyle.Render("Settings"))
	content.WriteString("\n\n")

	settings := []struct {
		name  string
		value string
		hint  string
	}{
		{"Voice", m.cfg.Voice, "(←/→) ♂/♀ voices"},
		{"Level", fmt.Sprintf("TOEIC%d", m.cfg.Level), "(←/→) 400-990"},
		{"Topic", m.cfg.Topic, "(←/→) 選択可能"},
		{"Length", fmt.Sprintf("%d words", m.cfg.Words), "(−/＋) 5-30"},
		{"Speed", fmt.Sprintf("%.1fx", m.cfg.Speed), "(←/→) 0.5-2.0"},
	}

	for i, setting := range settings {
		line := fmt.Sprintf("%-8s: %-12s %s", setting.name, setting.value, mutedStyle.Render(setting.hint))
		if i == m.settingsIndex {
			line = highlightStyle.Render("▶ ") + line
		} else {
			line = "  " + line
		}
		content.WriteString(line)
		content.WriteString("\n")
	}

	content.WriteString("\n")
	
	if m.message != "" {
		content.WriteString(successStyle.Render(m.message))
		content.WriteString("\n\n")
	}
	
	content.WriteString(mutedStyle.Render("[Enter] Save & Next Round   [Ctrl+S] Save   [Esc] Cancel"))

	return boxStyle.Width(50).Render(content.String())
}

func (m Model) viewHelp() string {
	var content strings.Builder
	content.WriteString(titleStyle.Render("📚 キーボードショートカット"))
	content.WriteString("\n\n")

	// Global shortcuts
	content.WriteString(headerStyle.Render("🌐 全体操作"))
	content.WriteString("\n")
	content.WriteString("  Ctrl+Q     → アプリケーション終了\n")
	content.WriteString("  Ctrl+H/?   → このヘルプ表示\n")
	content.WriteString("  Esc        → 前の状態に戻る\n")
	content.WriteString("\n")

	// State-specific shortcuts
	switch m.state {
	case StateWelcome:
		content.WriteString(headerStyle.Render("🏠 ウェルカム画面"))
		content.WriteString("\n")
		content.WriteString("  任意のキー  → 練習開始\n")
	case StateListening:
		content.WriteString(headerStyle.Render("🎧 入力画面"))
		content.WriteString("\n")
		content.WriteString("  Enter      → 入力内容を提出\n")
		content.WriteString("  Ctrl+R     → 音声を再生\n")
		content.WriteString("  Ctrl+S     → 設定画面を開く\n")
	case StateShowingResult:
		content.WriteString(headerStyle.Render("📊 結果画面"))
		content.WriteString("\n")
		content.WriteString("  N          → 次の問題\n")
		content.WriteString("  R          → 同じ文を再度練習\n")
		content.WriteString("  S          → 設定画面を開く\n")
	case StateSettings:
		content.WriteString(headerStyle.Render("⚙️ 設定画面"))
		content.WriteString("\n")
		content.WriteString("  ↑/↓        → 項目選択\n")
		content.WriteString("  ←/→        → 値調整\n")
		content.WriteString("  Enter      → 設定保存して練習開始\n")
	default:
		content.WriteString(headerStyle.Render("🔄 現在の状態"))
		content.WriteString("\n")
		content.WriteString("しばらくお待ちください...\n")
	}

	content.WriteString("\n")
	content.WriteString(headerStyle.Render("🔊 音声操作"))
	content.WriteString("\n")
	content.WriteString("🔊 音声は自動的に再生されます\n")
	content.WriteString("🔁 Ctrl+Rでいつでも再生可能\n")
	content.WriteString("🎧 ヘッドフォン推奨\n")

	content.WriteString("\n")
	content.WriteString(headerStyle.Render("📝 コツ＆ティップス"))
	content.WriteString("\n")
	content.WriteString("• 音声を止めてから入力すると集中しやすい\n")
	content.WriteString("• 大文字・小文字は区別されます\n")
	content.WriteString("• 句読点も正確に入力してください\n")
	content.WriteString("• TOEICレベルを上げると難しくなります\n")

	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("[Esc] 戻る  [Ctrl+Q] 終了"))

	return boxStyle.Width(75).Render(content.String())
}