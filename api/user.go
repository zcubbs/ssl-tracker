package api

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createUser(c *fiber.Ctx) error {

	return nil
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *Server) loginUser(c *fiber.Ctx) error {

	return nil
}
