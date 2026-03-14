package schemas

import (
	"wsinspect/backend/models"
)

type CreateMessageRequest struct {
	SessionID     uint                    `json:"session_id" binding:"required"`
	Direction     models.MessageDirection  `json:"direction" binding:"required"`
	Payload       string                  `json:"payload" binding:"required"`
	PayloadFormat models.PayloadFormat    `json:"payload_format"`
	Opcode        int                     `json:"opcode"`
}

type UpdateMessageRequest struct {
	Payload string `json:"payload"`
}

type MessageResponse struct {
	ID             uint                     `json:"id"`
	SessionID      uint                     `json:"session_id"`
	Timestamp      string                   `json:"timestamp"`
	Direction      models.MessageDirection  `json:"direction"`
	Opcode         int                      `json:"opcode"`
	Payload        string                   `json:"payload"`
	PayloadFormat  models.PayloadFormat     `json:"payload_format"`
	LatencyMs      *int64                  `json:"latency_ms,omitempty"`
	IsModified     bool                     `json:"is_modified"`
	OriginalPayload *string                 `json:"original_payload,omitempty"`
}

func ToMessageResponse(msg *models.Message) MessageResponse {
	resp := MessageResponse{
		ID:            msg.ID,
		SessionID:     msg.SessionID,
		Timestamp:     msg.Timestamp.Format("2006-01-02 15:04:05.000"),
		Direction:     msg.Direction,
		Opcode:        msg.Opcode,
		Payload:       msg.Payload,
		PayloadFormat: msg.PayloadFormat,
		IsModified:    msg.IsModified,
		LatencyMs:     msg.LatencyMs,
	}

	if msg.OriginalPayload != nil {
		resp.OriginalPayload = msg.OriginalPayload
	}

	return resp
}
