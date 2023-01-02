package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func (res ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Code)
	return nil
}

func NewErrorResponse(message string, code int) render.Renderer {
	return &ErrorResponse{Message: message, Code: code}
}
