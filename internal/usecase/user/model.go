package user

import "github.com/armiariyan/synapsis/internal/pkg/constants"

// * Requests
type (
	RegisterRequest struct {
		Name                 string `json:"name" validate:"required"`
		Email                string `json:"email" validate:"required,email"`
		PhoneNumber          string `json:"phoneNumber" validate:"required,number"`
		Password             string `json:"password" validate:"required"`
		PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
	}

	FindAllRequest struct {
	}
)

// * Responses
type (
	CreateResponse struct {
		constants.DefaultResponse
	}

	FindAllResponse struct {
		constants.DefaultResponse
	}
)
