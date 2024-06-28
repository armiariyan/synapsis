package products

import (
	"context"

	"github.com/armiariyan/synapsis/internal/pkg/constants"
)

type Service interface {
	FindByCategory(ctx context.Context, req FindByCategoryRequest) (res constants.DefaultResponse, err error)
	AddToCart(ctx context.Context, req AddToCartRequest) (res constants.DefaultResponse, err error)
}
