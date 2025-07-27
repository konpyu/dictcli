package service

import (
	"testing"

	"github.com/konpyu/dictcli/internal/types"
)

func TestParseGradingResponse(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    *types.Grade
		wantErr bool
	}{
		{
			name: "valid response",
			jsonStr: `{
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
				"japanese_explanation": "よくできています。",
				"alternative_expressions": ["I will go to the store"]
			}`,
			want: &types.Grade{
				WER:   0.15,
				Score: 85,
				Mistakes: []types.Mistake{
					{
						Position: 3,
						Expected: "going",
						Actual:   "gonna",
						Type:     "spelling",
					},
				},
				JapaneseExplanation:    "よくできています。",
				AlternativeExpressions: []string{"I will go to the store"},
			},
			wantErr: false,
		},
		{
			name: "invalid WER",
			jsonStr: `{
				"wer": 1.5,
				"score": 85,
				"mistakes": [],
				"japanese_explanation": "説明",
				"alternative_expressions": []
			}`,
			want: &types.Grade{
				WER:                    0,
				Score:                  85,
				Mistakes:               []types.Mistake{},
				JapaneseExplanation:    "説明",
				AlternativeExpressions: []string{},
			},
			wantErr: false,
		},
		{
			name: "invalid score",
			jsonStr: `{
				"wer": 0.2,
				"score": 150,
				"mistakes": [],
				"japanese_explanation": "説明",
				"alternative_expressions": []
			}`,
			want: &types.Grade{
				WER:                    0.2,
				Score:                  80,
				Mistakes:               []types.Mistake{},
				JapaneseExplanation:    "説明",
				AlternativeExpressions: []string{},
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			jsonStr: `{invalid json}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty string",
			jsonStr: ``,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGradingResponse(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGradingResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr {
				return
			}
			
			if got.WER != tt.want.WER {
				t.Errorf("parseGradingResponse() WER = %v, want %v", got.WER, tt.want.WER)
			}
			
			if got.Score != tt.want.Score {
				t.Errorf("parseGradingResponse() Score = %v, want %v", got.Score, tt.want.Score)
			}
			
			if got.JapaneseExplanation != tt.want.JapaneseExplanation {
				t.Errorf("parseGradingResponse() JapaneseExplanation = %v, want %v", 
					got.JapaneseExplanation, tt.want.JapaneseExplanation)
			}
			
			if len(got.Mistakes) != len(tt.want.Mistakes) {
				t.Errorf("parseGradingResponse() Mistakes length = %v, want %v", 
					len(got.Mistakes), len(tt.want.Mistakes))
			}
		})
	}
}