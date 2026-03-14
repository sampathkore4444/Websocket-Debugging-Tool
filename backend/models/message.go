package models

import (
	"time"
)

type MessageDirection string
type PayloadFormat string

const (
	MessageDirectionClientToServer MessageDirection = "client-to-server"
	MessageDirectionServerToClient MessageDirection = "server-to-client"

	PayloadFormatText   PayloadFormat = "text"
	PayloadFormatJSON   PayloadFormat = "json"
	PayloadFormatBinary PayloadFormat = "binary"
	PayloadFormatHex   PayloadFormat = "hex"
)

type Message struct {
	ID             uint            `gorm:"primarykey" json:"id"`
	SessionID      uint            `gorm:"index;not null" json:"session_id"`
	Timestamp      time.Time       `gorm:"autoCreateTime" json:"timestamp"`
	Direction      MessageDirection `gorm:"not null" json:"direction"`
	Opcode         int             `json:"opcode"`
	Payload        string          `gorm:"type:text" json:"payload"`
	PayloadFormat  PayloadFormat   `gorm:"default:text" json:"payload_format"`
	LatencyMs      *int64          `json:"latency_ms,omitempty"`
	IsModified     bool            `gorm:"default:false" json:"is_modified"`
	OriginalPayload *string         `gorm:"type:text" json:"original_payload,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}
