package handler

import (
	"net/http"

	"github.com/armiariyan/synapsis/internal/pkg/utils"
	cart "github.com/armiariyan/synapsis/internal/usecase/carts"
	"github.com/labstack/echo/v4"
)

type cartHandler struct {
	cartService cart.Service
}

func NewCartHandler() *cartHandler {
	return &cartHandler{}
}

func (h *cartHandler) SetCartService(service cart.Service) *cartHandler {
	h.cartService = service
	return h
}

func (h *cartHandler) Validate() *cartHandler {
	if h.cartService == nil {
		panic("cartService is nil")
	}
	return h
}

func (h *cartHandler) FindAll(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req cart.FindAllRequest
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	res, err := h.cartService.FindAll(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}

func (h *cartHandler) DeleteProductFromCarts(c echo.Context) (err error) {
	ctx := c.Request().Context()

	uuid := c.Param("uuid")
	if err = utils.ValidateUUID(ctx, uuid); err != nil {
		return
	}

	res, err := h.cartService.DeleteProduct(ctx, uuid)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}
