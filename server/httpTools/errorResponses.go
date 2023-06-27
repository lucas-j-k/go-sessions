package httpTools

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse implements chi Render.Renderer to return a standard format for error responses
type ErrResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// 400 - Bad Request
func ErrInvalidRequest() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 400,
		ErrorText:      "invalid_request",
	}
}

// 500 - Internal Server Error
func ErrInternalServer() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 500,
		ErrorText:      "internal_server_error",
	}
}

// 404 - Not found
func ErrNotFound() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 404,
		ErrorText:      "not_found",
	}
}

// 401 Unauthorized
func ErrUnauthorized() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 401,
		ErrorText:      "unauthorized",
	}
}
