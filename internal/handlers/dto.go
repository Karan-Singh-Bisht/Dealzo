//DTO fixes mass assignment problem

package handlers

import (
	"fmt"
	"strings"
	"time"
)

type CreateListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	City        string `json:"city"`
}

type CreateListingResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type ValidationError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

func (req CreateListingRequest) Validate() error {
	if strings.TrimSpace(req.Title) == "" {
		return &ValidationError{Field: "title", Msg: "field is required"}
	}
	if req.Price == 0 {
		return &ValidationError{Field: "price", Msg: "price should be greater than 0"}
	}
	if strings.TrimSpace(req.City) == "" {
		return &ValidationError{Field: "city", Msg: "city is required"}
	}
	if strings.TrimSpace(req.Description) == "" {
		return &ValidationError{Field: "description", Msg: "description is required"}
	}
	return nil
}
