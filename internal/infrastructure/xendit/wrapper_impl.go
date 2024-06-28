package xendit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/armiariyan/synapsis/internal/config"
	"github.com/armiariyan/synapsis/internal/pkg/rest"
)

type wrapper struct {
	client   rest.RestClient
	username string
	password string
}

func NewWrapper() *wrapper {
	return &wrapper{}
}

func (w *wrapper) Setup() Wrapper {
	restOptions := rest.Options{
		TagName: "xendit",
		Address: config.GetString("xendit.host"),
		Timeout: config.GetDuration("xendit.timeout"),
		SkipTLS: config.GetBool("xendit.skiptls"),
	}

	username := config.GetString("xendit.username")
	password := config.GetString("xendit.password")

	w.username = username
	w.password = password

	w.client = rest.New(restOptions)

	return w
}

func (w *wrapper) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (resp CreateInvoiceResponse, err error) {
	path := config.GetString("xendit.path.createInvoice")

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Basic %s", BasicAuth(w.username, w.password)))
	header.Add("Content-Type", "application/json")

	res, status, err := w.client.Post(ctx, path, header, req, false)
	if err != nil {
		return
	}

	if status != http.StatusOK {
		var errData Error
		json.Unmarshal(res, &errData)
		if errData.Message == "" {
			errData.Message = fmt.Sprintf("xendit servcer return http status %d", status)
		}
		err = fmt.Errorf(errData.Message)
		return
	}

	err = json.Unmarshal(res, &resp)

	return
}
