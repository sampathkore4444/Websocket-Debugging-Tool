package services

import (
	"testing"
	"time"

	"wsinspect/backend/models"
)

func TestCreateSession(t *testing.T) {
	service := NewSessionService(nil)

	session := &models.Session{
		Name:        "Test Session",
		TargetURL:   "wss://example.com/ws",
		Description: "Test description",
		CreatedAt:   time.Now(),
	}

	createdSession, err := service.CreateSession(session)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if createdSession.ID == "" {
		t.Error("Session ID should not be empty")
	}

	if createdSession.Name != session.Name {
		t.Errorf("Expected name %s, got %s", session.Name, createdSession.Name)
	}
}

func TestGetSession(t *testing.T) {
	service := NewSessionService(nil)

	session := &models.Session{
		Name:      "Test Session",
		TargetURL: "wss://example.com/ws",
		CreatedAt: time.Now(),
	}

	created, _ := service.CreateSession(session)

	retrieved, err := service.GetSession(created.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, retrieved.ID)
	}
}

func TestListSessions(t *testing.T) {
	service := NewSessionService(nil)

	sessions := []*models.Session{
		{Name: "Session 1", TargetURL: "wss://example1.com/ws", CreatedAt: time.Now()},
		{Name: "Session 2", TargetURL: "wss://example2.com/ws", CreatedAt: time.Now()},
		{Name: "Session 3", TargetURL: "wss://example3.com/ws", CreatedAt: time.Now()},
	}

	for _, s := range sessions {
		service.CreateSession(s)
	}

	list, err := service.ListSessions()
	if err != nil {
		t.Fatalf("Failed to list sessions: %v", err)
	}

	if len(list) < len(sessions) {
		t.Errorf("Expected at least %d sessions, got %d", len(sessions), len(list))
	}
}

func TestDeleteSession(t *testing.T) {
	service := NewSessionService(nil)

	session := &models.Session{
		Name:      "Test Session",
		TargetURL: "wss://example.com/ws",
		CreatedAt: time.Now(),
	}

	created, _ := service.CreateSession(session)

	err := service.DeleteSession(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	_, err = service.GetSession(created.ID)
	if err == nil {
		t.Error("Session should not exist after deletion")
	}
}

func TestUpdateSessionStatus(t *testing.T) {
	service := NewSessionService(nil)

	session := &models.Session{
		Name:      "Test Session",
		TargetURL: "wss://example.com/ws",
		CreatedAt: time.Now(),
	}

	created, _ := service.CreateSession(session)

	err := service.UpdateSessionStatus(created.ID, "connected")
	if err != nil {
		t.Fatalf("Failed to update session status: %v", err)
	}

	updated, _ := service.GetSession(created.ID)
	if updated.Status != "connected" {
		t.Errorf("Expected status 'connected', got '%s'", updated.Status)
	}
}
