package storage

import (
	"testing"
	"time"

	"github.com/konpyu/dictcli/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatisticsCalculation(t *testing.T) {
	t.Run("Calculate statistics with multiple sessions", func(t *testing.T) {
		history, err := NewHistory()
		require.NoError(t, err)
		// Clear any existing history data first
		_ = history.Clear()
		defer func() { _ = history.Clear() }()

		// Create test sessions
		now := time.Now()
		sessions := []*types.DictationSession{
			{
				ID:           "test1",
				Timestamp:    now.Add(-2 * 24 * time.Hour),
				Config:       types.Config{Topic: "Business", Level: 700},
				Sentence:     "Test sentence 1",
				UserInput:    "Test sentence 1",
				StartTime:    now.Add(-2 * 24 * time.Hour),
				EndTime:      now.Add(-2 * 24 * time.Hour).Add(30 * time.Second),
				DurationSecs: 30.0,
				Grade: &types.Grade{
					WER:   0.0,
					Score: 100,
					Mistakes: []types.Mistake{},
				},
			},
			{
				ID:           "test2",
				Timestamp:    now.Add(-1 * 24 * time.Hour),
				Config:       types.Config{Topic: "Business", Level: 700},
				Sentence:     "Test sentence 2",
				UserInput:    "Test centence 2", // Intentional mistake
				StartTime:    now.Add(-1 * 24 * time.Hour),
				EndTime:      now.Add(-1 * 24 * time.Hour).Add(25 * time.Second),
				DurationSecs: 25.0,
				Grade: &types.Grade{
					WER:   0.333,
					Score: 67,
					Mistakes: []types.Mistake{
						{Position: 1, Expected: "sentence", Actual: "centence", Type: "substitution"},
					},
				},
			},
			{
				ID:           "test3",
				Timestamp:    now,
				Config:       types.Config{Topic: "Travel", Level: 600},
				Sentence:     "I will go to the beach",
				UserInput:    "I will go the beach", // Missing "to"
				StartTime:    now,
				EndTime:      now.Add(20 * time.Second),
				DurationSecs: 20.0,
				Grade: &types.Grade{
					WER:   0.167,
					Score: 83,
					Mistakes: []types.Mistake{
						{Position: 3, Expected: "to", Actual: "", Type: "deletion"},
					},
				},
			},
		}

		// Save sessions
		for _, session := range sessions {
			err := history.SaveSession(session)
			require.NoError(t, err)
		}

		// Calculate statistics for last 7 days
		stats, err := history.CalculateStatistics(7)
		require.NoError(t, err)

		// Verify overall statistics
		assert.Equal(t, 3, stats.TotalSessions)
		assert.Equal(t, 3, stats.TotalRounds)
		assert.InDelta(t, 83.33, stats.AverageScore, 0.01)
		assert.InDelta(t, 0.167, stats.AverageWER, 0.001)

		// Verify topic breakdown
		assert.Len(t, stats.TopicBreakdown, 2)
		assert.Equal(t, 2, stats.TopicBreakdown["Business"].Count)
		assert.InDelta(t, 83.5, stats.TopicBreakdown["Business"].AverageScore, 0.01)
		assert.Equal(t, 1, stats.TopicBreakdown["Travel"].Count)
		assert.InDelta(t, 83.0, stats.TopicBreakdown["Travel"].AverageScore, 0.01)

		// Verify common mistakes (none should appear as each mistake only occurs once)
		assert.Len(t, stats.CommonMistakes, 0)

		// Verify recent progress
		assert.Len(t, stats.RecentProgress, 3)
		// Should be sorted by date (oldest first)
		assert.True(t, stats.RecentProgress[0].Date.Before(stats.RecentProgress[1].Date))
		assert.True(t, stats.RecentProgress[1].Date.Before(stats.RecentProgress[2].Date))
	})

	t.Run("Calculate statistics with common mistakes", func(t *testing.T) {
		history, err := NewHistory()
		require.NoError(t, err)
		// Clear any existing history data first
		_ = history.Clear()
		defer func() { _ = history.Clear() }()

		// Create sessions with repeated mistakes
		now := time.Now()
		sessions := []*types.DictationSession{
			{
				ID:        "test1",
				Timestamp: now,
				Config:    types.Config{Topic: "Business"},
				Grade: &types.Grade{
					WER:   0.2,
					Score: 80,
					Mistakes: []types.Mistake{
						{Expected: "the", Actual: "teh", Type: "substitution"},
						{Expected: "because", Actual: "becuase", Type: "substitution"},
					},
				},
			},
			{
				ID:        "test2",
				Timestamp: now,
				Config:    types.Config{Topic: "Business"},
				Grade: &types.Grade{
					WER:   0.1,
					Score: 90,
					Mistakes: []types.Mistake{
						{Expected: "the", Actual: "teh", Type: "substitution"},
						{Expected: "because", Actual: "becuase", Type: "substitution"},
					},
				},
			},
			{
				ID:        "test3",
				Timestamp: now,
				Config:    types.Config{Topic: "Business"},
				Grade: &types.Grade{
					WER:   0.15,
					Score: 85,
					Mistakes: []types.Mistake{
						{Expected: "the", Actual: "teh", Type: "substitution"},
					},
				},
			},
		}

		// Save sessions
		for _, session := range sessions {
			err := history.SaveSession(session)
			require.NoError(t, err)
		}

		// Calculate statistics
		stats, err := history.CalculateStatistics(30)
		require.NoError(t, err)

		// Verify common mistakes
		assert.Len(t, stats.CommonMistakes, 2)
		// Should be sorted by frequency (descending)
		assert.Equal(t, "the", stats.CommonMistakes[0].Expected)
		assert.Equal(t, "teh", stats.CommonMistakes[0].Actual)
		assert.Equal(t, 3, stats.CommonMistakes[0].Frequency)
		assert.Equal(t, "because", stats.CommonMistakes[1].Expected)
		assert.Equal(t, "becuase", stats.CommonMistakes[1].Actual)
		assert.Equal(t, 2, stats.CommonMistakes[1].Frequency)
	})

	t.Run("Calculate statistics with date filtering", func(t *testing.T) {
		history, err := NewHistory()
		require.NoError(t, err)
		// Clear any existing history data first
		_ = history.Clear()
		defer func() { _ = history.Clear() }()

		// Create sessions across different days
		now := time.Now()
		sessions := []*types.DictationSession{
			{
				ID:        "old1",
				Timestamp: now.Add(-10 * 24 * time.Hour), // 10 days ago
				Config:    types.Config{Topic: "Business"},
				Grade:     &types.Grade{WER: 0.1, Score: 90},
			},
			{
				ID:        "recent1",
				Timestamp: now.Add(-2 * 24 * time.Hour), // 2 days ago
				Config:    types.Config{Topic: "Business"},
				Grade:     &types.Grade{WER: 0.2, Score: 80},
			},
			{
				ID:        "recent2",
				Timestamp: now, // Today
				Config:    types.Config{Topic: "Business"},
				Grade:     &types.Grade{WER: 0.15, Score: 85},
			},
		}

		// Save sessions
		for _, session := range sessions {
			err := history.SaveSession(session)
			require.NoError(t, err)
		}

		// Calculate statistics for last 5 days (should exclude the old session)
		stats, err := history.CalculateStatistics(5)
		require.NoError(t, err)

		assert.Equal(t, 2, stats.TotalSessions)
		assert.InDelta(t, 82.5, stats.AverageScore, 0.01)
	})
}