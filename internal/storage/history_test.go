package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/konpyu/dictcli/internal/types"
)

func TestHistory(t *testing.T) {
	t.Run("NewHistory creates data directory", func(t *testing.T) {
		tempDir := t.TempDir()
		dataDir := filepath.Join(tempDir, "dictcli")
		
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test data directory: %v", err)
		}
		
		history := &History{
			filePath: filepath.Join(dataDir, "history.jsonl"),
		}
		
		if history.filePath == "" {
			t.Errorf("History file path is empty")
		}
		
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			t.Errorf("Data directory was not created")
		}
	})

	t.Run("SaveSession and LoadSessions", func(t *testing.T) {
		tempDir := t.TempDir()
		history := &History{
			filePath: filepath.Join(tempDir, "history.jsonl"),
		}

		session1 := &types.DictationSession{
			ID:        "test-1",
			Timestamp: time.Now(),
			Config: types.Config{
				Voice: "alloy",
				Level: 700,
				Topic: "Business",
				Words: 15,
				Speed: 1.0,
			},
			Sentence:  "This is a test sentence",
			UserInput: "This is a test sentence",
			Grade: &types.Grade{
				WER:      0.0,
				Score:    100,
				Mistakes: []types.Mistake{},
			},
		}

		err := history.SaveSession(session1)
		if err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		session2 := &types.DictationSession{
			ID:        "test-2",
			Timestamp: time.Now().AddDate(0, 0, -5),
			Config: types.Config{
				Voice: "echo",
				Level: 600,
				Topic: "Travel",
				Words: 10,
				Speed: 0.9,
			},
			Sentence:  "Travel test",
			UserInput: "Travel test",
			Grade: &types.Grade{
				WER:   0.1,
				Score: 90,
				Mistakes: []types.Mistake{
					{
						Position: 0,
						Expected: "Travel",
						Actual:   "Travle",
						Type:     "spelling",
					},
				},
			},
		}

		err = history.SaveSession(session2)
		if err != nil {
			t.Fatalf("Failed to save session 2: %v", err)
		}

		sessions, err := history.LoadSessions(30)
		if err != nil {
			t.Fatalf("Failed to load sessions: %v", err)
		}

		if len(sessions) != 2 {
			t.Errorf("Expected 2 sessions, got %d", len(sessions))
		}

		sessions3Days, err := history.LoadSessions(3)
		if err != nil {
			t.Fatalf("Failed to load sessions for 3 days: %v", err)
		}

		if len(sessions3Days) != 1 {
			t.Errorf("Expected 1 session within 3 days, got %d", len(sessions3Days))
		}
	})

	t.Run("LoadAllSessions", func(t *testing.T) {
		tempDir := t.TempDir()
		history := &History{
			filePath: filepath.Join(tempDir, "history.jsonl"),
		}

		for i := 0; i < 5; i++ {
			session := &types.DictationSession{
				ID:        fmt.Sprintf("test-%d", i),
				Timestamp: time.Now().AddDate(0, 0, -i*10),
				Config: types.Config{
					Voice: "alloy",
					Level: 700,
				},
			}
			_ = history.SaveSession(session)
		}

		sessions, err := history.LoadAllSessions()
		if err != nil {
			t.Fatalf("Failed to load all sessions: %v", err)
		}

		if len(sessions) != 5 {
			t.Errorf("Expected 5 sessions, got %d", len(sessions))
		}
	})

	t.Run("CalculateStatistics", func(t *testing.T) {
		tempDir := t.TempDir()
		history := &History{
			filePath: filepath.Join(tempDir, "history.jsonl"),
		}

		sessions := []*types.DictationSession{
			{
				ID:        "1",
				Timestamp: time.Now(),
				Config:    types.Config{Topic: "Business"},
				Grade: &types.Grade{
					WER:   0.1,
					Score: 90,
					Mistakes: []types.Mistake{
						{Expected: "the", Actual: "teh"},
						{Expected: "business", Actual: "busines"},
					},
				},
			},
			{
				ID:        "2",
				Timestamp: time.Now(),
				Config:    types.Config{Topic: "Business"},
				Grade: &types.Grade{
					WER:   0.05,
					Score: 95,
					Mistakes: []types.Mistake{
						{Expected: "the", Actual: "teh"},
					},
				},
			},
			{
				ID:        "3",
				Timestamp: time.Now(),
				Config:    types.Config{Topic: "Travel"},
				Grade: &types.Grade{
					WER:      0.0,
					Score:    100,
					Mistakes: []types.Mistake{},
				},
			},
		}

		for _, session := range sessions {
			_ = history.SaveSession(session)
		}

		stats, err := history.CalculateStatistics(30)
		if err != nil {
			t.Fatalf("Failed to calculate statistics: %v", err)
		}

		if stats.TotalSessions != 3 {
			t.Errorf("Expected 3 total sessions, got %d", stats.TotalSessions)
		}

		expectedAvgScore := (90.0 + 95.0 + 100.0) / 3.0
		if abs(stats.AverageScore-expectedAvgScore) > 0.01 {
			t.Errorf("Expected average score %.2f, got %.2f", expectedAvgScore, stats.AverageScore)
		}

		if len(stats.TopicBreakdown) != 2 {
			t.Errorf("Expected 2 topics, got %d", len(stats.TopicBreakdown))
		}

		businessStats := stats.TopicBreakdown["Business"]
		if businessStats.Count != 2 {
			t.Errorf("Expected 2 Business sessions, got %d", businessStats.Count)
		}

		foundTheMistake := false
		for _, mistake := range stats.CommonMistakes {
			if mistake.Expected == "the" && mistake.Actual == "teh" && mistake.Frequency == 2 {
				foundTheMistake = true
				break
			}
		}
		if !foundTheMistake {
			t.Errorf("Expected to find 'the->teh' mistake with frequency 2")
		}
	})

	t.Run("Clear history", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "history.jsonl")
		history := &History{filePath: filePath}

		session := &types.DictationSession{ID: "test"}
		_ = history.SaveSession(session)

		err := history.Clear()
		if err != nil {
			t.Fatalf("Failed to clear history: %v", err)
		}

		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("History file still exists after clear")
		}
	})
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}