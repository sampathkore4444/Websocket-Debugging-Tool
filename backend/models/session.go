package models

import (
	"time"
)

type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "active"
	SessionStatusClosed  SessionStatus = "closed"
	SessionStatusError   SessionStatus = "error"
)

type Session struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	ConnectionID string         `gorm:"uniqueIndex;not null" json:"connection_id"`
	ClientIP     string         `json:"client_ip"`
	ServerHost   string         `json:"server_host"`
	StartTime    time.Time     `gorm:"autoCreateTime" json:"start_time"`
	EndTime      *time.Time    `json:"end_time,omitempty"`
	Status       SessionStatus `gorm:"default:active" json:"status"`
	MessageCount int            `gorm:"default:0" json:"message_count"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

func (Session) TableName() string {
	return "sessions"
}
