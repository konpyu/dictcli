package test

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/konpyu/dictcli/internal/config"
	"github.com/konpyu/dictcli/internal/service"
	"github.com/konpyu/dictcli/internal/storage"
	"github.com/konpyu/dictcli/internal/types"
)

func BenchmarkServiceInitialization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := service.NewDictationService(false)
		if err != nil {
			b.Fatalf("Failed to create dictation service: %v", err)
		}
	}
}

func BenchmarkConfigLoad(b *testing.B) {
	tempDir := b.TempDir()
	_ = os.Setenv("XDG_CONFIG_HOME", tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := config.New()
		if err != nil {
			b.Fatalf("Failed to create config: %v", err)
		}
	}
}

func BenchmarkCacheOperations(b *testing.B) {
	tempDir := b.TempDir()
	_ = os.Setenv("XDG_CACHE_HOME", tempDir)

	cache, err := storage.NewAudioCache()
	if err != nil {
		b.Fatalf("Failed to create cache: %v", err)
	}

	testData := make([]byte, 50*1024) // 50KB typical audio file

	b.Run("Save", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := cache.Save("test sentence", "alloy", 1.0, testData)
			if err != nil {
				b.Fatalf("Failed to save to cache: %v", err)
			}
		}
	})

	// Save once for load test
	_ = cache.Save("test sentence", "alloy", 1.0, testData)

	b.Run("Load", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := cache.Load("test sentence", "alloy", 1.0)
			if err != nil {
				b.Fatalf("Failed to load from cache: %v", err)
			}
		}
	})
}

func BenchmarkHistoryOperations(b *testing.B) {
	tempDir := b.TempDir()
	_ = os.Setenv("XDG_DATA_HOME", tempDir)

	history, err := storage.NewHistory()
	if err != nil {
		b.Fatalf("Failed to create history: %v", err)
	}

	session := &types.DictationSession{
		ID:        "bench-session",
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
		Sentence:     "This is a test sentence for benchmarking.",
		AudioPath:    "/tmp/test.mp3",
		UserInput:    "This is a test sentence for benchmarking.",
		Grade: &types.Grade{
			Score:    95,
			WER:      0.05,
			Mistakes: []types.Mistake{},
		},
		ReplayCount:  0,
		DurationSecs: 60.0,
	}

	b.Run("SaveSession", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := history.SaveSession(session)
			if err != nil {
				b.Fatalf("Failed to save session: %v", err)
			}
		}
	})

	// Save some sessions for stats test
	for i := 0; i < 100; i++ {
		_ = history.SaveSession(session)
	}

	b.Run("CalculateStatistics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := history.CalculateStatistics(30)
			if err != nil {
				b.Fatalf("Failed to calculate statistics: %v", err)
			}
		}
	})
}

func TestStartupTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping startup time test in short mode")
	}

	tempDir := t.TempDir()
	_ = os.Setenv("XDG_CONFIG_HOME", tempDir)
	_ = os.Setenv("XDG_CACHE_HOME", tempDir)
	_ = os.Setenv("XDG_DATA_HOME", tempDir)

	start := time.Now()

	// Initialize all components like main does
	_, err := config.New()
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	_, err = service.NewDictationService(false)
	if err != nil {
		t.Fatalf("Failed to create dictation service: %v", err)
	}

	_, err = storage.NewAudioPlayer()
	if err != nil {
		t.Fatalf("Failed to create audio player: %v", err)
	}

	_, err = storage.NewHistory()
	if err != nil {
		t.Fatalf("Failed to create history: %v", err)
	}

	startupTime := time.Since(start)
	t.Logf("Startup time: %v", startupTime)

	// Should start up in less than 500ms
	if startupTime > 500*time.Millisecond {
		t.Errorf("Startup time too slow: %v (expected < 500ms)", startupTime)
	}
}

func TestResponseTimes(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping response time test")
	}

	tempDir := t.TempDir()
	_ = os.Setenv("XDG_CACHE_HOME", tempDir)

	dictationSvc, err := service.NewDictationService(false)
	if err != nil {
		t.Fatalf("Failed to create dictation service: %v", err)
	}

	cfg := &types.Config{
		Voice: "alloy",
		Level: 600,
		Topic: "Daily",
		Words: 10,
		Speed: 1.0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test sentence generation time
	start := time.Now()
	sentence, err := dictationSvc.GenerateSentence(ctx, cfg.Topic, cfg.Level, cfg.Words)
	if err != nil {
		t.Fatalf("Failed to generate sentence: %v", err)
	}
	sentenceTime := time.Since(start)
	t.Logf("Sentence generation time: %v", sentenceTime)

	// Test audio generation time
	start = time.Now()
	_, err = dictationSvc.GenerateAudio(ctx, sentence, cfg.Voice, cfg.Speed)
	if err != nil {
		t.Fatalf("Failed to generate audio: %v", err)
	}
	audioTime := time.Since(start)
	t.Logf("Audio generation time: %v", audioTime)

	// Test grading time
	start = time.Now()
	_, err = dictationSvc.GradeDictation(ctx, sentence, sentence)
	if err != nil {
		t.Fatalf("Failed to grade dictation: %v", err)
	}
	gradingTime := time.Since(start)
	t.Logf("Grading time: %v", gradingTime)

	// Check reasonable response times
	if sentenceTime > 10*time.Second {
		t.Errorf("Sentence generation too slow: %v", sentenceTime)
	}
	if audioTime > 15*time.Second {
		t.Errorf("Audio generation too slow: %v", audioTime)
	}
	if gradingTime > 10*time.Second {
		t.Errorf("Grading too slow: %v", gradingTime)
	}
}

