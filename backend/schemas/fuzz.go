package schemas

import (
	"wsinspect/backend/models"
)

type CreateFuzzTestRequest struct {
	SessionID uint                   `json:"session_id" binding:"required"`
	Name      string                `json:"name" binding:"required"`
	Strategy  models.FuzzStrategy   `json:"strategy"`
	Template  string                `json:"template" binding:"required"`
}

type UpdateFuzzTestRequest struct {
	Status models.FuzzStatus `json:"status"`
}

type FuzzTestResponse struct {
	ID           uint                 `json:"id"`
	SessionID    uint                 `json:"session_id"`
	Name         string               `json:"name"`
	Strategy     models.FuzzStrategy  `json:"strategy"`
	Status       models.FuzzStatus    `json:"status"`
	Template     string               `json:"template"`
	TestCount    int                  `json:"test_count"`
	SuccessCount int                  `json:"success_count"`
	FailCount    int                  `json:"fail_count"`
	StartedAt    *string              `json:"started_at,omitempty"`
	CompletedAt  *string              `json:"completed_at,omitempty"`
}

func ToFuzzTestResponse(fuzz *models.FuzzTest) FuzzTestResponse {
	resp := FuzzTestResponse{
		ID:           fuzz.ID,
		SessionID:    fuzz.SessionID,
		Name:         fuzz.Name,
		Strategy:     fuzz.Strategy,
		Status:       fuzz.Status,
		Template:     fuzz.Template,
		TestCount:    fuzz.TestCount,
		SuccessCount: fuzz.SuccessCount,
		FailCount:    fuzz.FailCount,
	}

	if fuzz.StartedAt != nil {
		startedAt := fuzz.StartedAt.Format("2006-01-02 15:04:05")
		resp.StartedAt = &startedAt
	}

	if fuzz.CompletedAt != nil {
		completedAt := fuzz.CompletedAt.Format("2006-01-02 15:04:05")
		resp.CompletedAt = &completedAt
	}

	return resp
}
