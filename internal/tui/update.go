package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/dictcli/internal/logger"
	"github.com/yourusername/dictcli/internal/types"
	"github.com/google/uuid"
)

// Key constants
const (
	keyEnter = "enter"
	keyEsc   = "esc"
	keyUp    = "up"
	keyDown  = "down"
	keyLeft  = "left"
	keyRight = "right"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Log state transitions
	oldState := m.state
	defer func() {
		if oldState != m.state {
			logger.Info("State transition: %s -> %s", m.getCurrentStateString(), m.getCurrentStateString())
		}
	}()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		logger.Debug("Window resized: %dx%d", m.width, m.height)
		return m, nil

	case tea.KeyMsg:
		// Global keyboard shortcuts
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			logger.Info("User requested quit (Ctrl+C)")
			m.state = StateQuitting
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "Q"))):
			if m.state != StateListening && m.state != StateSettings {
				logger.Info("User requested quit (Q)")
				m.state = StateQuitting
				return m, tea.Quit
			}
		}

		// State-specific keyboard handling
		switch m.state {
		case StateWelcome:
			m, cmd = m.updateWelcome(msg)
		case StateListening:
			m, cmd = m.updateListening(msg)
		case StateShowingResult:
			m, cmd = m.updateShowingResult(msg)
		case StateSettings:
			m, cmd = m.updateSettings(msg)
		default:
			// For other states, check for replay and settings shortcuts
			m, cmd = m.handleGlobalShortcuts(msg)
		}
		return m, cmd

	// Handle custom messages
	case GenerateSentenceMsg:
		return m.handleGenerateSentenceMsg(msg)

	case GenerateAudioMsg:
		return m.handleGenerateAudioMsg(msg)

	case PlayAudioMsg:
		return m.handlePlayAudioMsg(msg)

	case GradeDictationMsg:
		return m.handleGradeDictationMsg(msg)

	case SettingsSavedMsg:
		return m.handleSettingsSavedMsg()

	case ErrorMsg:
		m.error = msg.Err
		logger.Error("Error occurred: %v", msg.Err)
		// Show error and return to listening state
		m.state = StateListening
		return m, nil

	case TickMsg:
		// Handle any periodic updates
		return m, nil

	case spinner.TickMsg:
		if m.state == StateGenerating || m.state == StateGrading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update sub-components
	switch m.state {
	case StateListening:
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// updateWelcome handles updates in the welcome state
func (m Model) updateWelcome(msg tea.KeyMsg) (Model, tea.Cmd) {
	logger.Info("User pressed key in welcome screen, starting dictation")
	m.welcomeShown = true
	return m.startNewRound()
}

// updateListening handles updates in the listening state
func (m Model) updateListening(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case keyEnter:
		userInput := m.textInput.Value()
		if userInput == "" {
			logger.Debug("User submitted empty input")
			return m, nil
		}
		logger.Info("User submitted input: %s", userInput)
		m.state = StateGrading
		m.currentSession.UserInput = userInput
		// Store submission time in session
		
		// Start grading
		return m, m.gradeDictation()

	case "r", "R":
		logger.Info("User requested replay")
		m.replayCount++
		m.currentSession.ReplayCount = m.replayCount
		return m, m.playAudio()

	case "s", "S":
		logger.Info("User opened settings")
		m.openSettings()
		return m, nil

	default:
		// Let textinput handle other keys
		return m, nil
	}
}

// updateShowingResult handles updates in the result state
func (m Model) updateShowingResult(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "n", "N", keyEnter:
		logger.Info("User requested next round")
		return m.startNewRound()

	case "r", "R":
		logger.Info("User requested replay in result view")
		return m, m.playAudio()

	case "s", "S":
		logger.Info("User opened settings from result view")
		m.openSettings()
		return m, nil
	}
	return m, nil
}

// updateSettings handles updates in the settings state
func (m Model) updateSettings(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.settingsFocus > 0 {
			m.settingsFocus--
		}
		return m, nil

	case "down", "j":
		if m.settingsFocus < len(m.settingsOptions)-1 {
			m.settingsFocus++
		}
		return m, nil

	case "left", "h":
		m.adjustSetting(-1)
		return m, nil

	case "right", "l":
		m.adjustSetting(1)
		return m, nil

	case keyEnter:
		logger.Info("User saved settings")
		m.saveSettings()
		return m, m.saveSettingsCmd()

	case keyEsc:
		logger.Info("User cancelled settings")
		m.cancelSettings()
		return m, nil
	}
	return m, nil
}

// handleGlobalShortcuts handles shortcuts available in most states
func (m Model) handleGlobalShortcuts(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "r", "R":
		if m.currentAudioPath != "" && (m.state == StatePlaying || m.state == StateListening) {
			logger.Info("User requested replay")
			m.replayCount++
			if m.currentSession != nil {
				m.currentSession.ReplayCount = m.replayCount
			}
			return m, m.playAudio()
		}

	case "s", "S":
		if m.state != StateSettings && m.state != StateWelcome {
			logger.Info("User opened settings")
			m.openSettings()
			return m, nil
		}
	}
	return m, nil
}

// Message handlers
func (m Model) handleGenerateSentenceMsg(msg GenerateSentenceMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.error = fmt.Errorf("failed to generate sentence: %w", msg.Err)
		m.state = StateListening
		return m, nil
	}

	logger.Info("Generated sentence: %s", msg.Sentence)
	m.currentSentence = msg.Sentence
	m.currentSession.Sentence = msg.Sentence
	
	// Generate audio
	return m, m.generateAudio()
}

