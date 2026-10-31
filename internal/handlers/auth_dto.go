package handlers

import (
	"net/mail"
	"strings"
	"time"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func (req CreateUserRequest) AuthValidate() error {
	if strings.TrimSpace(req.Name) == "" {
		return &ValidationError{Field: "username", Msg: "Username is required"}
	}
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return &ValidationError{Field: "email", Msg: "Must be a valid email"}
	}
	if len(req.Password) < 8 {
		return &ValidationError{Field: "password", Msg: "Password must be atleast 8 characters"}
	}
	return nil
}
