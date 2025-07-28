package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level represents the severity of a log message
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of a log level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger is the main logging interface
type Logger struct {
	mu          sync.Mutex
	level       Level
	debugMode   bool
	fileWriter  io.Writer
	stdLogger   *log.Logger
	logFile     *os.File
	logDir      string
	maxFileSize int64
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init initializes the global logger instance
func Init(debugMode bool, logDir string) error {
	var err error
	once.Do(func() {
		defaultLogger, err = New(debugMode, logDir)
	})
	return err
}

// New creates a new logger instance
func New(debugMode bool, logDir string) (*Logger, error) {
	l := &Logger{
		level:       INFO,
		debugMode:   debugMode,
		logDir:      logDir,
		maxFileSize: 10 * 1024 * 1024, // 10MB
	}

	if debugMode {
		l.level = DEBUG
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	if err := l.openLogFile(); err != nil {
		return nil, err
	}

	return l, nil
}

// openLogFile creates or opens a log file with timestamp
func (l *Logger) openLogFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Close existing file if any
	if l.logFile != nil {
		_ = l.logFile.Close()
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("dictcli_%s.log", timestamp)
	filepath := filepath.Join(l.logDir, filename)

	// Open file for writing
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.logFile = file
	l.fileWriter = file
	l.stdLogger = log.New(l.fileWriter, "", 0)

	return nil
}

// checkRotation checks if log file needs rotation
func (l *Logger) checkRotation() error {
	if l.logFile == nil {
		return nil
	}

	info, err := l.logFile.Stat()
	if err != nil {
		return err
	}

	if info.Size() >= l.maxFileSize {
		return l.openLogFile()
	}

	return nil
}

// log writes a log message
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if rotation is needed
	if err := l.checkRotation(); err != nil {
		// Fall back to stderr if rotation fails
		fmt.Fprintf(os.Stderr, "Log rotation failed: %v\n", err)
	}

	// Format message
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, level.String(), message)

	// Write to file if available
	if l.stdLogger != nil {
		l.stdLogger.Println(logLine)
	}

	// Also write to stderr in debug mode for ERROR level
	if l.debugMode && level == ERROR {
		fmt.Fprintln(os.Stderr, logLine)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Close closes the logger and its resources
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Global logger functions

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, args...)
	}
}

// Info logs an info message using the default logger
func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, args...)
	}
}

// Warn logs a warning message using the default logger
func Warn(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, args...)
	}
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, args...)
	}
}

// Close closes the default logger
func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}