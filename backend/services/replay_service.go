package services

import (
	"log"
	"time"

	"wsinspect/backend/models"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ReplayMode string

const (
	ReplayModeExact  ReplayMode = "exact"
	ReplayModeFast  ReplayMode = "fast"
	ReplayModeEdit  ReplayMode = "edited"
)

type ReplayService struct {
	db *gorm.DB
}

func NewReplayService(db *gorm.DB) *ReplayService {
	return &ReplayService{db: db}
}

type ReplayRequest struct {
	SessionID uint      `json:"session_id" binding:"required"`
	Mode      ReplayMode `json:"mode"`
	TargetURL string    `json:"target_url" binding:"required"`
}

type ReplayResult struct {
	SessionID    uint      `json:"session_id"`
	Mode         ReplayMode `json:"mode"`
	TotalMessages int       `json:"total_messages"`
	SentMessages  int       `json:"sent_messages"`
	FailedMessages int      `json:"failed_messages"`
	Duration     time.Duration `json:"duration"`
	Success      bool      `json:"success"`
}

func (s *ReplayService) ReplaySession(req *ReplayRequest) (*ReplayResult, error) {
	startTime := time.Now()
	
	// Get all messages for the session
	var messages []models.Message
	if err := s.db.Where("session_id = ?", req.SessionID).
		Order("timestamp ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}

	result := &ReplayResult{
		SessionID:    req.SessionID,
		Mode:         req.Mode,
		TotalMessages: len(messages),
	}

	if len(messages) == 0 {
		result.Success = true
		return result, nil
	}

	// Connect to target server
	serverConn, _, err := websocket.DefaultDialer.Dial(req.TargetURL, nil)
	if err != nil {
		result.FailedMessages = len(messages)
		result.Success = false
		return result, err
	}
	defer serverConn.Close()

	var lastTimestamp time.Time
	for i, msg := range messages {
		// For exact mode, respect timing
		if req.Mode == ReplayModeExact && i > 0 {
			delay := msg.Timestamp.Sub(lastTimestamp)
			if delay > 0 {
				time.Sleep(delay)
			}
		}

		// Skip timing delays for fast mode
		if req.Mode == ReplayModeFast {
			// No delay
		}

		// Convert direction to websocket message type
		messageType := websocket.TextMessage
		if msg.Opcode > 0 {
			messageType = msg.Opcode
		}

		// For edited mode, use modified payload if available
		payload := msg.Payload
		if req.Mode == ReplayModeEdit && msg.IsModified && msg.OriginalPayload != nil {
			payload = *msg.OriginalPayload
		}

		// Send message to server
		if err := serverConn.WriteMessage(messageType, []byte(payload)); err != nil {
			result.FailedMessages++
			log.Printf("Failed to send message %d: %v", i, err)
			continue
		}

		result.SentMessages++
		lastTimestamp = msg.Timestamp
	}

	result.Duration = time.Since(startTime)
	result.Success = result.FailedMessages == 0

	return result, nil
}

func (s *ReplayService) GetReplayableMessages(sessionID uint) ([]models.Message, error) {
	var messages []models.Message
	err := s.db.Where("session_id = ? AND direction = ?", sessionID, models.MessageDirectionClientToServer).
		Order("timestamp ASC").
		Find(&messages).Error
	return messages, err
}
