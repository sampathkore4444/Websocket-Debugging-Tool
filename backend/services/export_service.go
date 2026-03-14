package services

import (
	"encoding/json"
	"fmt"
	"time"

	"wsinspect/backend/models"

	"gorm.io/gorm"
)

type ExportFormat string

const (
	ExportFormatJSON  ExportFormat = "json"
	ExportFormatNDJSON ExportFormat = "ndjson"
	ExportFormatBinary ExportFormat = "binary"
)

type ExportService struct {
	db *gorm.DB
}

func NewExportService(db *gorm.DB) *ExportService {
	return &ExportService{db: db}
}

type ExportRequest struct {
	SessionID uint         `json:"session_id" binding:"required"`
	Format    ExportFormat `json:"format"`
	Include   string       `json:"include"` // "all", "client-only", "server-only"
}

type ExportData struct {
	Session   *models.Session `json:"session"`
	Messages  []models.Message `json:"messages"`
	ExportTime time.Time      `json:"export_time"`
}

func (s *ExportService) ExportSession(req *ExportRequest) ([]byte, error) {
	format := req.Format
	if format == "" {
		format = ExportFormatJSON
	}

	include := req.Include
	if include == "" {
		include = "all"
	}

	// Get session
	var session models.Session
	if err := s.db.First(&session, req.SessionID).Error; err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Get messages
	var messages []models.Message
	query := s.db.Where("session_id = ?", req.SessionID).Order("timestamp ASC")

	switch include {
	case "client-only":
		query = query.Where("direction = ?", models.MessageDirectionClientToServer)
	case "server-only":
		query = query.Where("direction = ?", models.MessageDirectionServerToClient)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Export based on format
	switch format {
	case ExportFormatJSON:
		return s.exportJSON(&ExportData{
			Session:    &session,
			Messages:   messages,
			ExportTime: time.Now(),
		})
	case ExportFormatNDJSON:
		return s.exportNDJSON(&session, messages)
	default:
		return s.exportJSON(&ExportData{
			Session:    &session,
			Messages:   messages,
			ExportTime: time.Now(),
		})
	}
}

func (s *ExportService) exportJSON(data *ExportData) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func (s *ExportService) exportNDJSON(session *models.Session, messages []models.Message) ([]byte, error) {
	var result []byte
	
	// Session header
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}
	result = append(result, sessionJSON...)
	result = append(result, '\n')

	// Messages
	for _, msg := range messages {
		msgJSON, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		result = append(result, msgJSON...)
		result = append(result, '\n')
	}

	return result, nil
}

func (s *ExportService) ImportSession(data []byte) (*models.Session, error) {
	var exportData ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		return nil, fmt.Errorf("failed to parse import data: %w", err)
	}

	// Create new session with new ID
	newSession := &models.Session{
		ConnectionID: exportData.Session.ConnectionID + "-imported",
		ClientIP:    exportData.Session.ClientIP,
		ServerHost:  exportData.Session.ServerHost,
		Status:      models.SessionStatusClosed,
	}

	if err := s.db.Create(newSession).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Import messages with new session ID
	for _, msg := range exportData.Messages {
		msg.ID = 0
		msg.SessionID = newSession.ID
		if err := s.db.Create(&msg).Error; err != nil {
			continue
		}
	}

	return newSession, nil
}
