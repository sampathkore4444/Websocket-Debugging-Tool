package models

import (
	"time"
)

type FuzzStatus string
type FuzzStrategy string

const (
	FuzzStatusPending   FuzzStatus = "pending"
	FuzzStatusRunning   FuzzStatus = "running"
	FuzzStatusCompleted FuzzStatus = "completed"
	FuzzStatusFailed    FuzzStatus = "failed"

	FuzzStrategyRandom     FuzzStrategy = "random"
	FuzzStrategyMutation  FuzzStrategy = "mutation"
	FuzzStrategyBoundary  FuzzStrategy = "boundary"
	FuzzStrategyInvalid   FuzzStrategy = "invalid"
)

type FuzzTest struct {
	ID           uint         `gorm:"primarykey" json:"id"`
	SessionID    uint         `gorm:"index" json:"session_id"`
	Name         string       `json:"name"`
	Strategy     FuzzStrategy `gorm:"default:random" json:"strategy"`
	Status       FuzzStatus   `gorm:"default:pending" json:"status"`
	Template     string       `gorm:"type:text" json:"template"`
	TestCount    int          `gorm:"default:0" json:"test_count"`
	SuccessCount int          `gorm:"default:0" json:"success_count"`
	FailCount    int          `gorm:"default:0" json:"fail_count"`
	StartedAt    *time.Time   `json:"started_at,omitempty"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (FuzzTest) TableName() string {
	return "fuzz_tests"
}
