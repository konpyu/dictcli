package types

import (
	"time"
)

type Config struct {
	Voice   string  `json:"voice" mapstructure:"voice"`
	Level   int     `json:"level" mapstructure:"level"`
	Topic   string  `json:"topic" mapstructure:"topic"`
	Words   int     `json:"words" mapstructure:"words"`
	Speed   float64 `json:"speed" mapstructure:"speed"`
	NoCache bool    `json:"no_cache" mapstructure:"no_cache"`
	Debug   bool    `json:"debug" mapstructure:"debug"`
}

type DictationSession struct {
	ID             string        `json:"id"`
	Timestamp      time.Time     `json:"timestamp"`
	Config         Config        `json:"config"`
	Sentence       string        `json:"sentence"`
	AudioPath      string        `json:"audio_path"`
	UserInput      string        `json:"user_input"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	ReplayCount    int           `json:"replay_count"`
	Grade          *Grade        `json:"grade"`
	DurationSecs   float64       `json:"duration_secs"`
}

type Mistake struct {
	Position int    `json:"position"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Type     string `json:"type"`
}

type Grade struct {
	WER                  float64   `json:"wer"`
	Score                int       `json:"score"`
	Mistakes             []Mistake `json:"mistakes"`
	JapaneseExplanation  string    `json:"japanese_explanation"`
	AlternativeExpressions []string `json:"alternative_expressions"`
}

const (
	DefaultVoice  = "TOM"
	DefaultLevel  = 700
	DefaultTopic  = "Business"
	DefaultWords  = 15
	DefaultSpeed  = 1.0
)

const (
	VoiceAlloy   = "alloy"
	VoiceEcho    = "echo"
	VoiceFable   = "fable"
	VoiceOnyx    = "onyx"
	VoiceNova    = "nova"
	VoiceShimmer = "shimmer"
)

const (
	TopicBusiness   = "Business"
	TopicTravel     = "Travel"
	TopicDaily      = "Daily"
	TopicTechnology = "Technology"
	TopicHealth     = "Health"
)

const (
	MinLevel = 400
	MaxLevel = 990
	MinWords = 5
	MaxWords = 30
	MinSpeed = 0.5
	MaxSpeed = 2.0
)

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

func IsValidVoice(voice string) bool {
	switch voice {
	case VoiceAlloy, VoiceEcho, VoiceFable, VoiceOnyx, VoiceNova, VoiceShimmer:
		return true
	default:
		return false
	}
}

func IsValidTopic(topic string) bool {
	switch topic {
	case TopicBusiness, TopicTravel, TopicDaily, TopicTechnology, TopicHealth:
		return true
	default:
		return false
	}
}