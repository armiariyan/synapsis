package user

import (
	"context"

	"github.com/armiariyan/synapsis/internal/pkg/constants"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (res constants.DefaultResponse, err error)
	// FindAll(ctx context.Context, req FindAllRequest) (res FindAllResponse, err error)
}
