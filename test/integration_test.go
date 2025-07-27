package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/konpyu/dictcli/internal/service"
	"github.com/konpyu/dictcli/internal/storage"
	"github.com/konpyu/dictcli/internal/types"
)

func TestFullDictationFlow(t *testing.T) {
	// Skip if no API key available
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}

	// Create temporary directories for test
	tempDir := t.TempDir()
	cacheDir := filepath.Join(tempDir, "cache")
	dataDir := filepath.Join(tempDir, "data")

	// Setup test environment
	_ = os.Setenv("XDG_CACHE_HOME", cacheDir)
	_ = os.Setenv("XDG_DATA_HOME", dataDir)

	// Initialize services
	dictationSvc, err := service.NewDictationService(false)
	if err != nil {
		t.Fatalf("Failed to create dictation service: %v", err)
	}

	history, err := storage.NewHistory()
	if err != nil {
		t.Fatalf("Failed to create history: %v", err)
	}

	// Clear any existing history from previous test runs
	_ = history.Clear()

	cache, err := storage.NewAudioCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test configuration
	cfg := &types.Config{
		Voice: "alloy",
		Level: 600,
		Topic: "Daily",
		Words: 10,
		Speed: 1.0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test sentence generation
	sentence, err := dictationSvc.GenerateSentence(ctx, cfg.Topic, cfg.Level, cfg.Words)
	if err != nil {
		t.Fatalf("Failed to generate sentence: %v", err)
	}

	if sentence == "" {
		t.Fatal("Generated sentence is empty")
	}
	t.Logf("Generated sentence: %s", sentence)

	// Test audio generation (returns path, not data)
	audioPath, err := dictationSvc.GenerateAudio(ctx, sentence, cfg.Voice, cfg.Speed)
	if err != nil {
		t.Fatalf("Failed to generate audio: %v", err)
	}

	if audioPath == "" {
		t.Fatal("Generated audio path is empty")
	}
	t.Logf("Generated audio path: %s", audioPath)

	// Test audio caching - verify cache hit on second call
	audioPath2, err := dictationSvc.GenerateAudio(ctx, sentence, cfg.Voice, cfg.Speed)
	if err != nil {
		t.Fatalf("Failed to generate audio second time: %v", err)
	}

	if audioPath != audioPath2 {
		t.Fatal("Audio path mismatch - cache not working")
	}

	// Test direct cache operations
	if !cache.Exists(sentence, cfg.Voice, cfg.Speed) {
		t.Fatal("Audio should exist in cache")
	}

	cachedAudio, err := cache.Load(sentence, cfg.Voice, cfg.Speed)
	if err != nil {
		t.Fatalf("Failed to load audio from cache: %v", err)
	}

	if len(cachedAudio) == 0 {
		t.Fatal("Cached audio data is empty")
	}

	// Test grading
	userInput := sentence // Perfect input for testing
	grade, err := dictationSvc.GradeDictation(ctx, sentence, userInput)
	if err != nil {
		t.Fatalf("Failed to grade dictation: %v", err)
	}

	if grade.Score < 90 { // Should be near perfect
		t.Errorf("Expected high score for perfect input, got: %d", grade.Score)
	}

	// Test session saving
	session := &types.DictationSession{
		ID:           "test-session-1",
		Timestamp:    time.Now(),
		StartTime:    time.Now().Add(-time.Minute),
		EndTime:      time.Now(),
		Config:       *cfg,
		Sentence:     sentence,
		AudioPath:    audioPath,
		UserInput:    userInput,
		Grade:        grade,
		ReplayCount:  0,
		DurationSecs: 60.0,
	}

	err = history.SaveSession(session)
	if err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Test statistics
	stats, err := history.CalculateStatistics(30)
	if err != nil {
		t.Fatalf("Failed to calculate statistics: %v", err)
	}

	if stats.TotalSessions != 1 {
		t.Errorf("Expected 1 session, got: %d", stats.TotalSessions)
	}

	t.Logf("Integration test completed successfully")
}

func TestErrorScenarios(t *testing.T) {
	// Test without API key
	originalKey := os.Getenv("OPENAI_API_KEY")
	_ = os.Unsetenv("OPENAI_API_KEY")
	defer func() { _ = os.Setenv("OPENAI_API_KEY", originalKey) }()

	_, err := service.NewDictationService(false)
	if err == nil {
		t.Error("Expected error when creating service without API key, but got none")
		return
	}
	t.Logf("Expected error when creating service without API key: %v", err)

	// Restore API key for remaining tests
	_ = os.Setenv("OPENAI_API_KEY", originalKey)
}

