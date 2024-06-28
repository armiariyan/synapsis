package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/armiariyan/synapsis/internal/config"
	"github.com/armiariyan/synapsis/internal/domain/repositories"
	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/armiariyan/synapsis/internal/pkg/utils"
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

func (h *Handler) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		header := c.Request().Header

		bearerToken := header.Get("Authorization")
		if bearerToken == "" {
			c.Set("unauthorized", true)
			err := fmt.Errorf("unauthorized [0]")
			log.Error(ctx, "authentication failed", err)
			return err
		}

		sliceToken := strings.Split(bearerToken, "Bearer ")
		if len(sliceToken) < 2 {
			c.Set("unauthorized", true)
			err := fmt.Errorf("unauthorized [1]")
			log.Error(ctx, "authentication failed", err)
			return err
		}

		token := sliceToken[1]
		cl, err := utils.JwtVerify(token)
		if err != nil {
			c.Set("unauthorized", true)
			log.Error(ctx, "authentication failed", "failed to verify token", err)
			err = fmt.Errorf("unauthorized [2]")
			return err
		}

		claims := cl.Data.(utils.JWTClaimsData)

		userData, err := repositories.NewUser(h.synapsisDB).FindByUUID(ctx, claims.ID)
		if err != nil {
			c.Set("unauthorized", true)
			log.Error(ctx, "authentication failed", "failed to find user by id", err, claims.ID)
			err = fmt.Errorf("unauthorized [2]")
			return err
		}

		if token != userData.Token {
			c.Set("unauthorized", true)
			err = fmt.Errorf("token mismatch")
			log.Error(ctx, "authentication failed", err, claims.ID)
			return err
		}

		ctx = context.WithValue(ctx, string("user.id"), claims.ID)
		ctx = context.WithValue(ctx, string("user"), userData)

		c.SetRequest(c.Request().WithContext(ctx))
		log.Info(ctx, "authentication success", "authentication passed", userData)
		return next(c)
	}
}
