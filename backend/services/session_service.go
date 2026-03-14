package services

import (
	"time"

	"wsinspect/backend/models"
	"wsinspect/backend/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) CreateSession(req *schemas.CreateSessionRequest) (*models.Session, error) {
	session := &models.Session{
		ConnectionID: uuid.New().String(),
		ClientIP:     req.ClientIP,
		ServerHost:   req.ServerHost,
		Status:       models.SessionStatusActive,
		StartTime:    time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) GetSession(id uint) (*models.Session, error) {
	var session models.Session
	if err := s.db.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionService) GetSessionByConnectionID(connectionID string) (*models.Session, error) {
	var session models.Session
	if err := s.db.Where("connection_id = ?", connectionID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionService) ListSessions(limit, offset int) ([]models.Session, int64, error) {
	var sessions []models.Session
	var total int64

	s.db.Model(&models.Session{}).Count(&total)
	
	if err := s.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (s *SessionService) UpdateSession(id uint, req *schemas.UpdateSessionRequest) (*models.Session, error) {
	var session models.Session
	if err := s.db.First(&session, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Status != "" {
		updates["status"] = req.Status
		if req.Status == models.SessionStatusClosed {
			now := time.Now()
			updates["end_time"] = &now
		}
	}

	if err := s.db.Model(&session).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionService) DeleteSession(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete all messages for this session
		if err := tx.Where("session_id = ?", id).Delete(&models.Message{}).Error; err != nil {
			return err
		}
		// Delete all connections for this session
		if err := tx.Where("session_id = ?", id).Delete(&models.Connection{}).Error; err != nil {
			return err
		}
		// Delete the session
		return tx.Delete(&models.Session{}, id).Error
	})
}

func (s *SessionService) IncrementMessageCount(id uint) error {
	return s.db.Model(&models.Session{}).Where("id = ?", id).
		UpdateColumn("message_count", gorm.Expr("message_count + ?", 1)).Error
}