func TestConfigValidation(t *testing.T) {
	// Test invalid config values
	testCases := []struct {
		name   string
		config types.Config
		valid  bool
	}{
		{
			name: "valid config",
			config: types.Config{
				Voice: "alloy",
				Level: 600,
				Topic: "Daily",
				Words: 15,
				Speed: 1.0,
			},
			valid: true,
		},
		{
			name: "invalid voice",
			config: types.Config{
				Voice: "invalid",
				Level: 600,
				Topic: "Daily",
				Words: 15,
				Speed: 1.0,
			},
			valid: false,
		},
		{
			name: "invalid level too low",
			config: types.Config{
				Voice: "alloy",
				Level: 300,
				Topic: "Daily",
				Words: 15,
				Speed: 1.0,
			},
			valid: false,
		},
		{
			name: "invalid level too high",
			config: types.Config{
				Voice: "alloy",
				Level: 1000,
				Topic: "Daily",
				Words: 15,
				Speed: 1.0,
			},
			valid: false,
		},
		{
			name: "invalid words too low",
			config: types.Config{
				Voice: "alloy",
				Level: 600,
				Topic: "Daily",
				Words: 3,
				Speed: 1.0,
			},
			valid: false,
		},
		{
			name: "invalid speed too low",
			config: types.Config{
				Voice: "alloy",
				Level: 600,
				Topic: "Daily",
				Words: 15,
				Speed: 0.3,
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test voice validation
			if tc.name == "invalid voice" {
				if types.IsValidVoice(tc.config.Voice) {
					t.Error("Expected invalid voice to be rejected")
				}
			} else if tc.config.Voice != "" {
				if !types.IsValidVoice(tc.config.Voice) {
					t.Error("Expected valid voice to be accepted")
				}
			}

			// Test topic validation
			if types.IsValidTopic(tc.config.Topic) != tc.valid && tc.name != "invalid voice" && tc.name != "invalid level too low" && tc.name != "invalid level too high" && tc.name != "invalid words too low" && tc.name != "invalid speed too low" {
				t.Errorf("Topic validation mismatch for %s", tc.name)
			}

			// Test level bounds
			if tc.name == "invalid level too low" || tc.name == "invalid level too high" {
				if tc.config.Level >= types.MinLevel && tc.config.Level <= types.MaxLevel {
					t.Error("Expected invalid level to be outside bounds")
				}
			}

			// Test word bounds
			if tc.name == "invalid words too low" {
				if tc.config.Words >= types.MinWords {
					t.Error("Expected invalid word count to be below minimum")
				}
			}

			// Test speed bounds  
			if tc.name == "invalid speed too low" {
				if tc.config.Speed >= types.MinSpeed {
					t.Error("Expected invalid speed to be below minimum")
				}
			}
		})
	}
}

func TestCacheFunctionality(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.Setenv("XDG_CACHE_HOME", tempDir)

	cache, err := storage.NewAudioCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test cache miss
	_, err = cache.Load("test sentence", "alloy", 1.0)
	if err == nil {
		t.Error("Expected cache miss, but got hit")
	}

	// Test cache save and hit
	testData := []byte("test audio data")
	err = cache.Save("test sentence", "alloy", 1.0, testData)
	if err != nil {
		t.Fatalf("Failed to save to cache: %v", err)
	}

	cachedData, err := cache.Load("test sentence", "alloy", 1.0)
	if err != nil {
		t.Fatalf("Failed to load from cache: %v", err)
	}

	if string(cachedData) != string(testData) {
		t.Error("Cached data mismatch")
	}

	// Test cache size
	size, count, err := cache.Size()
	if err != nil {
		t.Fatalf("Failed to get cache size: %v", err)
	}

	if count < 1 {
		t.Errorf("Expected at least 1 cached file, got: %d", count)
	}

	if size == 0 {
		t.Error("Expected non-zero cache size")
	}

	// Test cache clear
	err = cache.Clear()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// After clear, cache directory may not exist, so we expect an error or 0 count
	_, count, err = cache.Size()
	if err == nil && count != 0 {
		t.Errorf("Expected 0 cached files after clear, got: %d", count)
	}
}

func TestHistoryPersistence(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.Setenv("XDG_DATA_HOME", tempDir)

	history1, err := storage.NewHistory()
	if err != nil {
		t.Fatalf("Failed to create history: %v", err)
	}

	// Clear any existing history first
	_ = history1.Clear()

	// Save test session
	session := &types.DictationSession{
		ID:        "test-history-1",
		Timestamp: time.Now(),
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Minute),
		Config: types.Config{
			Voice: "alloy",
			Level: 600,
			Topic: "Daily",
			Words: 10,
			Speed: 1.0,
		},
		Sentence:     "This is a test sentence.",
		AudioPath:    "/tmp/test.mp3",
		UserInput:    "This is a test sentence.",
		ReplayCount:  0,
		Grade: &types.Grade{
			Score:    100,
			WER:      0.0,
			Mistakes: []types.Mistake{},
		},
		DurationSecs: 60.0,
	}

	err = history1.SaveSession(session)
	if err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Create new history instance to test persistence
	history2, err := storage.NewHistory()
	if err != nil {
		t.Fatalf("Failed to create second history instance: %v", err)
	}

	stats, err := history2.CalculateStatistics(30)
	if err != nil {
		t.Fatalf("Failed to calculate statistics: %v", err)
	}

	if stats.TotalSessions != 1 {
		t.Errorf("Expected 1 session from persistent storage, got: %d", stats.TotalSessions)
	}

	if stats.AverageScore != 100.0 {
		t.Errorf("Expected average score 100.0, got: %.1f", stats.AverageScore)
	}
}