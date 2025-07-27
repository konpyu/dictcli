package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/konpyu/dictcli/internal/types"
)

type History struct {
	filePath string
}

func NewHistory() (*History, error) {
	dataDir := filepath.Join(xdg.DataHome, "dictcli")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &History{
		filePath: filepath.Join(dataDir, "history.jsonl"),
	}, nil
}

func (h *History) SaveSession(session *types.DictationSession) error {
	file, err := os.OpenFile(h.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(session); err != nil {
		return fmt.Errorf("failed to encode session: %w", err)
	}

	return nil
}

func (h *History) LoadSessions(days int) ([]*types.DictationSession, error) {
	file, err := os.Open(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*types.DictationSession{}, nil
		}
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var sessions []*types.DictationSession
	cutoffTime := time.Now().AddDate(0, 0, -days)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var session types.DictationSession
		if err := json.Unmarshal(scanner.Bytes(), &session); err != nil {
			continue
		}

		if session.Timestamp.After(cutoffTime) {
			sessions = append(sessions, &session)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	return sessions, nil
}

func (h *History) LoadAllSessions() ([]*types.DictationSession, error) {
	file, err := os.Open(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*types.DictationSession{}, nil
		}
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var sessions []*types.DictationSession
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var session types.DictationSession
		if err := json.Unmarshal(scanner.Bytes(), &session); err != nil {
			continue
		}
		sessions = append(sessions, &session)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	return sessions, nil
}

type Statistics struct {
	TotalSessions   int
	TotalRounds     int
	AverageScore    float64
	AverageWER      float64
	TopicBreakdown  map[string]*TopicStats
	CommonMistakes  []MistakeFrequency
	RecentProgress  []DailyStats
}

type TopicStats struct {
	Count        int
	AverageScore float64
	AverageWER   float64
}

type MistakeFrequency struct {
	Expected  string
	Actual    string
	Frequency int
}

type DailyStats struct {
	Date         time.Time
	SessionCount int
	AverageScore float64
	AverageWER   float64
}

func (h *History) CalculateStatistics(days int) (*Statistics, error) {
	sessions, err := h.LoadSessions(days)
	if err != nil {
		return nil, err
	}

	stats := &Statistics{
		TopicBreakdown: make(map[string]*TopicStats),
		CommonMistakes: []MistakeFrequency{},
		RecentProgress: []DailyStats{},
	}

	if len(sessions) == 0 {
		return stats, nil
	}

	mistakeMap := make(map[string]int)
	dailyMap := make(map[string]*DailyStats)
	var totalScore, totalWER float64

	for _, session := range sessions {
		stats.TotalSessions++
		
		if session.Grade != nil {
			totalScore += float64(session.Grade.Score)
			totalWER += session.Grade.WER

			if topicStats, exists := stats.TopicBreakdown[session.Config.Topic]; exists {
				topicStats.Count++
				topicStats.AverageScore = ((topicStats.AverageScore * float64(topicStats.Count-1)) + float64(session.Grade.Score)) / float64(topicStats.Count)
				topicStats.AverageWER = ((topicStats.AverageWER * float64(topicStats.Count-1)) + session.Grade.WER) / float64(topicStats.Count)
			} else {
				stats.TopicBreakdown[session.Config.Topic] = &TopicStats{
					Count:        1,
					AverageScore: float64(session.Grade.Score),
					AverageWER:   session.Grade.WER,
				}
			}

			for _, mistake := range session.Grade.Mistakes {
				key := fmt.Sprintf("%s->%s", mistake.Expected, mistake.Actual)
				mistakeMap[key]++
			}

			dateKey := session.Timestamp.Format("2006-01-02")
			if daily, exists := dailyMap[dateKey]; exists {
				daily.SessionCount++
				daily.AverageScore = ((daily.AverageScore * float64(daily.SessionCount-1)) + float64(session.Grade.Score)) / float64(daily.SessionCount)
				daily.AverageWER = ((daily.AverageWER * float64(daily.SessionCount-1)) + session.Grade.WER) / float64(daily.SessionCount)
			} else {
				dailyMap[dateKey] = &DailyStats{
					Date:         session.Timestamp.Truncate(24 * time.Hour),
					SessionCount: 1,
					AverageScore: float64(session.Grade.Score),
					AverageWER:   session.Grade.WER,
				}
			}
		}
	}

	if stats.TotalSessions > 0 {
		stats.AverageScore = totalScore / float64(stats.TotalSessions)
		stats.AverageWER = totalWER / float64(stats.TotalSessions)
	}

	for key, count := range mistakeMap {
		if count >= 2 {
			parts := splitMistakeKey(key)
			if len(parts) == 2 {
				stats.CommonMistakes = append(stats.CommonMistakes, MistakeFrequency{
					Expected:  parts[0],
					Actual:    parts[1],
					Frequency: count,
				})
			}
		}
	}

	for _, daily := range dailyMap {
		stats.RecentProgress = append(stats.RecentProgress, *daily)
	}

	return stats, nil
}

func splitMistakeKey(key string) []string {
	for i := 0; i < len(key)-1; i++ {
		if key[i:i+2] == "->" {
			return []string{key[:i], key[i+2:]}
		}
	}
	return []string{}
}

func (h *History) Clear() error {
	return os.Remove(h.filePath)
}