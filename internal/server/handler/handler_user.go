package handler

import (
	"net/http"

	"github.com/armiariyan/synapsis/internal/pkg/utils"
	"github.com/armiariyan/synapsis/internal/usecase/user"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) SetUserService(service user.Service) *userHandler {
	h.userService = service
	return h
}

func (h *userHandler) Validate() *userHandler {
	if h.userService == nil {
		panic("userService is nil")
	}
	return h
}

func (h *userHandler) Register(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req user.RegisterRequest
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	res, err := h.userService.Register(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}

func (h *userHandler) Login(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req user.LoginRequest
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	res, err := h.userService.Login(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}

func (h *userHandler) Checkout(c echo.Context) (err error) {
	ctx := c.Request().Context()

	res, err := h.userService.Checkout(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}
