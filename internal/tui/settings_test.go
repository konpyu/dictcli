package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/konpyu/dictcli/internal/config"
	"github.com/konpyu/dictcli/internal/service"
	"github.com/konpyu/dictcli/internal/storage"
	"github.com/konpyu/dictcli/internal/types"
)

func TestSettingsKeyHandling(t *testing.T) {
	// Create a minimal model for testing
	configManager := &config.Manager{}
	dictationService := &service.DictationService{}
	history := &storage.History{}
	audioPlayer := &storage.AudioPlayer{}
	
	model := New(configManager, dictationService, history, audioPlayer)
	model.state = StateSettings
	model.cfg = &types.Config{
		Voice: "alloy",
		Level: 700,
		Topic: "Business",
		Words: 15,
		Speed: 1.0,
	}
	model.settingsIndex = 0

	// Test Up key
	keyMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(keyMsg)
	m := updatedModel.(Model)
	
	// Settings index should stay at 0 (can't go below 0)
	if m.settingsIndex != 0 {
		t.Errorf("Expected settingsIndex to remain 0, got %d", m.settingsIndex)
	}

	// Test Down key
	keyMsg = tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = m.Update(keyMsg)
	m = updatedModel.(Model)
	
	// Settings index should move to 1
	if m.settingsIndex != 1 {
		t.Errorf("Expected settingsIndex to be 1, got %d", m.settingsIndex)
	}

	// Test Left key (should change voice)
	originalVoice := m.cfg.Voice
	keyMsg = tea.KeyMsg{Type: tea.KeyLeft}
	updatedModel, _ = m.Update(keyMsg)
	m = updatedModel.(Model)
	
	// Voice should have changed (we're on index 1 which is level, but let's test voice change too)
	m.settingsIndex = 0 // Set to voice field
	keyMsg = tea.KeyMsg{Type: tea.KeyLeft}
	updatedModel, _ = m.Update(keyMsg)
	m = updatedModel.(Model)
	
	if m.cfg.Voice == originalVoice {
		t.Errorf("Expected voice to change from %s", originalVoice)
	}

	// Test Esc key should change state back
	model.prevState = StateListening
	keyMsg = tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ = model.Update(keyMsg)
	m = updatedModel.(Model)
	
	if m.state != StateListening {
		t.Errorf("Expected state to change back to StateListening, got %v", m.state)
	}
}

func TestSettingsAdjustments(t *testing.T) {
	model := Model{
		cfg: &types.Config{
			Voice: "alloy",
			Level: 700,
			Topic: "Business", 
			Words: 15,
			Speed: 1.0,
		},
		settingsFields: []string{"voice", "level", "topic", "words", "speed"},
		settingsIndex:  1, // level
	}

	// Test level adjustment
	model.adjustSetting(1)
	if model.cfg.Level != 750 {
		t.Errorf("Expected level to be 750, got %d", model.cfg.Level)
	}

	model.adjustSetting(-1)
	if model.cfg.Level != 700 {
		t.Errorf("Expected level to be back to 700, got %d", model.cfg.Level)
	}

	// Test level bounds
	model.cfg.Level = 990
	model.adjustSetting(1)
	if model.cfg.Level != 990 {
		t.Errorf("Expected level to stay at 990 (max), got %d", model.cfg.Level)
	}

	model.cfg.Level = 400
	model.adjustSetting(-1)
	if model.cfg.Level != 400 {
		t.Errorf("Expected level to stay at 400 (min), got %d", model.cfg.Level)
	}
}