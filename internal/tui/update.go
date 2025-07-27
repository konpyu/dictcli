package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = min(msg.Width-4, 80)
		return m, nil

	case tea.KeyMsg:
		// Handle global key presses first
		if newModel, cmd, handled := m.handleGlobalKeys(msg); handled {
			return newModel, cmd
		}
		// Then let state-specific handlers process the key
		return m.handleStateKeys(msg)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case errMsg:
		m.err = msg.err
		m.state = m.prevState
		return m, nil

	case stateChangeMsg:
		return m.changeState(msg.newState)
	}

	switch m.state {
	case StateWelcome:
		return m.updateWelcome(msg)
	case StateGenerating:
		return m.updateGenerating(msg)
	case StatePlaying:
		return m.updatePlaying(msg)
	case StateListening:
		return m.updateListening(msg)
	case StateGrading:
		return m.updateGrading(msg)
	case StateShowingResult:
		return m.updateShowingResult(msg)
	case StateSettings:
		return m.updateSettings(msg)
	case StateHelp:
		return m.updateHelp(msg)
	}

	return m, nil
}

func (m Model) handleGlobalKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	// Ctrl+C always quits
	if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
		return m, tea.Quit, true
	}

	// Ctrl+Q always quits
	if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+q"))) {
		return m, tea.Quit, true
	}

	// Help shortcuts
	if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+h", "?"))) {
		newModel, cmd := m.changeState(StateHelp)
		return newModel, cmd, true
	}

	// Esc goes back to previous state (except from welcome)
	if key.Matches(msg, key.NewBinding(key.WithKeys("esc"))) && m.state != StateWelcome {
		if m.state == StateHelp || m.state == StateSettings {
			newModel, cmd := m.changeState(m.prevState)
			return newModel, cmd, true
		}
	}

	// Q/q quits in most states (except when typing, settings, or help)
	if m.state != StateListening && m.state != StateSettings && m.state != StateHelp && key.Matches(msg, key.NewBinding(key.WithKeys("q", "Q"))) {
		return m, tea.Quit, true
	}

	return m, nil, false
}

func (m Model) handleStateKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateWelcome:
		m.welcomeShown = true
		return m.changeState(StateGenerating)

	case StateListening:
		// Don't handle special keys when text input has focus and user is typing
		// Only handle special keys with Ctrl modifier
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+r"))) {
			return m.changeState(StatePlaying)
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))) {
			return m.changeState(StateSettings)
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+h", "?"))) {
			return m.changeState(StateHelp)
		}
		// Pass through to updateListening for text input
		return m.updateListening(msg)

	case StateShowingResult:
		switch strings.ToLower(msg.String()) {
		case "n":
			return m.changeState(StateGenerating)
		case "r":
			return m.changeState(StatePlaying)
		case "s":
			return m.changeState(StateSettings)
		case "q":
			return m, tea.Quit
		}

	case StateSettings:
		// Pass through to updateSettings for arrow keys, etc.
		return m.updateSettings(msg)

	case StateHelp:
		// Help state only responds to global keys (Esc, Ctrl+Q)
	}

	return m, nil
}

func (m Model) changeState(newState State) (Model, tea.Cmd) {
	m.prevState = m.state
	m.state = newState
	m.err = nil

	switch newState {
	case StateGenerating:
		return m, m.generateSentence
	case StatePlaying:
		if m.currentSession != nil && m.currentSession.AudioPath != "" {
			return m, m.playAudio
		}
		return m, m.generateAudio
	case StateGrading:
		return m, m.gradeDictation
	}

	return m, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}