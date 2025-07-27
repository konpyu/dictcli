package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/konpyu/dictcli/internal/types"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
	debug  bool
}

func NewOpenAIService(debug bool) (*OpenAIService, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is not set")
	}

	client := openai.NewClient(apiKey)
	
	return &OpenAIService{
		client: client,
		debug:  debug,
	}, nil
}

func (s *OpenAIService) withRetry(ctx context.Context, operation func() error) error {
	maxRetries := 3
	baseDelay := time.Second
	
	for i := 0; i < maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		
		if i == maxRetries-1 {
			return err
		}
		
		delay := baseDelay * time.Duration(1<<i)
		if s.debug {
			fmt.Printf("Retry %d/%d after %v due to error: %v\n", i+1, maxRetries, delay, err)
		}
		
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	
	return errors.New("max retries exceeded")
}

func (s *OpenAIService) GenerateSentence(ctx context.Context, config *types.Config) (string, error) {
	prompt := buildSentencePrompt(config)
	
	var response openai.ChatCompletionResponse
	err := s.withRetry(ctx, func() error {
		var err error
		response, err = s.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT4oMini,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "You are an English teacher creating sentences for dictation practice. Create natural, contextually appropriate sentences.",
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
				Temperature: 0.7,
				MaxTokens:   100,
			},
		)
		return err
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to generate sentence: %w", err)
	}
	
	if len(response.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}
	
	return response.Choices[0].Message.Content, nil
}

func (s *OpenAIService) GenerateAudio(ctx context.Context, text string, voice string, speed float64) ([]byte, error) {
	if !types.IsValidVoice(voice) {
		voice = types.VoiceOnyx
	}
	
	var audioData []byte
	err := s.withRetry(ctx, func() error {
		response, err := s.client.CreateSpeech(ctx, openai.CreateSpeechRequest{
			Model:          openai.TTSModel1,
			Input:          text,
			Voice:          openai.SpeechVoice(voice),
			ResponseFormat: openai.SpeechResponseFormatMp3,
			Speed:          speed,
		})
		if err != nil {
			return err
		}
		defer func() {
			_ = response.Close()
		}()
		
		buf := make([]byte, 0, 1024*1024) // 1MB initial capacity
		tmp := make([]byte, 1024)
		for {
			n, err := response.Read(tmp)
			if n > 0 {
				buf = append(buf, tmp[:n]...)
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
		audioData = buf
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to generate audio: %w", err)
	}
	
	return audioData, nil
}

func (s *OpenAIService) GradeDictation(ctx context.Context, reference, userInput string) (*types.Grade, error) {
	prompt := buildGradingPrompt(reference, userInput)
	
	var response openai.ChatCompletionResponse
	err := s.withRetry(ctx, func() error {
		var err error
		response, err = s.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT4oMini,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: gradingSystemPrompt,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
				Temperature:      0,
				MaxTokens:        500,
				ResponseFormat:   &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
			},
		)
		return err
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to grade dictation: %w", err)
	}
	
	if len(response.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}
	
	grade, err := parseGradingResponse(response.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse grading response: %w", err)
	}
	
	return grade, nil
}

func buildSentencePrompt(config *types.Config) string {
	levelDesc := fmt.Sprintf("TOEIC %d level", config.Level)
	return fmt.Sprintf(
		"Generate a single English sentence for dictation practice.\n"+
			"Requirements:\n"+
			"- Topic: %s\n"+
			"- Difficulty: %s (intermediate complexity)\n"+
			"- Length: approximately %d words\n"+
			"- Natural, conversational English\n"+
			"- No quotation marks\n"+
			"- Return only the sentence, nothing else",
		config.Topic, levelDesc, config.Words,
	)
}

const gradingSystemPrompt = `You are an English teacher grading dictation exercises for Japanese learners.
Compare the reference sentence with the user's input and provide detailed feedback.

Return a JSON object with this structure:
{
  "wer": 0.15,
  "score": 85,
  "mistakes": [
    {
      "position": 3,
      "expected": "going",
      "actual": "gonna",
      "type": "spelling"
    }
  ],
  "japanese_explanation": "全体的によくできています。「going」を「gonna」と書いてしまいましたが、これは口語的な表現です。正式な文章では「going」を使いましょう。",
  "alternative_expressions": [
    "I will go to the store tomorrow",
    "I'm planning to visit the store tomorrow"
  ]
}

Calculate WER (Word Error Rate) as: (insertions + deletions + substitutions) / total_words
Score is 100 * (1 - WER), rounded to nearest integer.
Provide explanations in natural Japanese.`

func buildGradingPrompt(reference, userInput string) string {
	return fmt.Sprintf(
		"Reference sentence: %s\n"+
			"User input: %s\n\n"+
			"Grade the dictation and provide feedback.",
		reference, userInput,
	)
}