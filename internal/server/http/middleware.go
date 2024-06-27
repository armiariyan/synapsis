package http

import (
	"net/http"
	"strings"

	"github.com/armiariyan/bepkg/response"
	"github.com/armiariyan/bepkg/utils"
	"github.com/armiariyan/logger"
	"github.com/armiariyan/synapsis/internal/config"
	"github.com/armiariyan/synapsis/internal/infrastructure/container"
	"github.com/armiariyan/synapsis/internal/pkg/constants"
	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddleware(server *echo.Echo, container *container.Container) {
	server.Use(SetLoggerMiddleware())
	server.Use(LoggerMiddleware())

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderXCSRFToken},
		ExposeHeaders:    []string{echo.HeaderContentLength, echo.HeaderAccessControlAllowOrigin, echo.HeaderContentDisposition},
		AllowCredentials: true,
	}))

	server.HTTPErrorHandler = errorHandler
	v := validator.New()
	// * custom register validation here
	// v.RegisterValidation("ISO8601Date", utility.IsISO8601Date)
	server.Validator = &DataValidator{ValidatorData: v}
}

func SetLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip Logging
			if isLoggingSkip(c) {
				return next(c)
			}

			threadId := c.Request().Header.Get(constants.HEADER_XID)
			if len(threadId) == 0 {
				threadId = utils.GenerateThreadId()
			}

			ctxLogger := logger.Context{
				ServiceName:    config.GetString("app.name"),
				ServiceVersion: config.GetString("app.version"),
				ServicePort:    config.GetInt("app.port"),
				ThreadID:       threadId,
				JourneyID:      c.Request().Header.Get(constants.HEADER_JID),
				Tag:            config.GetString("app.name"),
				ReqMethod:      c.Request().Method,
				ReqURI:         c.Request().URL.String(),
			}

			request := c.Request()

			ctx := log.SetContextFromEchoRequest(c)
			ctx = logger.InjectCtx(ctx, ctxLogger)
			c.SetRequest(request.WithContext(ctx))

			return next(c)
		}
	}
}

func LoggerMiddleware() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		ctx := c.Request().Context()
		if isLoggingRequestOnly(c) {
			log.TDR(ctx, reqBody, []byte{})
			return
		}
		// Skip Logging
		if isLoggingSkip(c) {
			return
		}

		if c.Get("skip-body-logging") != nil {
			log.TDR(ctx, []byte{}, resBody)
			return
		}

		log.TDR(ctx, reqBody, resBody)
	})
}

func isLoggingSkip(c echo.Context) bool {
	requestPath := c.Request().URL.String()
	skipPath := map[string]bool{
		"/": true,
	}

	return skipPath[requestPath]
}

func isLoggingRequestOnly(c echo.Context) bool {
	requestPath := c.Request().URL.String()
	return strings.Contains(requestPath, "/download")
}

func errorHandler(err error, c echo.Context) {
	// Need this, because somehow if default error handler use with echo body dump
	// It will be print response error twice
	if c.Get("error-handled") != nil {
		return
	}

	c.Set("error-handled", true)

	resp := constants.DefaultResponse{
		Status:  response.GeneralError,
		Message: err.Error(),
		Data:    struct{}{},
		Errors:  make([]string, 0),
	}

	if c.Get("invalid-format") != nil || strings.Contains(err.Error(), "Error:Field validation for") {
		resp.Status = constants.STATUS_BAD_REQUEST
		resp.Message = constants.MESSAGE_BAD_REQUEST
		resp.Errors = append(resp.Errors, strings.Split(err.Error(), "\n")...)
	}
	if !isLoggingSkip(c) {
		request := c.Request()

		ctx := log.SetErrorMessageFromEchoContext(c, err.Error())
		c.SetRequest(request.WithContext(ctx))

		log.Error(ctx, err.Error())
	}

	c.JSON(http.StatusOK, resp)
}

type DataValidator struct {
	ValidatorData *validator.Validate
}

func (cv *DataValidator) Validate(i interface{}) error {
	return cv.ValidatorData.Struct(i)
}
