package service

import (
	"context"
	"fmt"
	"time"

	"github.com/konpyu/dictcli/internal/storage"
	"github.com/konpyu/dictcli/internal/types"
)

type DictationService struct {
	openai *OpenAIService
	cache  *storage.AudioCache
	debug  bool
}

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

func (s *DictationService) GenerateSentence(topic string, level int, wordCount int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config := &types.Config{
		Topic:     topic,
		Level:     level,
		Words: wordCount,
	}

	return s.openai.GenerateSentence(ctx, config)
}

func (s *DictationService) GenerateAudio(text string, voice string, speed float64) (string, error) {
	if s.cache.Exists(text, voice, speed) {
		if s.debug {
			fmt.Println("Audio cache hit")
		}
		return s.cache.GetPath(text, voice, speed), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	audioData, err := s.openai.GenerateAudio(ctx, text, voice, speed)
	if err != nil {
		return "", fmt.Errorf("failed to generate audio: %w", err)
	}

	if err := s.cache.Save(text, voice, speed, audioData); err != nil {
		return "", fmt.Errorf("failed to save audio to cache: %w", err)
	}

	return s.cache.GetPath(text, voice, speed), nil
}

func (s *DictationService) GradeDictation(reference string, userInput string) (*types.Grade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.openai.GradeDictation(ctx, reference, userInput)
}