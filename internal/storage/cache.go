// Package storage provides data persistence and audio handling functionality.
package storage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// AudioCache manages local storage of TTS-generated audio files using SHA256-based keys.
// Files are stored in XDG-compliant cache directories for cross-platform compatibility.
type AudioCache struct {
	baseDir string // Base directory for cached audio files
}

// NewAudioCache creates a new audio cache instance with XDG-compliant directory structure.
func NewAudioCache() (*AudioCache, error) {
	cacheDir := filepath.Join(xdg.CacheHome, "dictcli", "audio")
	if err := os.MkdirAll(cacheDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &AudioCache{
		baseDir: cacheDir,
	}, nil
}

// generateKey creates a SHA256-based cache key from text, voice, and speed parameters.
func (c *AudioCache) generateKey(text string, voice string, speed float64) string {
	data := fmt.Sprintf("%s:%s:%.2f", text, voice, speed)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// GetPath returns the file system path for cached audio based on text, voice, and speed.
func (c *AudioCache) GetPath(text string, voice string, speed float64) string {
	key := c.generateKey(text, voice, speed)
	return filepath.Join(c.baseDir, key+".mp3")
}

// Exists checks if cached audio already exists for the given parameters.
func (c *AudioCache) Exists(text string, voice string, speed float64) bool {
	path := c.GetPath(text, voice, speed)
	_, err := os.Stat(path)
	return err == nil
}

// Save stores audio data to the cache with the given parameters as the key.
func (c *AudioCache) Save(text string, voice string, speed float64, data []byte) error {
	path := c.GetPath(text, voice, speed)
	return os.WriteFile(path, data, 0600)
}

// Load retrieves cached audio data for the given parameters.
func (c *AudioCache) Load(text string, voice string, speed float64) ([]byte, error) {
	path := c.GetPath(text, voice, speed)
	return os.ReadFile(path) // #nosec G304
}

// SaveFromReader saves audio data from an io.Reader to the cache, useful for streaming audio data.
func (c *AudioCache) SaveFromReader(text string, voice string, speed float64, reader io.Reader) error {
	path := c.GetPath(text, voice, speed)
	
	file, err := os.Create(path) // #nosec G304
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err := io.Copy(file, reader); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Clear removes all cached audio files and the cache directory.
func (c *AudioCache) Clear() error {
	return os.RemoveAll(c.baseDir)
}

// Size returns the total size in bytes and number of cached audio files.
func (c *AudioCache) Size() (int64, int, error) {
	var totalSize int64
	var fileCount int

	err := filepath.Walk(c.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".mp3" {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	return totalSize, fileCount, err
}