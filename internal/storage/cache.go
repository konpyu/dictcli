package storage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type AudioCache struct {
	baseDir string
}

func NewAudioCache() (*AudioCache, error) {
	cacheDir := filepath.Join(xdg.CacheHome, "dictcli", "audio")
	if err := os.MkdirAll(cacheDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &AudioCache{
		baseDir: cacheDir,
	}, nil
}

func (c *AudioCache) generateKey(text string, voice string, speed float64) string {
	data := fmt.Sprintf("%s:%s:%.2f", text, voice, speed)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func (c *AudioCache) GetPath(text string, voice string, speed float64) string {
	key := c.generateKey(text, voice, speed)
	return filepath.Join(c.baseDir, key+".mp3")
}

func (c *AudioCache) Exists(text string, voice string, speed float64) bool {
	path := c.GetPath(text, voice, speed)
	_, err := os.Stat(path)
	return err == nil
}

func (c *AudioCache) Save(text string, voice string, speed float64, data []byte) error {
	path := c.GetPath(text, voice, speed)
	return os.WriteFile(path, data, 0600)
}

func (c *AudioCache) Load(text string, voice string, speed float64) ([]byte, error) {
	path := c.GetPath(text, voice, speed)
	return os.ReadFile(path) // #nosec G304
}

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

func (c *AudioCache) Clear() error {
	return os.RemoveAll(c.baseDir)
}

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