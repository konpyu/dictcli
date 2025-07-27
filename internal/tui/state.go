package tui

type State int

const (
	StateWelcome State = iota
	StateGenerating
	StatePlaying
	StateListening
	StateGrading
	StateShowingResult
	StateSettings
)

func (s State) String() string {
	switch s {
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
	default:
		return "Unknown"
	}
}