// Package types defines the core data structures for DictCLI
package types

// Config represents the application configuration
type Config struct {
	// API Configuration
	OpenAIAPIKey string `json:"openai_api_key" mapstructure:"openai_api_key"`
	
	// Dictation Settings
	Voice      string  `json:"voice" mapstructure:"voice"`           // alloy/echo/fable/onyx/nova/shimmer
	Level      int     `json:"level" mapstructure:"level"`           // TOEIC level (400-990)
	Topic      string  `json:"topic" mapstructure:"topic"`           // Business/Travel/Daily/Technology/Health
	WordCount  int     `json:"word_count" mapstructure:"word_count"` // 5-30 words
	SpeechSpeed float64 `json:"speech_speed" mapstructure:"speech_speed"` // 0.5-2.0
	
	// Audio Settings
	NoCache bool `json:"no_cache" mapstructure:"no_cache"`
	
	// Debug Settings
	Debug bool `json:"debug" mapstructure:"debug"`
	
	// Paths (computed, not from config file)
	ConfigPath string `json:"-"`
	CachePath  string `json:"-"`
	DataPath   string `json:"-"`
	LogPath    string `json:"-"`
}

// Voice options
const (
	VoiceAlloy   = "alloy"
	VoiceEcho    = "echo"
	VoiceFable   = "fable"
	VoiceOnyx    = "onyx"
	VoiceNova    = "nova"
	VoiceShimmer = "shimmer"
)

// Topic options
const (
	TopicBusiness   = "Business"
	TopicTravel     = "Travel"
	TopicDaily      = "Daily"
	TopicTechnology = "Technology"
	TopicHealth     = "Health"
)

// Default configuration values
const (
	DefaultVoice       = VoiceAlloy
	DefaultLevel       = 700
	DefaultTopic       = TopicBusiness
	DefaultWordCount   = 15
	DefaultSpeechSpeed = 1.0
	
	MinLevel       = 400
	MaxLevel       = 990
	MinWordCount   = 5
	MaxWordCount   = 30
	MinSpeechSpeed = 0.5
	MaxSpeechSpeed = 2.0
)

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate voice
	validVoices := map[string]bool{
		VoiceAlloy: true, VoiceEcho: true, VoiceFable: true,
		VoiceOnyx: true, VoiceNova: true, VoiceShimmer: true,
	}
	if !validVoices[c.Voice] {
		c.Voice = DefaultVoice
	}
	
	// Validate level
	if c.Level < MinLevel || c.Level > MaxLevel {
		c.Level = DefaultLevel
	}
	
	// Validate topic
	validTopics := map[string]bool{
		TopicBusiness: true, TopicTravel: true, TopicDaily: true,
		TopicTechnology: true, TopicHealth: true,
	}
	if !validTopics[c.Topic] {
		c.Topic = DefaultTopic
	}
	
	// Validate word count
	if c.WordCount < MinWordCount || c.WordCount > MaxWordCount {
		c.WordCount = DefaultWordCount
	}
	
	// Validate speech speed
	if c.SpeechSpeed < MinSpeechSpeed || c.SpeechSpeed > MaxSpeechSpeed {
		c.SpeechSpeed = DefaultSpeechSpeed
	}
	
	return nil
}

// GetDefaultConfig returns a config with default values
func GetDefaultConfig() *Config {
	return &Config{
		Voice:       DefaultVoice,
		Level:       DefaultLevel,
		Topic:       DefaultTopic,
		WordCount:   DefaultWordCount,
		SpeechSpeed: DefaultSpeechSpeed,
		NoCache:     false,
		Debug:       false,
	}
}