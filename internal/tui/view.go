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
• Press 'R' to replay audio anytime`)

	return welcome + tips
}

func (m Model) viewGenerating() string {
	return boxStyle.Width(60).Render(
		m.header() + "\n\n" +
			fmt.Sprintf("%s Generating sentence...", m.spinner.View()),
	)
}

func (m Model) viewPlaying() string {
	status := fmt.Sprintf("%s Playing audio... (%.1fx speed)", m.spinner.View(), m.cfg.Speed)
	if m.currentSession != nil && m.currentSession.ReplayCount > 1 {
		status += fmt.Sprintf(" [Replay #%d]", m.currentSession.ReplayCount)
	}

	return boxStyle.Width(60).Render(
		m.header() + "\n\n" + status,
	)
}

func (m Model) viewListening() string {
	var content strings.Builder
	content.WriteString(m.header())
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("(🔊 [R]eplay  [S]ettings  [Q]uit)"))
	content.WriteString("\n\n")

	if m.err != nil {
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		content.WriteString("\n\n")
	}

	content.WriteString(m.textInput.View())

	return boxStyle.Width(80).Render(content.String())
}

func (m Model) viewGrading() string {
	return boxStyle.Width(60).Render(
		m.header() + "\n\n" +
			fmt.Sprintf("%s Grading your input...", m.spinner.View()),
	)
}

func (m Model) viewShowingResult() string {
	var content strings.Builder
	content.WriteString(m.header())
	content.WriteString("\n")

	if m.currentSession != nil && m.currentSession.Grade != nil {
		grade := m.currentSession.Grade
		scoreColor := "196"
		if grade.Score >= 80 {
			scoreColor = "46"
		} else if grade.Score >= 60 {
			scoreColor = "226"
		}

		content.WriteString(
			fmt.Sprintf("Score: %s  WER: %.2f",
				lipgloss.NewStyle().Foreground(lipgloss.Color(scoreColor)).Render(fmt.Sprintf("%d%%", grade.Score)),
				grade.WER,
			),
		)
		content.WriteString("  ")
		content.WriteString(mutedStyle.Render("[N]ext  [R]eplay  [S]ettings  [Q]uit"))
		content.WriteString("\n\n")

		if m.err != nil {
			content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
			content.WriteString("\n")
		}

		if len(grade.Mistakes) > 0 {
			content.WriteString("誤り: ")
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
		}

		content.WriteString(fmt.Sprintf("正解: %s\n", successStyle.Render(m.currentSession.Sentence)))

		if grade.JapaneseExplanation != "" {
			content.WriteString("\n")
			content.WriteString(grade.JapaneseExplanation)
			content.WriteString("\n")
		}

		if len(grade.AlternativeExpressions) > 0 {
			content.WriteString("\n")
			for i, alt := range grade.AlternativeExpressions {
				content.WriteString(fmt.Sprintf("別解%d: %s\n", i+1, alt))
			}
		}
	} else {
		content.WriteString(errorStyle.Render("No grading result available"))
	}

	return boxStyle.Width(80).Render(content.String())
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
	content.WriteString(mutedStyle.Render("[Enter] Save & Next Round   [Esc] Cancel"))

	return boxStyle.Width(50).Render(content.String())
}