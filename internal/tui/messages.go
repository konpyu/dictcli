package tui

import (
	"github.com/konpyu/dictcli/internal/types"
)

type generatedMsg struct {
	sentence string
	err      error
}

type audioGeneratedMsg struct {
	audioPath string
	err       error
}

type audioPlayedMsg struct {
	err error
}

type gradedMsg struct {
	grade *types.Grade
	err   error
}

type sessionSavedMsg struct {
	err error
}

type stateChangeMsg struct {
	newState State
}

type errMsg struct {
	err error
}

type clearMessageMsg struct{}