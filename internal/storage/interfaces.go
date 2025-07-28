package storage

import "github.com/yourusername/dictcli/internal/types"

// Storage defines the interface for persisting dictation data
type Storage interface {
	// SaveSession saves a completed dictation session
	SaveSession(session *types.DictationSession) error
	
	// GetSessions retrieves sessions based on filters
	GetSessions(filters SessionFilters) ([]*types.DictationSession, error)
	
	// GetStatistics calculates statistics for a given time period
	GetStatistics(period StatsPeriod) (*Statistics, error)
	
	// ClearHistory removes all stored sessions
	ClearHistory() error
}

// SessionFilters defines filters for retrieving sessions
type SessionFilters struct {
	StartDate    *string
	EndDate      *string
	Topic        *string
	MinScore     *int
	MaxScore     *int
	Limit        int
}

// StatsPeriod defines the time period for statistics
type StatsPeriod struct {
	Days int
}

// Statistics contains aggregated statistics
type Statistics struct {
	TotalSessions   int
	AverageScore    float64
	AverageWER      float64
	TotalTime       int // in seconds
	CommonMistakes  []MistakeCount
	ProgressByTopic map[string]TopicProgress
}

// MistakeCount represents a common mistake and its frequency
type MistakeCount struct {
	Expected string
	Actual   string
	Count    int
}

// TopicProgress represents progress for a specific topic
type TopicProgress struct {
	Sessions     int
	AverageScore float64
	AverageWER   float64
}