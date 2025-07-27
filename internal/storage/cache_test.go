package storage

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestAudioCache(t *testing.T) {
	t.Run("NewAudioCache creates cache directory", func(t *testing.T) {
		tempDir := t.TempDir()
		cache := &AudioCache{
			baseDir: filepath.Join(tempDir, "dictcli", "audio"),
		}
		
		err := os.MkdirAll(cache.baseDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test cache directory: %v", err)
		}
		
		if _, err := os.Stat(cache.baseDir); os.IsNotExist(err) {
			t.Errorf("Cache directory was not created")
		}
	})

	t.Run("generateKey creates consistent hash", func(t *testing.T) {
		cache := &AudioCache{}
		
		key1 := cache.generateKey("Hello world", "alloy", 1.0)
		key2 := cache.generateKey("Hello world", "alloy", 1.0)
		
		if key1 != key2 {
			t.Errorf("Same inputs produced different keys: %s vs %s", key1, key2)
		}
		
		key3 := cache.generateKey("Hello world", "echo", 1.0)
		if key1 == key3 {
			t.Errorf("Different voices produced same key")
		}
		
		key4 := cache.generateKey("Hello world", "alloy", 1.5)
		if key1 == key4 {
			t.Errorf("Different speeds produced same key")
		}
	})

	t.Run("Save and Load audio data", func(t *testing.T) {
		tempDir := t.TempDir()
		cache := &AudioCache{baseDir: tempDir}
		
		testData := []byte("test audio data")
		text := "Test sentence"
		voice := "alloy"
		speed := 1.0
		
		err := cache.Save(text, voice, speed, testData)
		if err != nil {
			t.Fatalf("Failed to save audio: %v", err)
		}
		
		if !cache.Exists(text, voice, speed) {
			t.Errorf("Cache reported file doesn't exist after save")
		}
		
		loadedData, err := cache.Load(text, voice, speed)
		if err != nil {
			t.Fatalf("Failed to load audio: %v", err)
		}
		
		if !bytes.Equal(loadedData, testData) {
			t.Errorf("Loaded data doesn't match saved data")
		}
	})

	t.Run("SaveFromReader", func(t *testing.T) {
		tempDir := t.TempDir()
		cache := &AudioCache{baseDir: tempDir}
		
		testData := []byte("test audio stream")
		reader := bytes.NewReader(testData)
		
		err := cache.SaveFromReader("Test", "alloy", 1.0, reader)
		if err != nil {
			t.Fatalf("Failed to save from reader: %v", err)
		}
		
		loadedData, err := cache.Load("Test", "alloy", 1.0)
		if err != nil {
			t.Fatalf("Failed to load audio: %v", err)
		}
		
		if !bytes.Equal(loadedData, testData) {
			t.Errorf("Loaded data doesn't match saved data")
		}
	})

	t.Run("Size calculation", func(t *testing.T) {
		tempDir := t.TempDir()
		cache := &AudioCache{baseDir: tempDir}
		
		_ = cache.Save("Test1", "alloy", 1.0, []byte("data1"))
		_ = cache.Save("Test2", "echo", 1.0, []byte("data22"))
		_ = cache.Save("Test3", "nova", 1.0, []byte("data333"))
		
		size, count, err := cache.Size()
		if err != nil {
			t.Fatalf("Failed to calculate size: %v", err)
		}
		
		if count != 3 {
			t.Errorf("Expected 3 files, got %d", count)
		}
		
		expectedSize := int64(5 + 6 + 7)
		if size != expectedSize {
			t.Errorf("Expected size %d, got %d", expectedSize, size)
		}
	})

	t.Run("Clear cache", func(t *testing.T) {
		tempDir := t.TempDir()
		cache := &AudioCache{baseDir: tempDir}
		
		_ = cache.Save("Test", "alloy", 1.0, []byte("data"))
		
		err := cache.Clear()
		if err != nil {
			t.Fatalf("Failed to clear cache: %v", err)
		}
		
		if _, err := os.Stat(cache.baseDir); !os.IsNotExist(err) {
			t.Errorf("Cache directory still exists after clear")
		}
	})
}