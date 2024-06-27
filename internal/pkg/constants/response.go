package constants

const (
	STATUS_SUCCESS      = "00"
	STATUS_BAD_REQUEST  = "400"
	STATUS_UNAUTHORIZED = "401"
	STATUS_FORBIDDEN    = "403"
	STATUS_CONFLICT     = "409"
)

const (
	MESSAGE_SUCCESS     = "success"
	MESSAGE_BAD_REQUEST = "invalid request format"
	MESSAGE_FAILED      = "something went wrong"
)

type DefaultResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

type PaginationData struct {
	Page        uint `json:"page"`
	TotalPages  uint `json:"totalPages"`
	TotalItems  uint `json:"totalItems"`
	Limit       uint `json:"limit"`
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
}

type PaginationResponseData struct {
	Results        interface{} `json:"results"`
	PaginationData `json:"pagination"`
}
