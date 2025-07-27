package service

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/konpyu/dictcli/internal/types"
)

func TestNewOpenAIService(t *testing.T) {
	originalKey := os.Getenv("OPENAI_API_KEY")
	defer func() {
		_ = os.Setenv("OPENAI_API_KEY", originalKey)
	}()
	
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "with API key",
			apiKey:  "test-key",
			wantErr: false,
		},
		{
			name:    "without API key",
			apiKey:  "",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv("OPENAI_API_KEY", tt.apiKey); err != nil {
				t.Fatalf("Failed to set env var: %v", err)
			}
			
			_, err := NewOpenAIService(false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOpenAIService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildSentencePrompt(t *testing.T) {
	config := &types.Config{
		Topic: types.TopicBusiness,
		Level: 700,
		Words: 15,
	}
	
	prompt := buildSentencePrompt(config)
	
	if prompt == "" {
		t.Error("buildSentencePrompt() returned empty string")
	}
	
	expectedParts := []string{
		"Business",
		"TOEIC 700",
		"15 words",
	}
	
	for _, part := range expectedParts {
		if !strings.Contains(prompt, part) {
			t.Errorf("buildSentencePrompt() missing expected part: %s", part)
		}
	}
}

func TestBuildGradingPrompt(t *testing.T) {
	reference := "I am going to the store"
	userInput := "I am going to the shop"
	
	prompt := buildGradingPrompt(reference, userInput)
	
	if !strings.Contains(prompt, reference) {
		t.Errorf("buildGradingPrompt() missing reference sentence")
	}
	
	if !strings.Contains(prompt, userInput) {
		t.Errorf("buildGradingPrompt() missing user input")
	}
}

func TestWithRetry(t *testing.T) {
	service := &OpenAIService{debug: false}
	ctx := context.Background()
	
	t.Run("success on first try", func(t *testing.T) {
		attempts := 0
		err := service.withRetry(ctx, func() error {
			attempts++
			return nil
		})
		
		if err != nil {
			t.Errorf("withRetry() unexpected error: %v", err)
		}
		
		if attempts != 1 {
			t.Errorf("withRetry() attempts = %d, want 1", attempts)
		}
	})
	
	t.Run("success after retry", func(t *testing.T) {
		attempts := 0
		err := service.withRetry(ctx, func() error {
			attempts++
			if attempts < 2 {
				// Return a timeout error which should be retryable
				return &testError{"timeout error", true}
			}
			return nil
		})
		
		if err != nil {
			t.Errorf("withRetry() unexpected error: %v", err)
		}
		
		if attempts != 2 {
			t.Errorf("withRetry() attempts = %d, want 2", attempts)
		}
	})
	
	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		
		// Sleep to ensure context times out
		time.Sleep(2 * time.Millisecond)
		
		err := service.withRetry(ctx, func() error {
			// Check if context is expired
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return &testError{"error", false}
			}
		})
		
		// The error should be context-related and formatted by formatError
		if err == nil {
			t.Errorf("withRetry() expected error, got nil")
		} else if !strings.Contains(err.Error(), "タイムアウト") && !strings.Contains(err.Error(), "キャンセル") {
			t.Errorf("withRetry() error = %v, expected to contain timeout or cancellation message", err)
		}
	})
}

type testError struct {
	msg string
	timeout bool
}

func (e *testError) Error() string {
	return e.msg
}

func (e *testError) Temporary() bool {
	return false // deprecated, not used
}

func (e *testError) Timeout() bool {
	return e.timeout
}


func TestOpenAIServiceIntegration(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("Skipping integration test: OPENAI_API_KEY not set")
	}
	
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	service, err := NewOpenAIService(true)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	config := &types.Config{
		Topic: types.TopicBusiness,
		Level: 700,
		Words: 10,
	}
	
	t.Run("GenerateSentence", func(t *testing.T) {
		sentence, err := service.GenerateSentence(ctx, config)
		if err != nil {
			t.Errorf("GenerateSentence() error: %v", err)
		}
		
		if sentence == "" {
			t.Error("GenerateSentence() returned empty sentence")
		}
		
		t.Logf("Generated sentence: %s", sentence)
	})
}