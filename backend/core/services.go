package core

import (
	"wsinspect/backend/services"

	"gorm.io/gorm"
)

func NewProxyService(db *gorm.DB) *services.ProxyService {
	return services.NewProxyService(db)
}

func NewSessionService(db *gorm.DB) *services.SessionService {
	return services.NewSessionService(db)
}

func NewReplayService(db *gorm.DB) *services.ReplayService {
	return services.NewReplayService(db)
}

func NewFuzzService(db *gorm.DB) *services.FuzzService {
	return services.NewFuzzService(db)
}
