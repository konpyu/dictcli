package service

import (
	"encoding/json"
	"fmt"

	"github.com/konpyu/dictcli/internal/types"
)

func parseGradingResponse(jsonStr string) (*types.Grade, error) {
	var grade types.Grade
	
	if err := json.Unmarshal([]byte(jsonStr), &grade); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	
	if grade.WER < 0 || grade.WER > 1 {
		grade.WER = 0
	}
	
	if grade.Score < 0 || grade.Score > 100 {
		grade.Score = int(100 * (1 - grade.WER))
	}
	
	return &grade, nil
}