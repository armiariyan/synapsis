package rest

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/armiariyan/logger"
	"github.com/armiariyan/rest"
	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/opentracing/opentracing-go"
)

type RestClient interface {
	Post(ctx context.Context, path string, header http.Header, request interface{}, isXML bool) (body []byte, statusCode int, err error)
	Patch(ctx context.Context, path string, header http.Header, request interface{}, isXML bool) (body []byte, statusCode int, err error)
	Get(ctx context.Context, path string, header http.Header) (body []byte, statusCode int, err error)
}

type client struct {
	options    Options
	httpClient *rest.DefaultHttpRequester
}

type restLogger struct {
	tag string
}

func New(options Options) RestClient {
	httpClient := &http.Client{
		Timeout: options.Timeout,
	}

	if options.SkipTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpClient.Transport = tr
	}

	requestLogger := &restLogger{
		tag: options.TagName,
	}

	restDefaultClient, _ := rest.DefaultClient(httpClient, rest.AddHook(requestLogger))

	return &client{
		options:    options,
		httpClient: restDefaultClient,
	}
}

func (r restLogger) BeforeRequest(ctx context.Context, data rest.HookData) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "restLogger.BeforeRequest")
	defer func() {
		span.Finish()
		ctx.Done()
	}()

	method := data.Request.Method

	log.T2(ctx, fmt.Sprintf("%s [Request]", method), data.URL, data.Request.Header, data.Request.Body)

}

func (r restLogger) AfterRequest(ctx context.Context, data rest.HookData) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "restLogger.AfterRequest")
	defer func() {
		span.Finish()
		ctx.Done()
	}()

	method := data.Request.Method
	log.T3(ctx, fmt.Sprintf("%s [Response]", method), data.StartTime, data.URL, data.Response.Body)
}

func (c *client) Post(ctx context.Context, path string, header http.Header, request interface{}, isXML bool) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	ctxLogger := logger.ExtractCtx(ctx)

	reqByte, _ := json.Marshal(request)
	if request == nil {
		reqByte = nil
	}
	if isXML {
		reqByte = []byte(fmt.Sprintf("%v", request))
	}
	resp, err := c.httpClient.Post(ctx, ctxLogger.ThreadID, url, header, reqByte)
	if err != nil {
		log.Error(ctx, "error post request", err.Error())
		return
	}

	body = resp.RespBody
	statusCode = resp.Raw.StatusCode

	return
}

func (c *client) Patch(ctx context.Context, path string, header http.Header, request interface{}, isXML bool) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	ctxLogger := logger.ExtractCtx(ctx)

	reqByte, _ := json.Marshal(request)
	if request == nil {
		reqByte = nil
	}
	if isXML {
		reqByte = []byte(fmt.Sprintf("%v", request))
	}
	resp, err := c.httpClient.Patch(ctx, ctxLogger.ThreadID, url, header, reqByte)
	if err != nil {
		log.Error(ctx, "error patch request", err.Error())
		return
	}

	body = resp.RespBody
	statusCode = resp.Raw.StatusCode

	return
}

func (c *client) Get(ctx context.Context, path string, header http.Header) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	ctxLogger := logger.ExtractCtx(ctx)

	resp, err := c.httpClient.Get(ctx, ctxLogger.ThreadID, url, header)
	if err != nil {
		log.Error(ctx, "error get request", err.Error())
		return
	}

	body = resp.RespBody
	statusCode = resp.Raw.StatusCode

	return
}
