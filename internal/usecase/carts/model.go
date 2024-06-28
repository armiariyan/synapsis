package carts

import "github.com/armiariyan/synapsis/internal/pkg/constants"

// * Requests
type (
	FindAllRequest struct {
		constants.PaginationRequest
	}

	DeleteItemRequest struct {
		Category string `query:"category" validate:"required,oneof=electronics books"`
		constants.PaginationRequest
	}
)

// * Responses
type (
	FindAllResponse struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Amount   float64 `json:"amount"`
		Category string  `json:"category"`
	}
)
