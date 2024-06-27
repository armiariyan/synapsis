package handler

import (
	"github.com/armiariyan/synapsis/internal/infrastructure/container"
	"gorm.io/gorm"
)

type Handler struct {
	synapsisDB         *gorm.DB
	healthCheckHandler *healthCheckHandler
	userHandler        *userHandler
}

func SetupHandler(container *container.Container) *Handler {
	return &Handler{
		synapsisDB:         container.SynapsisDB,
		healthCheckHandler: NewHealthCheckHandler().SetHealthCheckService(container.HealthCheckService).Validate(),
		userHandler:        NewUserHandler().SetUserService(container.UserService).Validate(),
	}
}

func (h *Handler) Validate() *Handler {
	if h.synapsisDB == nil {
		panic("synapsisDB is nil")
	}
	if h.healthCheckHandler == nil {
		panic("healthCheckHandler is nil")
	}
	if h.userHandler == nil {
		panic("userHandler is nil")
	}
	return h
}