func (m Model) handleGenerateAudioMsg(msg GenerateAudioMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.error = fmt.Errorf("failed to generate audio: %w", msg.Err)
		m.state = StateListening
		return m, nil
	}

	logger.Info("Generated audio: %s", msg.AudioPath)
	m.currentAudioPath = msg.AudioPath
	m.currentSession.AudioPath = msg.AudioPath
	m.state = StatePlaying
	
	// Play audio
	return m, m.playAudio()
}

func (m Model) handlePlayAudioMsg(msg PlayAudioMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.error = fmt.Errorf("failed to play audio: %w", msg.Err)
		logger.Error("Audio playback failed: %v", msg.Err)
	}

	// Move to listening state regardless of error
	m.state = StateListening
	m.textInput.Reset()
	m.textInput.Focus()
	return m, nil
}

func (m Model) handleGradeDictationMsg(msg GradeDictationMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.error = fmt.Errorf("failed to grade dictation: %w", msg.Err)
		m.state = StateListening
		return m, nil
	}

	logger.Info("Grading complete: Score=%d, WER=%.2f", msg.Grade.Score, msg.Grade.WER)
	m.currentGrade = msg.Grade
	m.currentSession.Grade = msg.Grade
	m.currentSession.Completed = true
	m.currentSession.EndTime = time.Now()
	
	// Save session to history
	if err := m.storage.SaveSession(m.currentSession); err != nil {
		logger.Error("Failed to save session: %v", err)
	}
	
	m.state = StateShowingResult
	return m, nil
}

func (m Model) handleSettingsSavedMsg() (tea.Model, tea.Cmd) {
	// Settings saved successfully
	logger.Info("Settings saved successfully")
	m.state = StateListening
	// Start a new round with updated settings
	return m.startNewRound()
}

// Helper methods
func (m *Model) startNewRound() (Model, tea.Cmd) {
	logger.Info("Starting new dictation round")
	
	// Reset session data
	m.currentSession = &types.DictationSession{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		ConfigUsed: *m.config,
	}
	m.currentSentence = ""
	m.currentAudioPath = ""
	m.currentGrade = nil
	m.error = nil
	m.replayCount = 0
	m.sessionStartTime = time.Now()
	
	// Clear text input
	m.textInput.Reset()
	
	// Start generating sentence
	m.state = StateGenerating
	return *m, m.generateSentence()
}

func (m *Model) openSettings() {
	m.tempConfig = &types.Config{
		Voice:       m.config.Voice,
		Level:       m.config.Level,
		Topic:       m.config.Topic,
		WordCount:   m.config.WordCount,
		SpeechSpeed: m.config.SpeechSpeed,
	}
	m.settingsFocus = 0
	m.state = StateSettings
}

func (m *Model) cancelSettings() {
	m.tempConfig = nil
	m.state = StateListening
}

func (m *Model) saveSettings() {
	if m.tempConfig != nil {
		m.config = m.tempConfig
		m.tempConfig = nil
	}
}

func (m *Model) adjustSetting(direction int) {
	if m.tempConfig == nil {
		return
	}

	switch m.settingsOptions[m.settingsFocus] {
	case "Voice":
		voices := []string{"alloy", "echo", "fable", "onyx", "nova", "shimmer"}
		currentIndex := 0
		for i, v := range voices {
			if v == m.tempConfig.Voice {
				currentIndex = i
				break
			}
		}
		newIndex := (currentIndex + direction + len(voices)) % len(voices)
		m.tempConfig.Voice = voices[newIndex]

	case "Level":
		// Adjust by 50 points
		newLevel := m.tempConfig.Level + (direction * 50)
		if newLevel >= 400 && newLevel <= 990 {
			m.tempConfig.Level = newLevel
		}

	case "Topic":
		topics := []string{"Business", "Travel", "Daily", "Technology", "Health"}
		currentIndex := 0
		for i, t := range topics {
			if t == m.tempConfig.Topic {
				currentIndex = i
				break
			}
		}
		newIndex := (currentIndex + direction + len(topics)) % len(topics)
		m.tempConfig.Topic = topics[newIndex]

	case "Length":
		newLength := m.tempConfig.WordCount + direction
		if newLength >= 5 && newLength <= 30 {
			m.tempConfig.WordCount = newLength
		}

	case "Speed":
		newSpeed := m.tempConfig.SpeechSpeed + (float64(direction) * 0.1)
		if newSpeed >= 0.5 && newSpeed <= 2.0 {
			m.tempConfig.SpeechSpeed = newSpeed
		}
	}
}

// Command generators
func (m Model) generateSentence() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		sentence, err := m.dictationService.GenerateSentence(ctx, m.config)
		return GenerateSentenceMsg{
			Sentence: sentence,
			Err:      err,
		}
	}
}

func (m Model) generateAudio() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		audioPath, err := m.dictationService.GenerateAudio(ctx, m.currentSentence, m.config)
		return GenerateAudioMsg{
			AudioPath: audioPath,
			Err:      err,
		}
	}
}

func (m Model) playAudio() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		err := m.audioPlayer.Play(ctx, m.currentAudioPath)
		return PlayAudioMsg{
			Success: err == nil,
			Err:     err,
		}
	}
}

func (m Model) gradeDictation() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		grade, err := m.dictationService.GradeDictation(ctx, m.currentSentence, m.currentSession.UserInput, m.config)
		return GradeDictationMsg{
			Grade: grade,
			Err:   err,
		}
	}
}

func (m Model) saveSettingsCmd() tea.Cmd {
	return func() tea.Msg {
		// Here we would save to config file
		// For now, just return success
		return SettingsSavedMsg{}
	}
}