package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/dictcli/internal/logger"
	"github.com/yourusername/dictcli/internal/service"
	"github.com/yourusername/dictcli/internal/storage"
	"github.com/yourusername/dictcli/internal/types"
)

// State represents the current state of the TUI
type State int

const (
	StateWelcome State = iota
	StateGenerating
	StatePlaying
	StateListening
	StateGrading
	StateShowingResult
	StateSettings
	StateQuitting
)

// Model represents the TUI model
type Model struct {
	// Current state
	state State

	// Services
	dictationService service.DictationService
	audioPlayer      service.AudioPlayer
	audioCache       service.AudioCache
	storage          storage.Storage

	// Configuration
	config *types.Config

	// Current session data
	currentSession   *types.DictationSession
	currentSentence  string
	currentAudioPath string
	currentGrade     *types.Grade

	// UI components
	textInput    textinput.Model
	spinner      spinner.Model
	settingsFocus int // Which setting is currently focused

	// UI state
	width           int
	height          int
	error           error
	welcomeShown    bool
	replayCount     int
	sessionStartTime time.Time

	// Settings modal specific
	tempConfig      *types.Config // Temporary config for settings modal
	settingsOptions []string
}

// New creates a new TUI model
func New(
	dictationService service.DictationService,
	audioPlayer service.AudioPlayer,
	audioCache service.AudioCache,
	storage storage.Storage,
	config *types.Config,
) Model {
	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Type what you hear..."
	ti.Focus()
	ti.CharLimit = 300
	ti.Width = 60

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	logger.Info("Creating new TUI model with config: Voice=%s, Level=%d, Topic=%s, WordCount=%d",
		config.Voice, config.Level, config.Topic, config.WordCount)

	return Model{
		state:            StateWelcome,
		dictationService: dictationService,
		audioPlayer:      audioPlayer,
		audioCache:       audioCache,
		storage:          storage,
		config:           config,
		textInput:        ti,
		spinner:          s,
		settingsFocus:    0,
		welcomeShown:     false,
		replayCount:      0,
		settingsOptions: []string{
			"Voice",
			"Level",
			"Topic",
			"Length",
			"Speed",
		},
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	logger.Info("Initializing TUI model")
	return tea.Batch(
		m.spinner.Tick,
		textinput.Blink,
	)
}

// getCurrentStateString returns a string representation of the current state
func (m Model) getCurrentStateString() string {
	switch m.state {
	case StateWelcome:
		return "Welcome"
	case StateGenerating:
		return "Generating"
	case StatePlaying:
		return "Playing"
	case StateListening:
		return "Listening"
	case StateGrading:
		return "Grading"
	case StateShowingResult:
		return "ShowingResult"
	case StateSettings:
		return "Settings"
	case StateQuitting:
		return "Quitting"
	default:
		return "Unknown"
	}
}