package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

)

var (
	logger     *Logger
	initOnce   sync.Once
	globalFile *os.File
)

type Logger struct {
	logger *log.Logger
	mu     sync.Mutex
	debug  bool
}

// Initialize sets up the logger with file output
func Initialize(debug bool) error {
	var initErr error
	initOnce.Do(func() {
		logDir := "logs"
		if err := os.MkdirAll(logDir, 0750); err != nil {
			initErr = fmt.Errorf("failed to create log directory: %w", err)
			return
		}

		logPath := filepath.Join(logDir, "dictcli.log")
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		if err != nil {
			initErr = fmt.Errorf("failed to open log file: %w", err)
			return
		}
		globalFile = file

		// Create logger that writes to file only (not stdout)
		logger = &Logger{
			logger: log.New(file, "", log.LstdFlags|log.Lmicroseconds),
			debug:  debug,
		}
	})
	return initErr
}

// Close closes the log file
func Close() {
	if globalFile != nil {
		globalFile.Close()
	}
}

// SetDebug enables or disables debug logging
func SetDebug(debug bool) {
	if logger != nil {
		logger.mu.Lock()
		logger.debug = debug
		logger.mu.Unlock()
	}
}

// Info logs an informational message (always logged)
func Info(format string, v ...interface{}) {
	if logger != nil {
		logger.mu.Lock()
		defer logger.mu.Unlock()
		logger.logger.Printf("[INFO] "+format, v...)
	}
}

// Debug logs a debug message (only if debug mode is enabled)
func Debug(format string, v ...interface{}) {
	if logger != nil {
		logger.mu.Lock()
		defer logger.mu.Unlock()
		if logger.debug {
			logger.logger.Printf("[DEBUG] "+format, v...)
		}
	}
}

// Error logs an error message (always logged)
func Error(format string, v ...interface{}) {
	if logger != nil {
		logger.mu.Lock()
		defer logger.mu.Unlock()
		logger.logger.Printf("[ERROR] "+format, v...)
	}
}

// Warn logs a warning message (always logged)
func Warn(format string, v ...interface{}) {
	if logger != nil {
		logger.mu.Lock()
		defer logger.mu.Unlock()
		logger.logger.Printf("[WARN] "+format, v...)
	}
}

// Writer returns an io.Writer for the logger
func Writer() io.Writer {
	if logger != nil {
		return logger.logger.Writer()
	}
	return os.Stderr
}