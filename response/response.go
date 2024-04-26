package response

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/fidusserver/apierror"
)

func Success(w http.ResponseWriter, body any, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

type ErrorResponseDto struct {
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, err error) {
	switch v := err.(type) {
	case apierror.ApiError:
		Success(w, ErrorResponseDto{Message: v.Error()}, v.Status())
	default:
		Success(w, ErrorResponseDto{Message: err.Error()}, http.StatusInternalServerError)
	}
}
