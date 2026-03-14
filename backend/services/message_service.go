package services

import (
	"encoding/json"
	"time"

	"wsinspect/backend/models"
	"wsinspect/backend/schemas"

	"gorm.io/gorm"
)

type MessageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) CreateMessage(req *schemas.CreateMessageRequest) (*models.Message, error) {
	payloadFormat := req.PayloadFormat
	if payloadFormat == "" {
		payloadFormat = models.PayloadFormatText
		// Try to detect JSON
		if isValidJSON(req.Payload) {
			payloadFormat = models.PayloadFormatJSON
		}
	}

	message := &models.Message{
		SessionID:     req.SessionID,
		Direction:     req.Direction,
		Payload:       req.Payload,
		PayloadFormat: payloadFormat,
		Opcode:        req.Opcode,
	}

	if err := s.db.Create(message).Error; err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) GetMessage(id uint) (*models.Message, error) {
	var message models.Message
	if err := s.db.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (s *MessageService) GetMessagesBySession(sessionID uint, limit, offset int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	s.db.Model(&models.Message{}).Where("session_id = ?", sessionID).Count(&total)

	if err := s.db.Where("session_id = ?", sessionID).
		Order("timestamp ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (s *MessageService) UpdateMessage(id uint, req *schemas.UpdateMessageRequest) (*models.Message, error) {
	var message models.Message
	if err := s.db.First(&message, id).Error; err != nil {
		return nil, err
	}

	// Store original payload before modification
	originalPayload := message.Payload

	updates := map[string]interface{}{
		"payload":            req.Payload,
		"is_modified":        true,
		"original_payload":   originalPayload,
	}

	// Update format if payload changed
	if isValidJSON(req.Payload) {
		updates["payload_format"] = models.PayloadFormatJSON
	} else {
		updates["payload_format"] = models.PayloadFormatText
	}

	if err := s.db.Model(&message).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *MessageService) DeleteMessage(id uint) error {
	return s.db.Delete(&models.Message{}, id).Error
}

func (s *MessageService) DeleteMessagesBySession(sessionID uint) error {
	return s.db.Where("session_id = ?", sessionID).Delete(&models.Message{}).Error
}

func (s *MessageService) RecordMessage(sessionID uint, direction models.MessageDirection, payload []byte, opcode int) (*models.Message, error) {
	payloadStr := string(payload)
	payloadFormat := models.PayloadFormatText

	if isValidJSON(payloadStr) {
		payloadFormat = models.PayloadFormatJSON
	}

	message := &models.Message{
		SessionID:     sessionID,
		Direction:     direction,
		Payload:       payloadStr,
		PayloadFormat: payloadFormat,
		Opcode:        opcode,
		Timestamp:     time.Now(),
	}

	if err := s.db.Create(message).Error; err != nil {
		return nil, err
	}

	return message, nil
}

// Helper function to validate JSON
func isValidJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
