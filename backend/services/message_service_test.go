package services

import (
	"testing"
	"time"

	"wsinspect/backend/models"
)

func TestSaveMessage(t *testing.T) {
	service := NewMessageService(nil)

	message := &models.Message{
		SessionID:  "test-session",
		Direction:  "outgoing",
		Content:    `{"type": "hello", "data": "test"}`,
		ContentType: "json",
		Timestamp:  time.Now(),
	}

	savedMessage, err := service.SaveMessage(message)
	if err != nil {
		t.Fatalf("Failed to save message: %v", err)
	}

	if savedMessage.ID == 0 {
		t.Error("Message ID should not be 0")
	}
}

func TestGetMessagesBySession(t *testing.T) {
	service := NewMessageService(nil)

	sessionID := "test-session-1"

	messages := []*models.Message{
		{SessionID: sessionID, Direction: "incoming", Content: "msg1", Timestamp: time.Now()},
		{SessionID: sessionID, Direction: "outgoing", Content: "msg2", Timestamp: time.Now()},
		{SessionID: sessionID, Direction: "incoming", Content: "msg3", Timestamp: time.Now()},
		{SessionID: "other-session", Direction: "incoming", Content: "msg4", Timestamp: time.Now()},
	}

	for _, m := range messages {
		service.SaveMessage(m)
	}

	sessionMessages, err := service.GetMessagesBySession(sessionID)
	if err != nil {
		t.Fatalf("Failed to get messages: %v", err)
	}

	if len(sessionMessages) != 3 {
		t.Errorf("Expected 3 messages for session %s, got %d", sessionID, len(sessionMessages))
	}
}

func TestClearSessionMessages(t *testing.T) {
	service := NewMessageService(nil)

	sessionID := "test-session-clear"

	service.SaveMessage(&models.Message{
		SessionID: sessionID, Direction: "incoming", Content: "msg1", Timestamp: time.Now(),
	})
	service.SaveMessage(&models.Message{
		SessionID: sessionID, Direction: "outgoing", Content: "msg2", Timestamp: time.Now(),
	})

	err := service.ClearSessionMessages(sessionID)
	if err != nil {
		t.Fatalf("Failed to clear messages: %v", err)
	}

	messages, _ := service.GetMessagesBySession(sessionID)
	if len(messages) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(messages))
	}
}

func TestFilterMessages(t *testing.T) {
	service := NewMessageService(nil)

	sessionID := "test-session-filter"

	service.SaveMessage(&models.Message{
		SessionID: sessionID, Direction: "incoming", Content: `{"type": "ping"}`, ContentType: "json", Timestamp: time.Now(),
	})
	service.SaveMessage(&models.Message{
		SessionID: sessionID, Direction: "outgoing", Content: `{"type": "pong"}`, ContentType: "json", Timestamp: time.Now(),
	})
	service.SaveMessage(&models.Message{
		SessionID: sessionID, Direction: "incoming", Content: `{"type": "data", "value": 123}`, ContentType: "json", Timestamp: time.Now(),
	})

	filtered, err := service.FilterMessages(sessionID, "incoming")
	if err != nil {
		t.Fatalf("Failed to filter messages: %v", err)
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 incoming messages, got %d", len(filtered))
	}
}

func TestSearchMessages(t *testing.T) {
	service := NewMessageService(nil)

	sessionID := "test-session-search"

	service.SaveMessage(&models.Message{
		SessionID: sessionID, Direction: "incoming", Content: "Hello World", Timestamp: time.Now(),
	})
	service.SaveMessage(&models.Message{
		SessionID: sessionID, Direction: "outgoing", Content: "Goodbye World", Timestamp: time.Now(),
	})
	service.SaveMessage(&models.Message{
		SessionID: sessionID, Direction: "incoming", Content: "Test Message", Timestamp: time.Now(),
	})

	results, err := service.SearchMessages(sessionID, "World")
	if err != nil {
		t.Fatalf("Failed to search messages: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results containing 'World', got %d", len(results))
	}
}
