package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/konpyu/dictcli/internal/config"
	"github.com/konpyu/dictcli/internal/service"
	"github.com/konpyu/dictcli/internal/storage"
	"github.com/konpyu/dictcli/internal/types"
)

type Model struct {
	state         State
	prevState     State
	welcomeShown  bool
	width         int
	height        int
	
	cfg           *types.Config
	configManager *config.Manager
	service       *service.DictationService
	history       *storage.History
	audioPlayer   *storage.AudioPlayer
	
	currentSession *types.DictationSession
	userInput      string
	
	spinner        spinner.Model
	textInput      textinput.Model
	
	err            error
	message        string
	
	settingsFields []string
	settingsIndex  int
}

func New(configManager *config.Manager, service *service.DictationService, history *storage.History, player *storage.AudioPlayer) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	ti := textinput.New()
	ti.Placeholder = "Type what you hear..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80
	
	return Model{
		state:         StateWelcome,
		cfg:           configManager.Get(),
		configManager: configManager,
		service:       service,
		history:       history,
		audioPlayer:   player,
		spinner:       s,
		textInput:     ti,
		settingsFields: []string{"voice", "level", "topic", "words", "speed"},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		textinput.Blink,
	)
}

// SetState sets the current state (for testing)
func (m *Model) SetState(state State) {
	m.state = state
}