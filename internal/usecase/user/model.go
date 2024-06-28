package user

// * Requests
type (
	RegisterRequest struct {
		Name                 string `json:"name" validate:"required"`
		Email                string `json:"email" validate:"required,email"`
		PhoneNumber          string `json:"phoneNumber" validate:"required,number"`
		Password             string `json:"password" validate:"required"`
		PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
)

// * Responses
type (
	LoginResponse struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		AccessCode  string `json:"token"`
		ExpiredAt   int64  `json:"expiredAt"`
	}

	CheckoutResponse struct {
		InvoiceID  string `json:"invoiceId"`
		InvoiceURL string `json:"invoiceUrl"`
	}
)
