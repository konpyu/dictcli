package types

import (
	"time"
)

// DictationSession represents a single dictation practice session
type DictationSession struct {
	// Session metadata
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
	ConfigUsed  Config    `json:"config_used"`

	// Generated content
	Sentence     string `json:"sentence"`
	AudioPath    string `json:"audio_path"`
	AudioCached  bool   `json:"audio_cached"`
	GenerationMS int64  `json:"generation_ms"`

	// User interaction
	UserInput   string `json:"user_input"`
	TimeTakenMS int64  `json:"time_taken_ms"`
	ReplayCount int    `json:"replay_count"`

	// Grading result
	Grade         *Grade `json:"grade,omitempty"`
	GradingTimeMS int64  `json:"grading_time_ms,omitempty"`

	// Session state
	Completed bool      `json:"completed"`
	EndTime   time.Time `json:"end_time,omitempty"`
}

// Duration returns the total session duration
func (s *DictationSession) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.Timestamp)
	}
	return s.EndTime.Sub(s.Timestamp)
}

// IsCompleted checks if the session was completed successfully
func (s *DictationSession) IsCompleted() bool {
	return s.Completed && s.Grade != nil
}