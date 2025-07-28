package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current state
func (m Model) View() string {
	if m.state == StateQuitting {
		return "Goodbye!\n"
	}

	switch m.state {
	case StateWelcome:
		return m.viewWelcome()
	case StateGenerating:
		return m.viewGenerating()
	case StatePlaying:
		return m.viewPlaying()
	case StateListening:
		return m.viewListening()
	case StateGrading:
		return m.viewGrading()
	case StateShowingResult:
		return m.viewShowingResult()
	case StateSettings:
		return m.viewSettings()
	default:
		return "Unknown state"
	}
}

// viewWelcome renders the welcome screen
func (m Model) viewWelcome() string {
	// Welcome container
	welcomeBox := containerStyle.Width(60).Align(lipgloss.Center).Render(
		lipgloss.JoinVertical(lipgloss.Center,
			welcomeTitleStyle.Render("✻ Welcome to DictCLI!"),
			"",
			welcomeSubtitleStyle.Render("LLM-powered English dictation practice"),
			welcomeSubtitleStyle.Render("for Japanese learners"),
			"",
			"Press any key to start...",
		),
	)

	// Tips section
	tips := []string{
		"Tips for getting started:",
		"",
		"• Listen carefully to the audio",
		"• Type what you hear",
		"• Get instant feedback in Japanese",
		"• Press 'R' to replay audio anytime",
	}
	tipsSection := tipStyle.Render(strings.Join(tips, "\n"))

	// Combine all elements
	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		welcomeBox,
		"",
		tipsSection,
	)
}

// viewGenerating renders the generating state
func (m Model) viewGenerating() string {
	content := lipgloss.JoinVertical(lipgloss.Left,
		m.renderTitleBar(),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left,
			m.spinner.View(),
			" Generating sentence...",
		),
	)

	return containerStyle.Render(content)
}

// viewPlaying renders the playing state
func (m Model) viewPlaying() string {
	content := lipgloss.JoinVertical(lipgloss.Left,
		m.renderTitleBar(),
		m.renderStatusBar("🔊 Playing... ", m.renderKeyHints([]string{"R:Replay", "S:Settings", "Q:Quit"})),
		"",
		"Listening to audio...",
	)

	return containerStyle.Render(content)
}

// viewListening renders the listening state
func (m Model) viewListening() string {
	var errorMsg string
	if m.error != nil {
		errorMsg = errorMsgStyle.Render(fmt.Sprintf("Error: %v", m.error))
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		m.renderTitleBar(),
		m.renderStatusBar("✍️  Type what you heard", m.renderKeyHints([]string{"Enter:Submit", "R:Replay", "S:Settings", "Q:Quit"})),
		"",
		inputStyle.Render(m.textInput.View()),
		errorMsg,
	)

	return containerStyle.Render(content)
}

// viewGrading renders the grading state
func (m Model) viewGrading() string {
	content := lipgloss.JoinVertical(lipgloss.Left,
		m.renderTitleBar(),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left,
			m.spinner.View(),
			" Grading your answer...",
		),
	)

	return containerStyle.Render(content)
}

