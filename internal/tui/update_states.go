package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/konpyu/dictcli/internal/types"
)

func (m Model) updateWelcome(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) updateGenerating(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case generatedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateListening
			return m, nil
		}
		m.currentSession = &types.DictationSession{
			ID:        generateSessionID(),
			Timestamp: time.Now(),
			Config:    *m.cfg,
			Sentence:  msg.sentence,
		}
		return m.changeState(StatePlaying)
	}
	return m, nil
}

func (m Model) updatePlaying(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case audioGeneratedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateListening
			return m, nil
		}
		if m.currentSession != nil {
			m.currentSession.AudioPath = msg.audioPath
		}
		return m, m.playAudio
		
	case audioPlayedMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		m.textInput.Reset()
		m.textInput.Focus()
		if m.currentSession != nil {
			m.currentSession.StartTime = time.Now()
		}
		return m.changeState(StateListening)
	}
	return m, nil
}

func (m Model) updateListening(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.userInput = strings.TrimSpace(m.textInput.Value())
			if m.userInput != "" && m.currentSession != nil {
				m.currentSession.UserInput = m.userInput
				m.currentSession.EndTime = time.Now()
				if m.currentSession.StartTime.IsZero() {
					m.currentSession.StartTime = time.Now()
				}
				m.currentSession.DurationSecs = m.currentSession.EndTime.Sub(m.currentSession.StartTime).Seconds()
				return m.changeState(StateGrading)
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) updateGrading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gradedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateShowingResult
			return m, nil
		}
		if m.currentSession != nil {
			m.currentSession.Grade = msg.grade
		}
		return m, tea.Batch(
			m.saveSession,
			func() tea.Msg { return stateChangeMsg{newState: StateShowingResult} },
		)
		
	case sessionSavedMsg:
		if msg.err != nil {
			m.err = msg.err
		}
	}
	return m, nil
}

func (m Model) updateShowingResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.settingsIndex > 0 {
				m.settingsIndex--
			}
		case tea.KeyDown:
			if m.settingsIndex < len(m.settingsFields)-1 {
				m.settingsIndex++
			}
		case tea.KeyLeft:
			m.adjustSetting(-1)
		case tea.KeyRight:
			m.adjustSetting(1)
		case tea.KeyEnter:
			return m.changeState(StateGenerating)
		}
	}
	return m, nil
}

func (m *Model) adjustSetting(delta int) {
	field := m.settingsFields[m.settingsIndex]
	switch field {
	case "voice":
		voices := []string{"alloy", "echo", "fable", "onyx", "nova", "shimmer"}
		currentIdx := indexOf(voices, m.cfg.Voice)
		newIdx := (currentIdx + delta + len(voices)) % len(voices)
		m.cfg.Voice = voices[newIdx]
		
	case "level":
		m.cfg.Level += delta * 50
		if m.cfg.Level < 400 {
			m.cfg.Level = 400
		} else if m.cfg.Level > 990 {
			m.cfg.Level = 990
		}
		
	case "topic":
		topics := []string{"Business", "Travel", "Daily", "Technology", "Health"}
		currentIdx := indexOf(topics, m.cfg.Topic)
		newIdx := (currentIdx + delta + len(topics)) % len(topics)
		m.cfg.Topic = topics[newIdx]
		
	case "words":
		m.cfg.Words += delta
		if m.cfg.Words < 5 {
			m.cfg.Words = 5
		} else if m.cfg.Words > 30 {
			m.cfg.Words = 30
		}
		
	case "speed":
		m.cfg.Speed += float64(delta) * 0.1
		if m.cfg.Speed < 0.5 {
			m.cfg.Speed = 0.5
		} else if m.cfg.Speed > 2.0 {
			m.cfg.Speed = 2.0
		}
	}
}

func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return 0
}