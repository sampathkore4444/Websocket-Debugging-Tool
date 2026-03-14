package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ScriptStep struct {
	Send    string `json:"send,omitempty"`
	Wait    int    `json:"wait,omitempty"`
	Connect string `json:"connect,omitempty"`
}

type Script struct {
	Steps []ScriptStep `json:"steps"`
}

type ScriptRunner struct {
	db           *gorm.DB
	sessionSvc   *SessionService
	messageSvc   *MessageService
}

func NewScriptRunner(db *gorm.DB, sessionSvc *SessionService, messageSvc *MessageService) *ScriptRunner {
	return &ScriptRunner{
		db:         db,
		sessionSvc: sessionSvc,
		messageSvc: messageSvc,
	}
}

type ScriptRequest struct {
	ScriptYAML string `json:"script_yaml" binding:"required"`
	TargetURL  string `json:"target_url" binding:"required"`
}

type ScriptResult struct {
	TotalSteps  int       `json:"total_steps"`
	Executed    int       `json:"executed"`
	Failed      int       `json:"failed"`
	Duration    time.Duration `json:"duration"`
	Success     bool      `json:"success"`
}

func (s *ScriptRunner) RunScript(req *ScriptRequest) (*ScriptResult, error) {
	startTime := time.Now()
	result := &ScriptResult{
		Success: true,
	}

	// Parse YAML script (simplified - in production use proper YAML parsing)
	var script Script
	if err := json.Unmarshal([]byte(req.ScriptYAML), &script); err != nil {
		// Try parsing as simple JSON
		if err := json.Unmarshal([]byte(req.ScriptYAML), &script); err != nil {
			result.Success = false
			result.Failed = 1
			return result, fmt.Errorf("failed to parse script: %w", err)
		}
	}

	result.TotalSteps = len(script.Steps)

	// Connect to target server
	serverConn, _, err := websocket.DefaultDialer.Dial(req.TargetURL, nil)
	if err != nil {
		result.Success = false
		result.Failed = len(script.Steps)
		return result, err
	}
	defer serverConn.Close()

	// Execute each step
	for _, step := range script.Steps {
		// Handle connect step
		if step.Connect != "" {
			serverConn, _, err = websocket.DefaultDialer.Dial(step.Connect, nil)
			if err != nil {
				result.Failed++
				continue
			}
		}

		// Handle send step
		if step.Send != "" {
			messageType := websocket.TextMessage
			if err := serverConn.WriteMessage(messageType, []byte(step.Send)); err != nil {
				result.Failed++
				continue
			}
			result.Executed++
		}

		// Handle wait step
		if step.Wait > 0 {
			time.Sleep(time.Duration(step.Wait) * time.Millisecond)
		}
	}

	result.Duration = time.Since(startTime)

	return result, nil
}

// Script example:
// {
//   "steps": [
//     {"send": "{\"type\": \"login\"}"},
//     {"wait": 1000},
//     {"send": "{\"type\": \"subscribe\", \"channel\": \"orders\"}"}
//   ]
// }