func TestMemoryUsage(t *testing.T) {
	var m1, m2 runtime.MemStats

	// Measure initial memory
	runtime.GC()
	runtime.ReadMemStats(&m1)

	tempDir := t.TempDir()
	_ = os.Setenv("XDG_CONFIG_HOME", tempDir)
	_ = os.Setenv("XDG_CACHE_HOME", tempDir)
	_ = os.Setenv("XDG_DATA_HOME", tempDir)

	// Initialize services
	_, err := config.New()
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	_, err = service.NewDictationService(false)
	if err != nil {
		t.Fatalf("Failed to create dictation service: %v", err)
	}

	_, err = storage.NewAudioPlayer()
	if err != nil {
		t.Fatalf("Failed to create audio player: %v", err)
	}

	history, err := storage.NewHistory()
	if err != nil {
		t.Fatalf("Failed to create history: %v", err)
	}

	cache, err := storage.NewAudioCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Simulate some usage
	testData := make([]byte, 100*1024) // 100KB
	for i := 0; i < 10; i++ {
		_ = cache.Save("test", "alloy", float64(i), testData)

		session := &types.DictationSession{
			ID:        fmt.Sprintf("test-%d", i),
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
			Sentence:     "Test sentence",
			AudioPath:    "/tmp/test.mp3",
			UserInput:    "Test sentence",
			Grade: &types.Grade{
				Score: 100,
				WER:   0.0,
			},
			ReplayCount:  0,
			DurationSecs: 60.0,
		}
		_ = history.SaveSession(session)
	}

	// Measure final memory
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Handle potential overflow in memory calculation
	var memUsed float64
	if m2.Alloc >= m1.Alloc {
		memUsed = float64(m2.Alloc-m1.Alloc) / 1024 / 1024 // MB
	} else {
		memUsed = 0 // Handle wraparound
	}
	t.Logf("Memory usage: %.2f MB", memUsed)

	// Should use less than 50MB for basic operations
	if memUsed > 50 {
		t.Errorf("Memory usage too high: %.2f MB (expected < 50MB)", memUsed)
	}
}

func TestConcurrentOperations(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.Setenv("XDG_CACHE_HOME", tempDir)
	_ = os.Setenv("XDG_DATA_HOME", tempDir)

	cache, err := storage.NewAudioCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	history, err := storage.NewHistory()
	if err != nil {
		t.Fatalf("Failed to create history: %v", err)
	}

	// Clear any existing history first
	_ = history.Clear()

	// Test concurrent cache operations
	t.Run("ConcurrentCache", func(t *testing.T) {
		done := make(chan bool, 10)
		testData := []byte("test data")

		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Save
				err := cache.Save("test", "alloy", float64(id), testData)
				if err != nil {
					t.Errorf("Failed to save in goroutine %d: %v", id, err)
					return
				}

				// Load
				_, err = cache.Load("test", "alloy", float64(id))
				if err != nil {
					t.Errorf("Failed to load in goroutine %d: %v", id, err)
					return
				}
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	// Test concurrent history operations
	t.Run("ConcurrentHistory", func(t *testing.T) {
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				session := &types.DictationSession{
					ID:        fmt.Sprintf("concurrent-%d", id),
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
					Sentence:     "Test sentence",
					AudioPath:    "/tmp/test.mp3",
					UserInput:    "Test sentence",
					Grade: &types.Grade{
						Score: 100,
						WER:   0.0,
					},
					ReplayCount:  0,
					DurationSecs: 60.0,
				}

				err := history.SaveSession(session)
				if err != nil {
					t.Errorf("Failed to save session in goroutine %d: %v", id, err)
					return
				}
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}

		// Verify all sessions were saved
		stats, err := history.CalculateStatistics(30)
		if err != nil {
			t.Fatalf("Failed to calculate statistics: %v", err)
		}

		if stats.TotalSessions != 10 {
			t.Errorf("Expected 10 sessions, got: %d", stats.TotalSessions)
		}
	})
}