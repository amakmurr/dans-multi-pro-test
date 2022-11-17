package internal

import (
	"errors"
	"net/http"

	v1 "github.com/amakmurr/dans-multi-pro-test/pkg/openapi"

	"github.com/go-chi/render"
)

var (
	ErrDataNotFound   = errors.New("data not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrValidation     = errors.New("validation error")
	ErrInternalSystem = errors.New("internal system error")
)

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	resp := &v1.DefaultErrorResponse{}
	httpCode := http.StatusInternalServerError
	resp.Message = err.Error()

	switch {
	case errors.Is(err, ErrDataNotFound):
		httpCode = http.StatusNotFound
	case errors.Is(err, ErrUnauthorized):
		httpCode = http.StatusUnauthorized
	case errors.Is(err, ErrValidation):
		httpCode = http.StatusBadRequest
	case errors.Is(err, ErrInternalSystem):
		httpCode = http.StatusInternalServerError
	}

	render.Status(r, httpCode)
	render.JSON(w, r, resp)
}
