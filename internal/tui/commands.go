package tui

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/konpyu/dictcli/internal/logging"
)

func generateSessionID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (m Model) generateSentence() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	logging.Info("Generating sentence - Topic: %s, Level: %d, Words: %d", m.cfg.Topic, m.cfg.Level, m.cfg.Words)
	sentence, err := m.service.GenerateSentence(ctx, m.cfg.Topic, m.cfg.Level, m.cfg.Words)
	if err != nil {
		logging.Error("Failed to generate sentence: %v", err)
	} else {
		logging.Info("Generated sentence successfully")
	}
	return generatedMsg{sentence: sentence, err: err}
}

func (m Model) generateAudio() tea.Msg {
	if m.currentSession == nil || m.currentSession.Sentence == "" {
		return audioGeneratedMsg{err: fmt.Errorf("音声を生成する文がありません")}
	}
	
	logging.Info("Generating audio - Voice: %s, Speed: %.1f", m.cfg.Voice, m.cfg.Speed)
	
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	
	audioPath, err := m.service.GenerateAudio(ctx, m.currentSession.Sentence, m.cfg.Voice, m.cfg.Speed)
	return audioGeneratedMsg{audioPath: audioPath, err: err}
}

func (m Model) playAudio() tea.Msg {
	if m.currentSession == nil || m.currentSession.AudioPath == "" {
		return audioPlayedMsg{err: fmt.Errorf("再生する音声ファイルがありません")}
	}
	
	if m.currentSession.ReplayCount > 0 {
		m.currentSession.ReplayCount++
	} else {
		m.currentSession.ReplayCount = 1
	}
	
	err := m.audioPlayer.Play(m.currentSession.AudioPath)
	if err != nil {
		return audioPlayedMsg{err: fmt.Errorf("音声再生エラー: %w", err)}
	}
	return audioPlayedMsg{err: nil}
}

func (m Model) gradeDictation() tea.Msg {
	if m.currentSession == nil {
		return gradedMsg{err: fmt.Errorf("採点するセッションがありません")}
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	grade, err := m.service.GradeDictation(ctx, m.currentSession.Sentence, m.currentSession.UserInput)
	return gradedMsg{grade: grade, err: err}
}

func (m Model) saveSession() tea.Msg {
	if m.currentSession == nil {
		logging.Warn("No session to save")
		return sessionSavedMsg{err: fmt.Errorf("保存するセッションがありません")}
	}
	
	logging.Info("Saving session - ID: %s", m.currentSession.ID)
	err := m.history.SaveSession(m.currentSession)
	if err != nil {
		logging.Error("Failed to save session: %v", err)
		return sessionSavedMsg{err: fmt.Errorf("セッションの保存に失敗しました: %w", err)}
	}
	logging.Info("Session saved successfully")
	return sessionSavedMsg{err: nil}
}