package types

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   Config
	}{
		{
			name: "valid config",
			config: Config{
				Level: 700,
				Words: 15,
				Speed: 1.0,
			},
			want: Config{
				Level: 700,
				Words: 15,
				Speed: 1.0,
			},
		},
		{
			name: "level too low",
			config: Config{
				Level: 300,
				Words: 15,
				Speed: 1.0,
			},
			want: Config{
				Level: DefaultLevel,
				Words: 15,
				Speed: 1.0,
			},
		},
		{
			name: "level too high",
			config: Config{
				Level: 1000,
				Words: 15,
				Speed: 1.0,
			},
			want: Config{
				Level: DefaultLevel,
				Words: 15,
				Speed: 1.0,
			},
		},
		{
			name: "words too low",
			config: Config{
				Level: 700,
				Words: 3,
				Speed: 1.0,
			},
			want: Config{
				Level: 700,
				Words: DefaultWords,
				Speed: 1.0,
			},
		},
		{
			name: "words too high",
			config: Config{
				Level: 700,
				Words: 50,
				Speed: 1.0,
			},
			want: Config{
				Level: 700,
				Words: DefaultWords,
				Speed: 1.0,
			},
		},
		{
			name: "speed too low",
			config: Config{
				Level: 700,
				Words: 15,
				Speed: 0.3,
			},
			want: Config{
				Level: 700,
				Words: 15,
				Speed: DefaultSpeed,
			},
		},
		{
			name: "speed too high",
			config: Config{
				Level: 700,
				Words: 15,
				Speed: 3.0,
			},
			want: Config{
				Level: 700,
				Words: 15,
				Speed: DefaultSpeed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.config
			ValidateConfig(&cfg)
			
			if cfg.Level != tt.want.Level {
				t.Errorf("Level = %v, want %v", cfg.Level, tt.want.Level)
			}
			if cfg.Words != tt.want.Words {
				t.Errorf("Words = %v, want %v", cfg.Words, tt.want.Words)
			}
			if cfg.Speed != tt.want.Speed {
				t.Errorf("Speed = %v, want %v", cfg.Speed, tt.want.Speed)
			}
		})
	}
}

func TestIsValidVoice(t *testing.T) {
	tests := []struct {
		voice string
		want  bool
	}{
		{VoiceAlloy, true},
		{VoiceEcho, true},
		{VoiceFable, true},
		{VoiceOnyx, true},
		{VoiceNova, true},
		{VoiceShimmer, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.voice, func(t *testing.T) {
			if got := IsValidVoice(tt.voice); got != tt.want {
				t.Errorf("IsValidVoice(%v) = %v, want %v", tt.voice, got, tt.want)
			}
		})
	}
}

func TestIsValidTopic(t *testing.T) {
	tests := []struct {
		topic string
		want  bool
	}{
		{TopicBusiness, true},
		{TopicTravel, true},
		{TopicDaily, true},
		{TopicTechnology, true},
		{TopicHealth, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.topic, func(t *testing.T) {
			if got := IsValidTopic(tt.topic); got != tt.want {
				t.Errorf("IsValidTopic(%v) = %v, want %v", tt.topic, got, tt.want)
			}
		})
	}
}