// Package service provides the core business logic for dictation practice.
package service

import (
	"context"
	"fmt"
	"log"

	"github.com/konpyu/dictcli/internal/storage"
	"github.com/konpyu/dictcli/internal/types"
)

// DictationService orchestrates the complete dictation practice flow,
// coordinating sentence generation, audio synthesis, and grading.
type DictationService struct {
	openai *OpenAIService        // OpenAI API integration
	cache  *storage.AudioCache   // Audio file caching
	debug  bool                  // Debug logging enabled
}

// NewDictationService creates a new dictation service with OpenAI integration and audio caching.
func NewDictationService(debug bool) (*DictationService, error) {
	openaiService, err := NewOpenAIService(debug)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI service: %w", err)
	}

	cache, err := storage.NewAudioCache()
	if err != nil {
		return nil, fmt.Errorf("failed to create audio cache: %w", err)
	}

	return &DictationService{
		openai: openaiService,
		cache:  cache,
		debug:  debug,
	}, nil
}

// GenerateSentence creates a new sentence for dictation practice based on the specified parameters.
// The sentence complexity is adjusted based on the TOEIC level, and content is themed around the specified topic.
func (s *DictationService) GenerateSentence(ctx context.Context, topic string, level int, wordCount int) (string, error) {
	config := &types.Config{
		Topic: topic,
		Level: level,
		Words: wordCount,
	}

	return s.openai.GenerateSentence(ctx, config)
}

// GenerateAudio creates audio for the given text using TTS, with intelligent caching to minimize API calls.
// Returns the file path to the generated audio. If audio is already cached, returns immediately.
func (s *DictationService) GenerateAudio(ctx context.Context, text string, voice string, speed float64) (string, error) {
	if s.cache.Exists(text, voice, speed) {
		if s.debug {
			log.Printf("[Cache] Audio cache HIT for voice=%s, speed=%.1f, text_len=%d", voice, speed, len(text))
		}
		return s.cache.GetPath(text, voice, speed), nil
	}

	if s.debug {
		log.Printf("[Cache] Audio cache MISS for voice=%s, speed=%.1f, text_len=%d", voice, speed, len(text))
	}

	audioData, err := s.openai.GenerateAudio(ctx, text, voice, speed)
	if err != nil {
		return "", fmt.Errorf("failed to generate audio: %w", err)
	}

	if err := s.cache.Save(text, voice, speed, audioData); err != nil {
		if s.debug {
			log.Printf("[Cache] Failed to save audio to cache: %v", err)
		}
		return "", fmt.Errorf("failed to save audio to cache: %w", err)
	}

	if s.debug {
		log.Printf("[Cache] Audio saved to cache, size=%d bytes", len(audioData))
	}

	return s.cache.GetPath(text, voice, speed), nil
}

// GradeDictation evaluates the user's input against the reference sentence,
// providing detailed feedback including score, error analysis, and Japanese explanations.
func (s *DictationService) GradeDictation(ctx context.Context, reference string, userInput string) (*types.Grade, error) {
	return s.openai.GradeDictation(ctx, reference, userInput)
}