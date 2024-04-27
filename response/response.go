package response

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/fidusserver/apierror"
)

func Json(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

type ErrorResponseDto struct {
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, err error) {
	switch v := err.(type) {
	case apierror.ApiError:
		Json(w, v.Status(), ErrorResponseDto{Message: v.Error()})
	default:
		Json(w, http.StatusInternalServerError, ErrorResponseDto{Message: err.Error()})
	}
}
