package services

import (
	"log"
	"net/http"
	"sync"
	"time"

	"wsinspect/backend/models"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ProxyService struct {
	db           *gorm.DB
	upgrader     websocket.Upgrader
	connections  map[string]*ProxyConnection
	mu           sync.RWMutex
	sessionSvc   *SessionService
	messageSvc   *MessageService
}

type ProxyConnection struct {
	ID          string
	SessionID   uint
	ClientConn  *websocket.Conn
	ServerConn  *websocket.Conn
	TargetURL   string
	IsActive    bool
	StartedAt   time.Time
}

func NewProxyService(db *gorm.DB) *ProxyService {
	ps := &ProxyService{
		db:          db,
		connections: make(map[string]*ProxyConnection),
	}
	ps.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return ps
}

func (s *ProxyService) SetServices(sessionSvc *SessionService, messageSvc *MessageService) {
	s.sessionSvc = sessionSvc
	s.messageSvc = messageSvc
}

func (s *ProxyService) HandleWebSocket(w http.ResponseWriter, r *http.Request, targetURL string) error {
	// Create new session for this connection
	session, err := s.sessionSvc.CreateSession(&CreateSessionRequest{
		ClientIP:   r.RemoteAddr,
		ServerHost: targetURL,
	})
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		return err
	}

	// Upgrade client connection
	clientConn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade client connection: %v", err)
		return err
	}

	// Connect to target server
	serverConn, _, err := websocket.DefaultDialer.Dial(targetURL, nil)
	if err != nil {
		log.Printf("Failed to connect to target server: %v", err)
		clientConn.Close()
		return err
	}

	// Create proxy connection
	proxyConn := &ProxyConnection{
		ID:         session.ConnectionID,
		SessionID:  session.ID,
		ClientConn: clientConn,
		ServerConn: serverConn,
		TargetURL:  targetURL,
		IsActive:   true,
		StartedAt:  time.Now(),
	}

	s.mu.Lock()
	s.connections[session.ConnectionID] = proxyConn
	s.mu.Unlock()

	// Start bidirectional proxy
	go s.proxyToServer(proxyConn)
	go s.proxyToClient(proxyConn)

	return nil
}

func (s *ProxyService) proxyToServer(pc *ProxyConnection) {
	defer func() {
		pc.ClientConn.Close()
		s.closeConnection(pc.ID)
	}()

	for {
		messageType, message, err := pc.ClientConn.ReadMessage()
		if err != nil {
			break
		}

		// Record message to database
		if s.messageSvc != nil {
			s.messageSvc.RecordMessage(pc.SessionID, models.MessageDirectionClientToServer, message, messageType)
			s.sessionSvc.IncrementMessageCount(pc.SessionID)
		}

		// Forward to server
		if err := pc.ServerConn.WriteMessage(messageType, message); err != nil {
			break
		}
	}
}

func (s *ProxyService) proxyToClient(pc *ProxyConnection) {
	defer func() {
		pc.ServerConn.Close()
		s.closeConnection(pc.ID)
	}()

	for {
		messageType, message, err := pc.ServerConn.ReadMessage()
		if err != nil {
			break
		}

		// Record message to database
		if s.messageSvc != nil {
			s.messageSvc.RecordMessage(pc.SessionID, models.MessageDirectionServerToClient, message, messageType)
			s.sessionSvc.IncrementMessageCount(pc.SessionID)
		}

		// Forward to client
		if err := pc.ClientConn.WriteMessage(messageType, message); err != nil {
			break
		}
	}
}

func (s *ProxyService) closeConnection(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conn, ok := s.connections[id]; ok {
		conn.IsActive = false
		s.sessionSvc.UpdateSession(conn.SessionID, &UpdateSessionStatusRequest{
			Status: models.SessionStatusClosed,
		})
		delete(s.connections, id)
	}
}

func (s *ProxyService) GetActiveConnections() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.connections)
}

func (s *ProxyService) GetConnection(id string) (*ProxyConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conn, ok := s.connections[id]
	return conn, ok
}

func (s *ProxyService) StopConnection(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conn, ok := s.connections[id]; ok {
		conn.ClientConn.Close()
		conn.ServerConn.Close()
		conn.IsActive = false
		return nil
	}

	return nil
}

// Helper types for session service
type CreateSessionRequest struct {
	ClientIP   string
	ServerHost string
}

type UpdateSessionStatusRequest struct {
	Status models.SessionStatus
}
