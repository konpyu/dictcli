// Package types defines core data structures and validation for DictCLI.
package types

import (
	"time"
)

// Config represents user configuration settings for dictation practice.
type Config struct {
	Voice   string  `json:"voice" mapstructure:"voice"`       // TTS voice (alloy, echo, fable, onyx, nova, shimmer)
	Level   int     `json:"level" mapstructure:"level"`       // TOEIC level (400-990)
	Topic   string  `json:"topic" mapstructure:"topic"`       // Topic category (Business, Travel, etc.)
	Words   int     `json:"words" mapstructure:"words"`       // Word count per sentence (5-30)
	Speed   float64 `json:"speed" mapstructure:"speed"`       // Speech speed multiplier (0.5-2.0)
	NoCache bool    `json:"no_cache" mapstructure:"no_cache"` // Disable audio caching
	Debug   bool    `json:"debug" mapstructure:"debug"`       // Enable debug logging
}

// DictationSession represents a complete dictation practice session.
type DictationSession struct {
	ID           string    `json:"id"`           // Unique session identifier
	Timestamp    time.Time `json:"timestamp"`    // Session creation time
	Config       Config    `json:"config"`       // Configuration used for this session
	Sentence     string    `json:"sentence"`     // Generated sentence for dictation
	AudioPath    string    `json:"audio_path"`   // Path to cached audio file
	UserInput    string    `json:"user_input"`   // User's dictation input
	StartTime    time.Time `json:"start_time"`   // When user started typing
	EndTime      time.Time `json:"end_time"`     // When user submitted answer
	ReplayCount  int       `json:"replay_count"` // Number of times audio was replayed
	Grade        *Grade    `json:"grade"`        // Grading results (nil if not graded)
	DurationSecs float64   `json:"duration_secs"` // Total session duration in seconds
}

// Mistake represents a single dictation error with position and type information.
type Mistake struct {
	Position int    `json:"position"` // Word position in sentence (0-based)
	Expected string `json:"expected"` // Correct word/phrase
	Actual   string `json:"actual"`   // User's input
	Type     string `json:"type"`     // Error type (substitution, insertion, deletion)
}

// Grade represents the assessment results of a dictation session.
type Grade struct {
	WER                    float64   `json:"wer"`                      // Word Error Rate (0.0-1.0)
	Score                  int       `json:"score"`                    // Overall score (0-100)
	Mistakes               []Mistake `json:"mistakes"`                 // Detailed list of errors
	JapaneseExplanation    string    `json:"japanese_explanation"`     // Explanation in Japanese
	AlternativeExpressions []string  `json:"alternative_expressions"`  // Alternative ways to express the same idea
}

// Default configuration values.
const (
	DefaultVoice = "alloy"     // Default TTS voice
	DefaultLevel = 700         // Default TOEIC level
	DefaultTopic = "Business"  // Default topic category  
	DefaultWords = 15          // Default word count per sentence
	DefaultSpeed = 1.0         // Default speech speed multiplier
)

// Available TTS voice options.
const (
	VoiceAlloy   = "alloy"   // Balanced, natural voice
	VoiceEcho    = "echo"    // Clear, articulate voice
	VoiceFable   = "fable"   // Warm, engaging voice
	VoiceOnyx    = "onyx"    // Deep, authoritative voice
	VoiceNova    = "nova"    // Energetic, youthful voice
	VoiceShimmer = "shimmer" // Gentle, calming voice
)

// Available topic categories for sentence generation.
const (
	TopicBusiness   = "Business"   // Business and workplace scenarios
	TopicTravel     = "Travel"     // Travel and tourism contexts
	TopicDaily      = "Daily"      // Everyday life situations
	TopicTechnology = "Technology" // Technology and digital topics
	TopicHealth     = "Health"     // Health and medical contexts
)

// Configuration value constraints.
const (
	MinLevel = 400 // Minimum TOEIC level
	MaxLevel = 990 // Maximum TOEIC level
	MinWords = 5   // Minimum words per sentence
	MaxWords = 30  // Maximum words per sentence
	MinSpeed = 0.5 // Minimum speech speed multiplier
	MaxSpeed = 2.0 // Maximum speech speed multiplier
)

// ValidVoices contains all supported TTS voice options.
var ValidVoices = []string{
	VoiceAlloy,
	VoiceEcho,
	VoiceFable,
	VoiceOnyx,
	VoiceNova,
	VoiceShimmer,
}

// ValidTopics contains all supported topic categories.
var ValidTopics = []string{
	TopicBusiness,
	TopicTravel,
	TopicDaily,
	TopicTechnology,
	TopicHealth,
}

// ValidateConfig validates and corrects configuration values to ensure they are within acceptable ranges.
// Invalid values are automatically corrected to defaults rather than returning an error.
func ValidateConfig(c *Config) error {
	if c.Level < MinLevel || c.Level > MaxLevel {
		c.Level = DefaultLevel
	}
	if c.Words < MinWords || c.Words > MaxWords {
		c.Words = DefaultWords
	}
	if c.Speed < MinSpeed || c.Speed > MaxSpeed {
		c.Speed = DefaultSpeed
	}
	return nil
}

// IsValidVoice checks if the given voice is supported by the TTS system.
func IsValidVoice(voice string) bool {
	switch voice {
	case VoiceAlloy, VoiceEcho, VoiceFable, VoiceOnyx, VoiceNova, VoiceShimmer:
		return true
	default:
		return false
	}
}

// IsValidTopic checks if the given topic is supported for sentence generation.
func IsValidTopic(topic string) bool {
	switch topic {
	case TopicBusiness, TopicTravel, TopicDaily, TopicTechnology, TopicHealth:
		return true
	default:
		return false
	}
}