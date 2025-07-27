package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/konpyu/dictcli/internal/types"
	"github.com/spf13/viper"
)

type Manager struct {
	viper *viper.Viper
	cfg   *types.Config
}

func New() (*Manager, error) {
	v := viper.New()
	
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	
	configDir := filepath.Join(xdg.ConfigHome, "dictcli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	v.AddConfigPath(configDir)
	v.AddConfigPath(".")
	
	v.SetEnvPrefix("DICTCLI")
	v.AutomaticEnv()
	
	setDefaults(v)
	
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	var cfg types.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	if err := types.ValidateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	
	return &Manager{
		viper: v,
		cfg:   &cfg,
	}, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("voice", types.DefaultVoice)
	v.SetDefault("level", types.DefaultLevel)
	v.SetDefault("topic", types.DefaultTopic)
	v.SetDefault("words", types.DefaultWords)
	v.SetDefault("speed", types.DefaultSpeed)
	v.SetDefault("no_cache", false)
	v.SetDefault("debug", false)
}

func (m *Manager) Get() *types.Config {
	return m.cfg
}

func (m *Manager) Set(key string, value interface{}) error {
	m.viper.Set(key, value)
	
	if err := m.viper.Unmarshal(m.cfg); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}
	
	return types.ValidateConfig(m.cfg)
}

func (m *Manager) Save() error {
	configDir := filepath.Join(xdg.ConfigHome, "dictcli")
	configPath := filepath.Join(configDir, "config.yaml")
	
	if err := m.viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	return nil
}

func (m *Manager) GetConfigPath() string {
	return m.viper.ConfigFileUsed()
}

func (m *Manager) SetFromFlags(voice, topic string, level, words int, speed float64, noCache, debug bool) {
	if voice != "" && types.IsValidVoice(voice) {
		m.cfg.Voice = voice
	}
	if topic != "" && types.IsValidTopic(topic) {
		m.cfg.Topic = topic
	}
	if level >= types.MinLevel && level <= types.MaxLevel {
		m.cfg.Level = level
	}
	if words >= types.MinWords && words <= types.MaxWords {
		m.cfg.Words = words
	}
	if speed >= types.MinSpeed && speed <= types.MaxSpeed {
		m.cfg.Speed = speed
	}
	
	m.cfg.NoCache = noCache
	m.cfg.Debug = debug
}

func GetDefaultConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "dictcli", "config.yaml")
}

func Load() (*types.Config, error) {
	manager, err := New()
	if err != nil {
		return nil, err
	}
	return manager.Get(), nil
}