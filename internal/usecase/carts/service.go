package carts

import (
	"context"

	"github.com/armiariyan/synapsis/internal/pkg/constants"
)

type Service interface {
	FindAll(ctx context.Context, req FindAllRequest) (res constants.DefaultResponse, err error)
	DeleteProduct(ctx context.Context, uuid string) (res constants.DefaultResponse, err error)
}
