package models

import (
	"time"
)

type ConnectionStatus string

const (
	ConnectionStatusPending   ConnectionStatus = "pending"
	ConnectionStatusActive   ConnectionStatus = "active"
	ConnectionStatusClosed   ConnectionStatus = "closed"
	ConnectionStatusError    ConnectionStatus = "error"
)

type Connection struct {
	ID           uint             `gorm:"primarykey" json:"id"`
	SessionID    uint             `gorm:"index" json:"session_id"`
	ConnectionID string           `gorm:"uniqueIndex;not null" json:"connection_id"`
	ClientIP     string           `json:"client_ip"`
	ClientPort   string           `json:"client_port"`
	ServerHost   string           `json:"server_host"`
	ServerPort   string           `json:"server_port"`
	Protocol     string           `gorm:"default:ws" json:"protocol"`
	Status       ConnectionStatus `gorm:"default:pending" json:"status"`
	IsSecure    bool             `gorm:"default:false" json:"is_secure"`
	StartedAt   time.Time        `gorm:"autoCreateTime" json:"started_at"`
	EndedAt     *time.Time       `json:"ended_at,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func (Connection) TableName() string {
	return "connections"
}
