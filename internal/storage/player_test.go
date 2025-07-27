package storage

import (
	"os"
	"runtime"
	"testing"
)

func TestAudioPlayer(t *testing.T) {
	t.Run("NewAudioPlayer detects platform", func(t *testing.T) {
		player, err := NewAudioPlayer()
		
		switch runtime.GOOS {
		case "darwin":
			if err != nil {
				t.Fatalf("Failed to create audio player on macOS: %v", err)
			}
			if player.command != "afplay" {
				t.Errorf("Expected afplay on macOS, got %s", player.command)
			}
		case "linux":
			if err != nil {
				t.Logf("Audio player not available on Linux: %v", err)
				t.Skip("No audio player found on Linux")
			}
		case "windows":
			if err != nil {
				t.Fatalf("Failed to create audio player on Windows: %v", err)
			}
			if player.command != "powershell" {
				t.Errorf("Expected powershell on Windows, got %s", player.command)
			}
		default:
			if err == nil {
				t.Errorf("Expected error on unsupported platform %s", runtime.GOOS)
			}
		}
	})

	t.Run("Play with non-existent file", func(t *testing.T) {
		player, err := NewAudioPlayer()
		if err != nil {
			t.Skip("No audio player available")
		}

		err = player.Play("/non/existent/file.mp3")
		if err == nil {
			t.Errorf("Expected error when playing non-existent file")
		}
	})

	t.Run("PlayAsync with non-existent file", func(t *testing.T) {
		player, err := NewAudioPlayer()
		if err != nil {
			t.Skip("No audio player available")
		}

		err = player.PlayAsync("/non/existent/file.mp3")
		if err == nil {
			t.Errorf("Expected error when playing non-existent file")
		}
	})

	t.Run("IsAudioPlayerAvailable", func(t *testing.T) {
		available := IsAudioPlayerAvailable()
		
		switch runtime.GOOS {
		case "darwin", "windows":
			if !available {
				t.Errorf("Expected audio player to be available on %s", runtime.GOOS)
			}
		case "linux":
			t.Logf("Audio player available on Linux: %v", available)
		}
	})

	t.Run("Play with valid dummy file", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping actual audio playback in short mode")
		}

		player, err := NewAudioPlayer()
		if err != nil {
			t.Skip("No audio player available")
		}

		tempFile, err := os.CreateTemp("", "test*.mp3")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer func() {
			_ = os.Remove(tempFile.Name())
		}()
		_ = tempFile.Close()

		err = player.Play(tempFile.Name())
		if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
			if err == nil {
				t.Logf("Note: Audio player started but may fail on invalid MP3 content")
			}
		}
	})
}