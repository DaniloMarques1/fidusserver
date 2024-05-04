package response

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
)

func Json(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func Error(w http.ResponseWriter, err error) {
	switch v := err.(type) {
	case apierror.ApiError:
		Json(w, v.Status(), dtos.ErrorResponseDto{Message: v.Error()})
	default:
		Json(w, http.StatusInternalServerError, dtos.ErrorResponseDto{Message: err.Error()})
	}
}
