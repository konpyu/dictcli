package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/konpyu/dictcli/internal/logging"
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
	case clearMessageMsg:
		m.message = ""
		return m, nil
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
			logging.Debug("Enter pressed in settings, current config - Voice: %s, Level: %d, Topic: %s, Words: %d, Speed: %.1f", 
				m.cfg.Voice, m.cfg.Level, m.cfg.Topic, m.cfg.Words, m.cfg.Speed)
			if err := m.saveConfig(); err != nil {
				m.err = err
				return m, nil
			}
			return m.changeState(StateGenerating)
		case tea.KeyCtrlS:
			logging.Info("Ctrl+S pressed - saving settings")
			if err := m.saveConfig(); err != nil {
				m.err = err
				return m, nil
			}
			m.message = "Settings saved successfully!"
			return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
				return clearMessageMsg{}
			})
		}
	}
	return m, nil
}

func (m *Model) adjustSetting(delta int) {
	field := m.settingsFields[m.settingsIndex]
	logging.Debug("Adjusting setting %s (delta: %d)", field, delta)
	
	switch field {
	case "voice":
		voices := []string{"alloy", "echo", "fable", "onyx", "nova", "shimmer"}
		currentIdx := indexOf(voices, m.cfg.Voice)
		newIdx := (currentIdx + delta + len(voices)) % len(voices)
		oldVoice := m.cfg.Voice
		m.cfg.Voice = voices[newIdx]
		logging.Debug("Voice changed: %s -> %s", oldVoice, m.cfg.Voice)
		
	case "level":
		oldLevel := m.cfg.Level
		m.cfg.Level += delta * 50
		if m.cfg.Level < 400 {
			m.cfg.Level = 400
		} else if m.cfg.Level > 990 {
			m.cfg.Level = 990
		}
		logging.Debug("Level changed: %d -> %d", oldLevel, m.cfg.Level)
		
	case "topic":
		topics := []string{"Business", "Travel", "Daily", "Technology", "Health"}
		currentIdx := indexOf(topics, m.cfg.Topic)
		newIdx := (currentIdx + delta + len(topics)) % len(topics)
		oldTopic := m.cfg.Topic
		m.cfg.Topic = topics[newIdx]
		logging.Debug("Topic changed: %s -> %s", oldTopic, m.cfg.Topic)
		
	case "words":
		oldWords := m.cfg.Words
		m.cfg.Words += delta
		if m.cfg.Words < 5 {
			m.cfg.Words = 5
		} else if m.cfg.Words > 30 {
			m.cfg.Words = 30
		}
		logging.Debug("Words changed: %d -> %d", oldWords, m.cfg.Words)
		
	case "speed":
		oldSpeed := m.cfg.Speed
		m.cfg.Speed += float64(delta) * 0.1
		if m.cfg.Speed < 0.5 {
			m.cfg.Speed = 0.5
		} else if m.cfg.Speed > 2.0 {
			m.cfg.Speed = 2.0
		}
		logging.Debug("Speed changed: %.1f -> %.1f", oldSpeed, m.cfg.Speed)
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

func (m *Model) saveConfig() error {
	// Debug: Print config before saving
	logging.Debug("Saving config - Voice: %s, Level: %d, Topic: %s, Words: %d, Speed: %.1f", 
		m.cfg.Voice, m.cfg.Level, m.cfg.Topic, m.cfg.Words, m.cfg.Speed)
	
	if err := m.configManager.Set("voice", m.cfg.Voice); err != nil {
		return err
	}
	if err := m.configManager.Set("level", m.cfg.Level); err != nil {
		return err
	}
	if err := m.configManager.Set("topic", m.cfg.Topic); err != nil {
		return err
	}
	if err := m.configManager.Set("words", m.cfg.Words); err != nil {
		return err
	}
	if err := m.configManager.Set("speed", m.cfg.Speed); err != nil {
		return err
	}
	
	if err := m.configManager.Save(); err != nil {
		return err
	}
	
	// CRITICAL FIX: Update the model's config pointer to reflect the saved values
	m.cfg = m.configManager.Get()
	
	// Debug: Print config after saving
	logging.Debug("Config after save - Voice: %s, Level: %d, Topic: %s, Words: %d, Speed: %.1f", 
		m.cfg.Voice, m.cfg.Level, m.cfg.Topic, m.cfg.Words, m.cfg.Speed)
	
	return nil
}

func (m Model) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}