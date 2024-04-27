package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/response"
	"github.com/danilomarques1/fidusserver/services"
)

func StorePassword(w http.ResponseWriter, r *http.Request) {
	masterId := r.Context().Value("masterId").(string) // TODO: protect it
	body := &dtos.StorePasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("Error parsing json body %v\n", err)
		response.Error(w, apierror.InternalServerError(err))
		return
	}
	storePasswordService := services.NewStorePasswordService()
	// TODO: get id from the token
	if err := storePasswordService.Execute(masterId, body); err != nil {
		log.Printf("Error storing the password %v\n", err)
		response.Error(w, err)
		return
	}

	response.Json(w, http.StatusNoContent, nil)
}
