package tui

import (
	"github.com/yourusername/dictcli/internal/types"
)

// Message types for Bubble Tea
type GenerateSentenceMsg struct {
	Sentence string
	Err      error
}

type GenerateAudioMsg struct {
	AudioPath string
	Err       error
}

type PlayAudioMsg struct {
	Success bool
	Err     error
}

type GradeDictationMsg struct {
	Grade *types.Grade
	Err   error
}

type SettingsSavedMsg struct{}

type ErrorMsg struct {
	Err error
}

type TickMsg struct{}