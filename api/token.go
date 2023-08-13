package api

import (
	"encoding/json"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	db "github.com/zcubbs/tlz/db/sqlc"
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) renewAccessToken(c *fiber.Ctx) error {
	var req renewAccessTokenRequest
	if c.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   msg,
		})
	}

	// Decode the request body into struct and failed if any error occur
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		msg := "Request body contains badly-formed JSON"
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	// Validate the request body
	err = s.validate.validator.Struct(req)
	if err != nil {
		log.Error("failed to validate request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid request body",
		})
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		log.Error("failed to verify refresh token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid refresh token",
		})
	}

	session, err := s.store.GetSession(c.Context(), refreshPayload.ID)
	if err != nil {
		log.Error("failed to get session", "error", err)
		if errors.Is(err, db.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "invalid refresh token",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "could not get session",
		})
	}

	if session.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid refresh token",
		})
	}

	if session.ID != refreshPayload.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid refresh token",
		})
	}

	if session.RefreshToken != req.RefreshToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid refresh token",
		})
	}

	if session.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid refresh token",
		})
	}

	user, err := s.store.GetUser(c.Context(), session.UserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "could not get user",
		})
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		user.ID,
		s.cfg.AccessTokenDuration,
	)
	if err != nil {
		log.Error("failed to create access token", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "could not create access token",
		})
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return c.Status(fiber.StatusOK).JSON(&rsp)
}
