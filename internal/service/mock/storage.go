package mock

import (
	"github.com/yourusername/dictcli/internal/logger"
	"github.com/yourusername/dictcli/internal/storage"
	"github.com/yourusername/dictcli/internal/types"
)

// MockStorage provides a mock implementation of the Storage interface
type MockStorage struct {
	sessions []*types.DictationSession
	logger   *logger.Logger
}

// NewMockStorage creates a new mock storage
func NewMockStorage(logger *logger.Logger) *MockStorage {
	return &MockStorage{
		sessions: make([]*types.DictationSession, 0),
		logger:   logger,
	}
}

// SaveSession saves a completed dictation session
func (m *MockStorage) SaveSession(session *types.DictationSession) error {
	m.sessions = append(m.sessions, session)
	m.logger.Info("MockStorage: Saved session - id: %s, score: %d", session.ID, session.Grade.Score)
	return nil
}

// GetSessions retrieves sessions based on filters
func (m *MockStorage) GetSessions(filters storage.SessionFilters) ([]*types.DictationSession, error) {
	// Simple implementation - return all sessions for now
	// TODO: Apply actual filtering
	m.logger.Debug("MockStorage: Retrieved sessions - count: %d", len(m.sessions))
	return m.sessions, nil
}

// GetStatistics calculates statistics for a given time period
func (m *MockStorage) GetStatistics(period storage.StatsPeriod) (*storage.Statistics, error) {
	if len(m.sessions) == 0 {
		return &storage.Statistics{}, nil
	}
	
	totalScore := 0.0
	totalWER := 0.0
	
	for _, session := range m.sessions {
		if session.Grade != nil {
			totalScore += float64(session.Grade.Score)
			totalWER += session.Grade.WER
		}
	}
	
	stats := &storage.Statistics{
		TotalSessions:   len(m.sessions),
		AverageScore:    totalScore / float64(len(m.sessions)),
		AverageWER:      totalWER / float64(len(m.sessions)),
		TotalTime:       len(m.sessions) * 60, // Mock: 1 minute per session
		CommonMistakes:  []storage.MistakeCount{},
		ProgressByTopic: make(map[string]storage.TopicProgress),
	}
	
	m.logger.Debug("MockStorage: Generated statistics - sessions: %d, avg_score: %.2f", stats.TotalSessions, stats.AverageScore)
	
	return stats, nil
}

// ClearHistory removes all stored sessions
func (m *MockStorage) ClearHistory() error {
	count := len(m.sessions)
	m.sessions = make([]*types.DictationSession, 0)
	m.logger.Info("MockStorage: Cleared history - sessions removed: %d", count)
	return nil
}