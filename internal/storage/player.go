package storage

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type AudioPlayer struct {
	command string
	args    []string
}

func NewAudioPlayer() (*AudioPlayer, error) {
	player := &AudioPlayer{}

	switch runtime.GOOS {
	case "darwin":
		player.command = "afplay"
		player.args = []string{}
	case "linux":
		if _, err := exec.LookPath("mpg123"); err == nil {
			player.command = "mpg123"
			player.args = []string{"-q"}
		} else if _, err := exec.LookPath("play"); err == nil {
			player.command = "play"
			player.args = []string{"-q"}
		} else if _, err := exec.LookPath("ffplay"); err == nil {
			player.command = "ffplay"
			player.args = []string{"-nodisp", "-autoexit", "-loglevel", "quiet"}
		} else {
			return nil, fmt.Errorf("no audio player found. Please install mpg123, sox, or ffmpeg")
		}
	case "windows":
		player.command = "powershell"
		player.args = []string{"-c", "(New-Object Media.SoundPlayer '%s').PlaySync()"}
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return player, nil
}

func (p *AudioPlayer) Play(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("audio file not found: %w", err)
	}

	var args []string
	if runtime.GOOS == "windows" {
		args = []string{p.args[0], fmt.Sprintf(p.args[1], filePath)}
	} else {
		args = append(p.args, filePath)
	}

	// G204: This is intentional - we need to launch system audio players
	cmd := exec.Command(p.command, args...) // #nosec G204
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start audio player: %w", err)
	}

	return cmd.Wait()
}

func (p *AudioPlayer) PlayAsync(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("audio file not found: %w", err)
	}

	var args []string
	if runtime.GOOS == "windows" {
		args = []string{p.args[0], fmt.Sprintf(p.args[1], filePath)}
	} else {
		args = append(p.args, filePath)
	}

	// G204: This is intentional - we need to launch system audio players
	cmd := exec.Command(p.command, args...) // #nosec G204
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start audio player: %w", err)
	}

	go func() {
		_ = cmd.Wait()
	}()

	return nil
}

func IsAudioPlayerAvailable() bool {
	_, err := NewAudioPlayer()
	return err == nil
}