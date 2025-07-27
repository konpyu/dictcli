package config

import (
	"testing"

	"github.com/konpyu/dictcli/internal/types"
)

func TestConfigManagerLogic(t *testing.T) {
	t.Run("Config manager creation and defaults", func(t *testing.T) {
		manager, err := New()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		cfg := manager.Get()
		if cfg == nil {
			t.Fatal("expected config to be non-nil")
		}

		if !types.IsValidVoice(cfg.Voice) {
			t.Errorf("voice %s is not valid", cfg.Voice)
		}
		
		if cfg.Level < types.MinLevel || cfg.Level > types.MaxLevel {
			t.Errorf("level %d is out of range", cfg.Level)
		}
		
		if !types.IsValidTopic(cfg.Topic) {
			t.Errorf("topic %s is not valid", cfg.Topic)
		}
	})

	t.Run("Set config values", func(t *testing.T) {
		manager, err := New()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if err := manager.Set("voice", "echo"); err != nil {
			t.Fatalf("expected no error setting voice, got %v", err)
		}

		if err := manager.Set("level", 850); err != nil {
			t.Fatalf("expected no error setting level, got %v", err)
		}

		cfg := manager.Get()
		if cfg.Voice != "echo" {
			t.Errorf("expected voice echo, got %s", cfg.Voice)
		}
		if cfg.Level != 850 {
			t.Errorf("expected level 850, got %d", cfg.Level)
		}
	})

	t.Run("SetFromFlags validates values", func(t *testing.T) {
		manager, err := New()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		manager.SetFromFlags("nova", "Technology", 600, 20, 1.5, true, false)

		cfg := manager.Get()
		if cfg.Voice != "nova" {
			t.Errorf("expected voice nova, got %s", cfg.Voice)
		}
		if cfg.Topic != "Technology" {
			t.Errorf("expected topic Technology, got %s", cfg.Topic)
		}
		if cfg.Level != 600 {
			t.Errorf("expected level 600, got %d", cfg.Level)
		}
		if cfg.Words != 20 {
			t.Errorf("expected words 20, got %d", cfg.Words)
		}
		if cfg.Speed != 1.5 {
			t.Errorf("expected speed 1.5, got %f", cfg.Speed)
		}
		if !cfg.NoCache {
			t.Error("expected NoCache to be true")
		}
	})

	t.Run("SetFromFlags ignores invalid values", func(t *testing.T) {
		manager, err := New()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		oldCfg := *manager.Get()

		manager.SetFromFlags("invalid_voice", "invalid_topic", 50, 50, 5.0, false, false)

		cfg := manager.Get()
		if cfg.Voice != oldCfg.Voice {
			t.Errorf("voice should not have changed from %s to %s", oldCfg.Voice, cfg.Voice)
		}
		if cfg.Topic != oldCfg.Topic {
			t.Errorf("topic should not have changed from %s to %s", oldCfg.Topic, cfg.Topic)
		}
		if cfg.Level != oldCfg.Level {
			t.Errorf("level should not have changed from %d to %d", oldCfg.Level, cfg.Level)
		}
		if cfg.Words != oldCfg.Words {
			t.Errorf("words should not have changed from %d to %d", oldCfg.Words, cfg.Words)
		}
		if cfg.Speed != oldCfg.Speed {
			t.Errorf("speed should not have changed from %f to %f", oldCfg.Speed, cfg.Speed)
		}
	})
}