// viewShowingResult renders the result state
func (m Model) viewShowingResult() string {
	if m.currentGrade == nil {
		return m.viewListening() // Fallback
	}

	// Score display
	scoreText := fmt.Sprintf("Score: %d%%  WER: %.2f", m.currentGrade.Score, m.currentGrade.WER)
	scoreDisplay := scoreStyle.Render(scoreText)

	// Build result content
	var resultLines []string

	// Add mistakes if any
	if len(m.currentGrade.Mistakes) > 0 {
		mistakeTexts := []string{}
		for _, mistake := range m.currentGrade.Mistakes {
			mistakeText := fmt.Sprintf("%s → %s", 
				mistakeStyle.Render(mistake.Expected),
				mistakeStyle.Render(mistake.Actual))
			mistakeTexts = append(mistakeTexts, mistakeText)
		}
		resultLines = append(resultLines, "誤り: "+strings.Join(mistakeTexts, ", "))
	}

	// Add correct answer
	resultLines = append(resultLines, correctStyle.Render("正解: "+m.currentSentence))

	// Add Japanese explanation
	if m.currentGrade.JapaneseExplanation != "" {
		resultLines = append(resultLines, "", explanationStyle.Render(m.currentGrade.JapaneseExplanation))
	}

	// Add alternatives
	if len(m.currentGrade.AlternativeExpressions) > 0 {
		resultLines = append(resultLines, "")
		for i, alt := range m.currentGrade.AlternativeExpressions {
			resultLines = append(resultLines, alternativeStyle.Render(fmt.Sprintf("別解%d: %s", i+1, alt)))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		m.renderTitleBar(),
		m.renderStatusBar(scoreDisplay, m.renderKeyHints([]string{"N:Next", "R:Replay", "S:Settings", "Q:Quit"})),
		"",
		strings.Join(resultLines, "\n"),
	)

	return containerStyle.Render(content)
}

// viewSettings renders the settings modal
func (m Model) viewSettings() string {
	cfg := m.tempConfig
	if cfg == nil {
		cfg = m.config
	}

	var settingsLines []string
	settingsLines = append(settingsLines, settingsTitleStyle.Render("Settings"))
	settingsLines = append(settingsLines, "")

	// Voice setting
	voiceItem := fmt.Sprintf("Voice  : %-10s (←/→)", cfg.Voice)
	if m.settingsFocus == 0 {
		settingsLines = append(settingsLines, selectedSettingStyle.Render(voiceItem))
	} else {
		settingsLines = append(settingsLines, settingsItemStyle.Render(voiceItem))
	}

	// Level setting
	levelItem := fmt.Sprintf("Level  : TOEIC%-4d (←/→)", cfg.Level)
	if m.settingsFocus == 1 {
		settingsLines = append(settingsLines, selectedSettingStyle.Render(levelItem))
	} else {
		settingsLines = append(settingsLines, settingsItemStyle.Render(levelItem))
	}

	// Topic setting
	topicItem := fmt.Sprintf("Topic  : %-10s (←/→)", cfg.Topic)
	if m.settingsFocus == 2 {
		settingsLines = append(settingsLines, selectedSettingStyle.Render(topicItem))
	} else {
		settingsLines = append(settingsLines, settingsItemStyle.Render(topicItem))
	}

	// Length setting
	lengthItem := fmt.Sprintf("Length : %-2d words   (←/→)", cfg.WordCount)
	if m.settingsFocus == 3 {
		settingsLines = append(settingsLines, selectedSettingStyle.Render(lengthItem))
	} else {
		settingsLines = append(settingsLines, settingsItemStyle.Render(lengthItem))
	}

	// Speed setting
	speedItem := fmt.Sprintf("Speed  : %.1fx        (←/→)", cfg.SpeechSpeed)
	if m.settingsFocus == 4 {
		settingsLines = append(settingsLines, selectedSettingStyle.Render(speedItem))
	} else {
		settingsLines = append(settingsLines, settingsItemStyle.Render(speedItem))
	}

	settingsLines = append(settingsLines, "")
	settingsLines = append(settingsLines, strings.Repeat("─", 40))
	settingsLines = append(settingsLines, keyHintStyle.Render("[Enter] Save & Next Round   [Esc] Cancel"))

	return settingsContainerStyle.Render(strings.Join(settingsLines, "\n"))
}

// Helper methods
func (m Model) renderTitleBar() string {
	title := fmt.Sprintf("DictCLI - Voice:%s | Level:TOEIC%d | Topic:%s | Length:%dw",
		m.config.Voice, m.config.Level, m.config.Topic, m.config.WordCount)
	return titleBarStyle.Width(m.getContentWidth()).Render(title)
}

func (m Model) renderStatusBar(left, right string) string {
	width := m.getContentWidth()
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	padding := width - leftWidth - rightWidth

	if padding < 0 {
		padding = 0
	}

	return statusBarStyle.Width(width).Render(
		left + strings.Repeat(" ", padding) + right,
	)
}

func (m Model) renderKeyHints(hints []string) string {
	var rendered []string
	for _, hint := range hints {
		parts := strings.Split(hint, ":")
		if len(parts) == 2 {
			key := helpKeyStyle.Render(parts[0])
			desc := parts[1]
			rendered = append(rendered, key+":"+desc)
		} else {
			rendered = append(rendered, hint)
		}
	}
	return strings.Join(rendered, "  ")
}

func (m Model) getContentWidth() int {
	if m.width > 0 {
		return m.width - 4 // Account for padding and borders
	}
	return 80 // Default width
}

// Add missing style
var settingsTitleStyle = lipgloss.NewStyle().
	Foreground(primaryColor).
	Bold(true).
	Align(lipgloss.Center)