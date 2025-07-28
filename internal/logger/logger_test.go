package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Create temporary directory for test logs
	tmpDir := filepath.Join(os.TempDir(), "dictcli_test_logs")
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Initialize logger
	err := Init(true, tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() { _ = Close() }()

	// Test different log levels
	Debug("This is a debug message: %s", "test")
	Info("This is an info message: %d", 42)
	Warn("This is a warning message: %v", true)
	Error("This is an error message: %s", "error test")

	// Test TUI logger
	tuiLogger := InitTUILogger()
	
	// Test state transitions
	tuiLogger.StateTransition("StateWelcome", "StateGenerating")
	tuiLogger.StateTransition("StateGenerating", "StatePlaying")
	
	// Test key press
	tuiLogger.KeyPress("R", "StatePlaying")
	tuiLogger.KeyPress("Enter", "StateListening")
	
	// Test user input
	tuiLogger.UserInput("This is my answer", "StateListening")
	tuiLogger.UserInput("This is a very long answer that should be truncated in the log file to avoid excessively long log entries", "StateListening")
	
	// Test API calls
	tuiLogger.APICall("OpenAI", "GenerateSentence", 1234)
	tuiLogger.APIError("OpenAI", "GenerateAudio", os.ErrPermission)
	
	// Test audio playback
	tuiLogger.AudioPlayback("start", "/tmp/audio_123.mp3")
	tuiLogger.AudioPlayback("stop", "/tmp/audio_123.mp3")
	
	// Test settings
	tuiLogger.Settings("voice", "alloy", "echo")
	tuiLogger.Settings("level", "600", "700")
	
	// Test grade
	tuiLogger.Grade(0.15, 85, 3)
	
	// Test session
	tuiLogger.Session("start", map[string]interface{}{
		"id":    "session_123",
		"topic": "Business",
		"level": 700,
	})
	tuiLogger.Session("end", map[string]interface{}{
		"duration": "5m30s",
		"rounds":   5,
	})

	// Give time for async writes
	time.Sleep(100 * time.Millisecond)

	// Check if log files were created
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read log directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("No log files created")
	}

	// Read and verify log content
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".log" {
			filePath := filepath.Join(tmpDir, file.Name())
			content, err := os.ReadFile(filepath.Clean(filePath))
			if err != nil {
				t.Errorf("Failed to read log file: %v", err)
				continue
			}

			// Check for expected content
			contentStr := string(content)
			expectedStrings := []string{
				"[DEBUG]",
				"[INFO]",
				"[WARN]",
				"[ERROR]",
				"STATE_TRANSITION",
				"KEY_PRESS",
				"USER_INPUT",
				"API_CALL",
				"API_ERROR",
				"AUDIO",
				"SETTINGS_CHANGE",
				"GRADE_RESULT",
				"SESSION_START",
				"SESSION_END",
			}

			for _, expected := range expectedStrings {
				if !contains(contentStr, expected) {
					t.Errorf("Log file missing expected content: %s", expected)
				}
			}

			t.Logf("Log file content:\n%s", contentStr)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}