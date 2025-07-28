package mock

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
	
	"github.com/yourusername/dictcli/internal/types"
)

// nolint:gosec // Using math/rand for mock data generation is acceptable
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// MockDictationService provides mock implementations for testing the TUI
type MockDictationService struct {
	// Configurable delays to simulate API latency
	GenerationDelay time.Duration
	AudioDelay      time.Duration
	GradingDelay    time.Duration
	
	// Test data
	sentences map[string][]string // sentences by topic
	gradeTemplates []GradeTemplate
}

// GradeTemplate defines a template for generating mock grades
type GradeTemplate struct {
	WER                     float64
	Score                   int
	JapaneseExplanation     string
	AlternativeExpressions  []string
	SampleMistakes          []types.Mistake
}

// NewMockDictationService creates a new mock service with predefined test data
func NewMockDictationService() *MockDictationService {
	return &MockDictationService{
		GenerationDelay: 800 * time.Millisecond,
		AudioDelay:      600 * time.Millisecond,
		GradingDelay:    1200 * time.Millisecond,
		sentences:       getSampleSentences(),
		gradeTemplates:  getGradeTemplates(),
	}
}

// GenerateSentence returns a predefined sentence based on configuration
func (m *MockDictationService) GenerateSentence(ctx context.Context, config *types.Config) (string, error) {
	// Simulate API delay
	select {
	case <-time.After(m.GenerationDelay):
	case <-ctx.Done():
		return "", ctx.Err()
	}
	
	// Get sentences for the topic
	sentences, exists := m.sentences[config.Topic]
	if !exists {
		sentences = m.sentences[types.TopicBusiness] // fallback
	}
	
	// Filter by approximate word count (±3 words)
	var suitable []string
	for _, sentence := range sentences {
		wordCount := len(splitWords(sentence))
		if abs(wordCount-config.WordCount) <= 3 {
			suitable = append(suitable, sentence)
		}
	}
	
	if len(suitable) == 0 {
		suitable = sentences // fallback to any sentence
	}
	
	// Return random suitable sentence
	// nolint:gosec // Using math/rand for mock data generation is acceptable
	return suitable[rng.Intn(len(suitable))], nil
}

// GenerateAudio returns a mock audio file path
func (m *MockDictationService) GenerateAudio(ctx context.Context, text string, config *types.Config) (string, error) {
	// Simulate API delay
	select {
	case <-time.After(m.AudioDelay):
	case <-ctx.Done():
		return "", ctx.Err()
	}
	
	// Generate a mock file path (doesn't actually exist)
	filename := fmt.Sprintf("mock_audio_%d.mp3", time.Now().Unix())
	audioPath := filepath.Join(config.CachePath, filename)
	
	return audioPath, nil
}

// GradeDictation evaluates user input and returns grading results
func (m *MockDictationService) GradeDictation(ctx context.Context, correct, userInput string, config *types.Config) (*types.Grade, error) {
	// Simulate API delay
	select {
	case <-time.After(m.GradingDelay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	// Simple mock grading logic
	correctWords := splitWords(correct)
	userWords := splitWords(userInput)
	
	// Calculate basic metrics
	totalWords := len(correctWords)
	if totalWords == 0 {
		return &types.Grade{Score: 0, WER: 1.0}, nil
	}
	
	// Simple word-by-word comparison
	correctCount := 0
	var mistakes []types.Mistake
	
	maxLen := max(len(correctWords), len(userWords))
	for i := 0; i < maxLen; i++ {
		if i < len(correctWords) && i < len(userWords) {
			if correctWords[i] == userWords[i] {
				correctCount++
			} else {
				// Substitution
				mistakes = append(mistakes, types.Mistake{
					Position: i,
					Expected: correctWords[i],
					Actual:   userWords[i],
					Type:     types.MistakeTypeSubstitution,
					JapaneseNote: getJapaneseMistakeNote(correctWords[i], userWords[i]),
				})
			}
		} else if i >= len(userWords) {
			// Deletion (user missing word)
			mistakes = append(mistakes, types.Mistake{
				Position: i,
				Expected: correctWords[i],
				Actual:   "",
				Type:     types.MistakeTypeDeletion,
				JapaneseNote: fmt.Sprintf("「%s」が不足しています", correctWords[i]),
			})
		} else {
			// Insertion (user added extra word)
			mistakes = append(mistakes, types.Mistake{
				Position: i,
				Expected: "",
				Actual:   userWords[i],
				Type:     types.MistakeTypeInsertion,
				JapaneseNote: fmt.Sprintf("「%s」は不要です", userWords[i]),
			})
		}
	}
	
	// Calculate WER and score
	errorCount := len(mistakes)
	wer := float64(errorCount) / float64(totalWords)
	if wer > 1.0 {
		wer = 1.0
	}
	
	score := int((1.0 - wer) * 100)
	if score < 0 {
		score = 0
	}
	
	// Select appropriate grade template for Japanese explanation
	template := m.selectGradeTemplate(wer)
	
	grade := &types.Grade{
		WER:                    wer,
		Score:                  score,
		Mistakes:               mistakes,
		JapaneseExplanation:    template.JapaneseExplanation,
		AlternativeExpressions: template.AlternativeExpressions,
		CorrectAnswer:          correct,
		WordCount:              totalWords,
		CorrectWords:           correctCount,
		IncorrectWords:         errorCount,
	}
	
	return grade, nil
}

// Helper functions

func splitWords(text string) []string {
	// Simple word splitting (could be improved)
	words := make([]string, 0)
	current := ""
	
	for _, r := range text {
		if r == ' ' || r == '\t' || r == '\n' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	
	if current != "" {
		words = append(words, current)
	}
	
	return words
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m *MockDictationService) selectGradeTemplate(wer float64) GradeTemplate {
	// Select template based on WER
	for _, template := range m.gradeTemplates {
		if wer <= template.WER {
			return template
		}
	}
	// Fallback to last template
	return m.gradeTemplates[len(m.gradeTemplates)-1]
}

func getJapaneseMistakeNote(expected, actual string) string {
	commonMistakes := map[string]map[string]string{
		"a": {
			"an": "冠詞の使い分けに注意しましょう",
			"the": "不定冠詞と定冠詞の違いです",
		},
		"is": {
			"are": "単数・複数の主語に応じた動詞の活用です",
			"was": "現在形と過去形の使い分けです",
		},
		"their": {
			"there": "所有格代名詞と副詞の違いです",
			"they're": "所有格と短縮形の違いです",
		},
	}
	
	if expectedMap, exists := commonMistakes[expected]; exists {
		if note, exists := expectedMap[actual]; exists {
			return note
		}
	}
	
	return fmt.Sprintf("「%s」→「%s」の聞き取りミス", expected, actual)
}