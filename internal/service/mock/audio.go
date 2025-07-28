package mock

import (
	"context"
	"fmt"
	"sync"
	"time"
	
	"github.com/yourusername/dictcli/internal/logger"
)

// MockAudioPlayer simulates audio playback without actually playing audio
type MockAudioPlayer struct {
	mu       sync.Mutex
	playing  bool
	stopChan chan struct{}
	logger   *logger.Logger
}

// NewMockAudioPlayer creates a new mock audio player
func NewMockAudioPlayer(logger *logger.Logger) *MockAudioPlayer {
	return &MockAudioPlayer{
		logger: logger,
	}
}

// Play simulates playing an audio file
func (m *MockAudioPlayer) Play(ctx context.Context, audioPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.playing {
		return fmt.Errorf("audio is already playing")
	}
	
	m.logger.Info("MockAudioPlayer: Starting playback: %s", audioPath)
	
	m.playing = true
	m.stopChan = make(chan struct{})
	
	// Simulate audio playback in a goroutine
	go func() {
		defer func() {
			m.mu.Lock()
			m.playing = false
			m.mu.Unlock()
			m.logger.Info("MockAudioPlayer: Playback finished: %s", audioPath)
		}()
		
		// Simulate realistic audio duration (2-8 seconds based on sentence length)
		duration := m.estimateAudioDuration(audioPath)
		
		select {
		case <-time.After(duration):
			// Normal completion
		case <-m.stopChan:
			// Stopped manually
			m.logger.Info("MockAudioPlayer: Playback stopped: %s", audioPath)
		case <-ctx.Done():
			// Context cancelled
			m.logger.Info("MockAudioPlayer: Playback cancelled: %s", audioPath)
		}
	}()
	
	return nil
}

// Stop stops the currently playing audio
func (m *MockAudioPlayer) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if !m.playing {
		return fmt.Errorf("no audio is currently playing")
	}
	
	close(m.stopChan)
	return nil
}

// IsPlaying returns whether audio is currently playing
func (m *MockAudioPlayer) IsPlaying() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.playing
}

// estimateAudioDuration estimates audio duration based on filename/path
func (m *MockAudioPlayer) estimateAudioDuration(audioPath string) time.Duration {
	// Simulate realistic durations based on mock file timestamps
	// In a real implementation, this would be based on actual audio length
	
	// Add some variation based on the path
	variation := time.Duration(len(audioPath)%3) * time.Second
	
	return baseTimeout + variation
}

// MockAudioCache simulates audio file caching
type MockAudioCache struct {
	cache  map[string]string // key -> file path
	mu     sync.RWMutex
	logger *logger.Logger
}

// NewMockAudioCache creates a new mock audio cache
func NewMockAudioCache(logger *logger.Logger) *MockAudioCache {
	return &MockAudioCache{
		cache:  make(map[string]string),
		logger: logger,
	}
}

// Get retrieves a cached audio file if it exists
func (m *MockAudioCache) Get(text string, voice string, speed float64) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	key := m.generateCacheKey(text, voice, speed)
	path, exists := m.cache[key]
	
	if exists {
		m.logger.Debug("MockAudioCache: Cache hit - key: %s, path: %s", key, path)
	} else {
		m.logger.Debug("MockAudioCache: Cache miss - key: %s", key)
	}
	
	return path, exists
}

// Put stores an audio file in the cache
func (m *MockAudioCache) Put(text string, voice string, speed float64, audioPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	key := m.generateCacheKey(text, voice, speed)
	m.cache[key] = audioPath
	
	m.logger.Debug("MockAudioCache: Cached audio - key: %s, path: %s", key, audioPath)
	
	return nil
}

// Clear removes all cached audio files
func (m *MockAudioCache) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	count := len(m.cache)
	m.cache = make(map[string]string)
	
	m.logger.Info("MockAudioCache: Cleared cache - files removed: %d", count)
	
	return nil
}

// Size returns the total size of cached files in bytes (simulated)
func (m *MockAudioCache) Size() (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Simulate cache size (approx 50KB per file)
	return int64(len(m.cache) * 50 * 1024), nil
}

// generateCacheKey creates a cache key from audio parameters
func (m *MockAudioCache) generateCacheKey(text string, voice string, speed float64) string {
	return fmt.Sprintf("%s_%s_%.1f", hashString(text), voice, speed)
}

// hashString creates a simple hash of a string (for mock purposes)
func hashString(s string) string {
	hash := uint32(0)
	for _, c := range s {
		hash = hash*31 + uint32(c)
	}
	return fmt.Sprintf("%08x", hash)
}

// baseTimeout is the base duration for audio playback simulation
const baseTimeout = 3 * time.Second