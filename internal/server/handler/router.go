package handler

import (
	"github.com/armiariyan/synapsis/internal/infrastructure/container"
	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo, cnt *container.Container) {
	h := SetupHandler(cnt).Validate()

	e.GET("/", h.healthCheckHandler.HealthCheck)

	v1 := e.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/register", h.userHandler.Register)
		}
	}
}
