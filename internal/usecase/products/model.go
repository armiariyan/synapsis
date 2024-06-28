package products

import "github.com/armiariyan/synapsis/internal/pkg/constants"

// * Requests
type (
	FindByCategoryRequest struct {
		Category string `query:"category" validate:"required,oneof=electronics books"`
		constants.PaginationRequest
	}

	AddToCartRequest struct {
		ProductID string  `json:"productId" validate:"required"`
		Amount    float64 `json:"amount" validate:"required,number"` // * amount of product
	}
)

// * Responses
type (
	FindByCategoryResponse struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		Price        float64 `json:"price"`
		CategoryName string  `json:"categoryName"`
	}
)
