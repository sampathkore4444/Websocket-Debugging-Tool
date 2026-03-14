package schemas

import (
	"wsinspect/backend/models"
)

type CreateSessionRequest struct {
	ClientIP   string `json:"client_ip" binding:"required"`
	ServerHost string `json:"server_host" binding:"required"`
}

type UpdateSessionRequest struct {
	Status models.SessionStatus `json:"status"`
}

type SessionResponse struct {
	ID           uint                   `json:"id"`
	ConnectionID string                 `json:"connection_id"`
	ClientIP     string                 `json:"client_ip"`
	ServerHost   string                 `json:"server_host"`
	Status       models.SessionStatus   `json:"status"`
	MessageCount int                    `json:"message_count"`
	StartTime    string                 `json:"start_time"`
	EndTime      *string                `json:"end_time,omitempty"`
}

func ToSessionResponse(session *models.Session) SessionResponse {
	resp := SessionResponse{
		ID:           session.ID,
		ConnectionID: session.ConnectionID,
		ClientIP:     session.ClientIP,
		ServerHost:   session.ServerHost,
		Status:       session.Status,
		MessageCount: session.MessageCount,
		StartTime:    session.StartTime.Format("2006-01-02 15:04:05"),
	}

	if session.EndTime != nil {
		endTime := session.EndTime.Format("2006-01-02 15:04:05")
		resp.EndTime = &endTime
	}

	return resp
}
