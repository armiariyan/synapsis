package xendit

import "time"

// * Request
type (
	CreateInvoiceRequest struct {
		ExternalID string        `json:"external_id"`
		Amount     int           `json:"amount"` // * total price amount
		Items      []ItemInvoice `json:"items"`
	}

	ItemInvoice struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
		Category string `json:"category"`
		URL      string `json:"url"`
	}

	Error struct {
		ErrorCode string `json:"error_code"`
		Message   string `json:"message"`
	}
)

// * Response
type (
	CreateInvoiceResponse struct {
		ID         string        `json:"id"`
		ExternalID string        `json:"external_id"`
		InvoiceURL string        `json:"invoice_url"`
		UserID     string        `json:"user_id"`
		Status     string        `json:"status"`
		Amount     int           `json:"amount"`
		ExpiryDate time.Time     `json:"expiry_date"`
		Created    time.Time     `json:"created"`
		Updated    time.Time     `json:"updated"`
		Currency   string        `json:"currency"`
		Items      []ItemInvoice `json:"items"`
		Metadata   interface{}   `json:"metadata"`
	}
)
