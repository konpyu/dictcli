package logger

import (
	"fmt"
	"strings"
)

// TUILogger provides specialized logging for TUI state transitions
type TUILogger struct {
	logger *Logger
}

// NewTUILogger creates a new TUI logger
func NewTUILogger(l *Logger) *TUILogger {
	if l == nil {
		l = defaultLogger
	}
	return &TUILogger{logger: l}
}

// StateTransition logs a state transition
func (t *TUILogger) StateTransition(from, to string) {
	if t.logger == nil {
		return
	}
	t.logger.Debug("STATE_TRANSITION: %s -> %s", from, to)
}

// KeyPress logs a key press event
func (t *TUILogger) KeyPress(key string, state string) {
	if t.logger == nil {
		return
	}
	t.logger.Debug("KEY_PRESS: key='%s' state='%s'", key, state)
}

// UserInput logs user input
func (t *TUILogger) UserInput(input string, state string) {
	if t.logger == nil {
		return
	}
	// Truncate long input for logging
	if len(input) > 50 {
		input = input[:50] + "..."
	}
	t.logger.Debug("USER_INPUT: input='%s' state='%s'", input, state)
}

// APICall logs an API call
func (t *TUILogger) APICall(service, method string, duration int64) {
	if t.logger == nil {
		return
	}
	t.logger.Info("API_CALL: service='%s' method='%s' duration=%dms", service, method, duration)
}

// APIError logs an API error
func (t *TUILogger) APIError(service, method string, err error) {
	if t.logger == nil {
		return
	}
	t.logger.Error("API_ERROR: service='%s' method='%s' error='%v'", service, method, err)
}

// AudioPlayback logs audio playback events
func (t *TUILogger) AudioPlayback(action string, file string) {
	if t.logger == nil {
		return
	}
	t.logger.Debug("AUDIO: action='%s' file='%s'", action, file)
}

// Settings logs settings changes
func (t *TUILogger) Settings(field, oldValue, newValue string) {
	if t.logger == nil {
		return
	}
	t.logger.Info("SETTINGS_CHANGE: field='%s' old='%s' new='%s'", field, oldValue, newValue)
}

// Grade logs grading results
func (t *TUILogger) Grade(wer float64, score int, mistakes int) {
	if t.logger == nil {
		return
	}
	t.logger.Info("GRADE_RESULT: wer=%.3f score=%d mistakes=%d", wer, score, mistakes)
}

// Session logs session events
func (t *TUILogger) Session(event string, details map[string]interface{}) {
	if t.logger == nil {
		return
	}
	var parts []string
	for k, v := range details {
		parts = append(parts, fmt.Sprintf("%s='%v'", k, v))
	}
	t.logger.Info("SESSION_%s: %s", strings.ToUpper(event), strings.Join(parts, " "))
}

// Global TUI logger instance
var globalTUILogger *TUILogger

// InitTUILogger initializes the global TUI logger
func InitTUILogger() *TUILogger {
	globalTUILogger = NewTUILogger(defaultLogger)
	return globalTUILogger
}

// GetTUILogger returns the global TUI logger
func GetTUILogger() *TUILogger {
	if globalTUILogger == nil {
		globalTUILogger = NewTUILogger(defaultLogger)
	}
	return globalTUILogger
}