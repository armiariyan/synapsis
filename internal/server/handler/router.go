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
			users.POST("/login", h.userHandler.Login)
			users.POST("/checkout", h.userHandler.Checkout, h.Authentication)

		}

		products := v1.Group("/products", h.Authentication)
		{
			products.GET("", h.productHandler.FindAllProductsByCategory)
			products.POST("/add-to-cart", h.productHandler.AddToCart)
		}

		carts := v1.Group("/carts", h.Authentication)
		{
			carts.GET("", h.cartHandler.FindAll)
			carts.DELETE("/:uuid", h.cartHandler.DeleteProductFromCarts)
		}
	}
}
