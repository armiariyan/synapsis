package handler

import (
	"net/http"

	"github.com/armiariyan/synapsis/internal/pkg/utils"
	product "github.com/armiariyan/synapsis/internal/usecase/products"
	"github.com/labstack/echo/v4"
)

type productHandler struct {
	productService product.Service
}

func NewProductHandler() *productHandler {
	return &productHandler{}
}

func (h *productHandler) SetProductService(service product.Service) *productHandler {
	h.productService = service
	return h
}

func (h *productHandler) Validate() *productHandler {
	if h.productService == nil {
		panic("productService is nil")
	}
	return h
}

func (h *productHandler) FindAllProductsByCategory(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req product.FindByCategoryRequest
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	res, err := h.productService.FindByCategory(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}

func (h *productHandler) AddToCart(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req product.AddToCartRequest
	if err = utils.Validate(c, &req); err != nil {
		return
	}

	res, err := h.productService.AddToCart(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}
