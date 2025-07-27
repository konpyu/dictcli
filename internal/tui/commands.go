package tui

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func generateSessionID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (m Model) generateSentence() tea.Msg {
	sentence, err := m.service.GenerateSentence(m.cfg.Topic, m.cfg.Level, m.cfg.Words)
	return generatedMsg{sentence: sentence, err: err}
}

func (m Model) generateAudio() tea.Msg {
	if m.currentSession == nil || m.currentSession.Sentence == "" {
		return audioGeneratedMsg{err: fmt.Errorf("no sentence to generate audio for")}
	}
	
	audioPath, err := m.service.GenerateAudio(m.currentSession.Sentence, m.cfg.Voice, m.cfg.Speed)
	return audioGeneratedMsg{audioPath: audioPath, err: err}
}

func (m Model) playAudio() tea.Msg {
	if m.currentSession == nil || m.currentSession.AudioPath == "" {
		return audioPlayedMsg{err: fmt.Errorf("no audio to play")}
	}
	
	if m.currentSession.ReplayCount > 0 {
		m.currentSession.ReplayCount++
	} else {
		m.currentSession.ReplayCount = 1
	}
	
	err := m.audioPlayer.Play(m.currentSession.AudioPath)
	return audioPlayedMsg{err: err}
}

func (m Model) gradeDictation() tea.Msg {
	if m.currentSession == nil {
		return gradedMsg{err: fmt.Errorf("no session to grade")}
	}
	
	grade, err := m.service.GradeDictation(m.currentSession.Sentence, m.currentSession.UserInput)
	return gradedMsg{grade: grade, err: err}
}

func (m Model) saveSession() tea.Msg {
	if m.currentSession == nil {
		return sessionSavedMsg{err: fmt.Errorf("no session to save")}
	}
	
	err := m.history.SaveSession(m.currentSession)
	return sessionSavedMsg{err: err}
}