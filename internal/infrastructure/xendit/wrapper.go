package xendit

import (
	"context"
)

type Wrapper interface {
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (resp CreateInvoiceResponse, err error)
}
