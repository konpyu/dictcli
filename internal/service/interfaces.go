// Package service defines the core service interfaces for DictCLI
package service

import (
	"context"
	
	"github.com/yourusername/dictcli/internal/types"
)

// DictationService handles sentence generation, audio creation, and grading
type DictationService interface {
	// GenerateSentence creates an English sentence based on the given configuration
	GenerateSentence(ctx context.Context, config *types.Config) (string, error)
	
	// GenerateAudio converts text to speech and returns the audio file path
	GenerateAudio(ctx context.Context, text string, config *types.Config) (string, error)
	
	// GradeDictation evaluates user input against the correct answer
	GradeDictation(ctx context.Context, correct, userInput string, config *types.Config) (*types.Grade, error)
}

// AudioPlayer handles audio playback
type AudioPlayer interface {
	// Play plays the audio file at the given path
	Play(ctx context.Context, audioPath string) error
	
	// Stop stops the currently playing audio
	Stop() error
	
	// IsPlaying returns whether audio is currently playing
	IsPlaying() bool
}

// AudioCache manages audio file caching
type AudioCache interface {
	// Get retrieves a cached audio file if it exists
	Get(text string, voice string, speed float64) (string, bool)
	
	// Put stores an audio file in the cache
	Put(text string, voice string, speed float64, audioPath string) error
	
	// Clear removes all cached audio files
	Clear() error
	
	// Size returns the total size of cached files in bytes
	Size() (int64, error)
}

// Storage handles data persistence
type Storage interface {
	// SaveSession appends a session to the history
	SaveSession(session *types.DictationSession) error
	
	// GetSessions retrieves sessions within a time range
	GetSessions(from, to string) ([]*types.DictationSession, error)
	
	// GetStatistics calculates statistics for a given period
	GetStatistics(days int) (*Statistics, error)
	
	// GetCommonMistakes returns the most common mistakes
	GetCommonMistakes(limit int) ([]MistakePattern, error)
}

// Statistics represents aggregated learning statistics
type Statistics struct {
	TotalSessions    int     `json:"total_sessions"`
	TotalRounds      int     `json:"total_rounds"`
	AverageScore     float64 `json:"average_score"`
	AverageWER       float64 `json:"average_wer"`
	TotalTimeMinutes int     `json:"total_time_minutes"`
	
	// Progress by topic
	TopicStats map[string]*TopicStatistics `json:"topic_stats"`
	
	// Progress by level
	LevelStats map[int]*LevelStatistics `json:"level_stats"`
	
	// Time-based progress
	DailyProgress []DailyStats `json:"daily_progress"`
}

// TopicStatistics represents statistics for a specific topic
type TopicStatistics struct {
	Sessions     int     `json:"sessions"`
	AverageScore float64 `json:"average_score"`
	AverageWER   float64 `json:"average_wer"`
}

// LevelStatistics represents statistics for a specific level
type LevelStatistics struct {
	Sessions     int     `json:"sessions"`
	AverageScore float64 `json:"average_score"`
	AverageWER   float64 `json:"average_wer"`
}

// DailyStats represents statistics for a single day
type DailyStats struct {
	Date         string  `json:"date"`
	Sessions     int     `json:"sessions"`
	AverageScore float64 `json:"average_score"`
	AverageWER   float64 `json:"average_wer"`
}

// MistakePattern represents a common mistake pattern
type MistakePattern struct {
	Expected  string `json:"expected"`
	Actual    string `json:"actual"`
	Count     int    `json:"count"`
	Frequency float64 `json:"frequency"` // Percentage of occurrences
}