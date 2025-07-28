package types

// Grade represents the grading result of a dictation attempt
type Grade struct {
	// Core metrics
	WER   float64 `json:"wer"`   // Word Error Rate (0.0-1.0)
	Score int     `json:"score"` // Score out of 100

	// Error analysis
	Mistakes []Mistake `json:"mistakes"`
	
	// Japanese feedback
	JapaneseExplanation string   `json:"japanese_explanation"`
	AlternativeExpressions []string `json:"alternative_expressions"`
	
	// Correct answer
	CorrectAnswer string `json:"correct_answer"`
	
	// Additional metrics
	WordCount      int `json:"word_count"`
	CorrectWords   int `json:"correct_words"`
	IncorrectWords int `json:"incorrect_words"`
}

// Mistake represents a single error in the dictation
type Mistake struct {
	Position int    `json:"position"` // Word position (0-based)
	Expected string `json:"expected"` // What was expected
	Actual   string `json:"actual"`   // What user typed
	Type     string `json:"type"`     // Error type: "insertion", "deletion", "substitution"
	
	// Japanese explanation for this specific mistake
	JapaneseNote string `json:"japanese_note,omitempty"`
}

// MistakeType constants
const (
	MistakeTypeInsertion    = "insertion"
	MistakeTypeDeletion     = "deletion"
	MistakeTypeSubstitution = "substitution"
)

// IsPerfect checks if the grade is perfect (no mistakes)
func (g *Grade) IsPerfect() bool {
	return g.WER == 0.0 && g.Score == 100
}

// GetAccuracy returns accuracy percentage (0-100)
func (g *Grade) GetAccuracy() float64 {
	if g.WordCount == 0 {
		return 0.0
	}
	return float64(g.CorrectWords) / float64(g.WordCount) * 100
}