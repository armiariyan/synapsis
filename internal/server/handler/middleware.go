package handler

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/armiariyan/synapsis/internal/config"
	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/labstack/echo/v4"
)

func (h *Handler) BasicAuth(prefix string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			basicToken := c.Request().Header.Get("Authorization")
			if basicToken == "" {
				c.Set("unauthorized", true)
				err := fmt.Errorf("authorization header is empty")
				log.Error(ctx, "authentication failed", err)
				return err
			}
			sliceToken := strings.Split(basicToken, "Basic ")
			if len(sliceToken) < 2 {
				c.Set("unauthorized", true)
				err := fmt.Errorf("basic token slice is not greater than 2")
				log.Error(ctx, "authentication failed", err)
				return err
			}
			token := sliceToken[1]
			username := config.GetString(fmt.Sprintf("basicAuth.%s.username", prefix))
			password := config.GetString(fmt.Sprintf("basicAuth.%s.password", prefix))
			if authToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))); authToken != token {
				c.Set("forbidden", true)
				err := fmt.Errorf("invalid basic auth provided")
				log.Error(ctx, "authentication failed", prefix, err, token, authToken)
				return err
			}
			log.Info(ctx, "authentication success", "basic auth passed", prefix)
			return next(c)
		}
	}
}
