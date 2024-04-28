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
	masterId, ok := r.Context().Value("masterId").(string)
	if !ok {
		response.Json(w, http.StatusForbidden, nil)
		return
	}
	body := &dtos.StorePasswordRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("Error parsing json body %v\n", err)
		response.Error(w, apierror.InternalServerError(err))
		return
	}
	storePasswordService := services.NewStorePasswordService()
	if err := storePasswordService.Execute(masterId, body); err != nil {
		log.Printf("Error storing the password %v\n", err)
		response.Error(w, err)
		return
	}

	response.Json(w, http.StatusNoContent, nil)
}

func RetrievePassword(w http.ResponseWriter, r *http.Request) {
	masterId, ok := r.Context().Value("masterId").(string)
	if !ok {
		response.Json(w, http.StatusForbidden, nil)
		return
	}
	key := r.URL.Query().Get("key")
	if len(key) == 0 {
		log.Println("Key is empty")
		response.Error(w, apierror.InvalidKey())
		return
	}
	retrievePassword := services.NewRetrievePasswordService()
	password, err := retrievePassword.Execute(masterId, key)
	if err != nil {
		log.Printf("Error retrieving password %v\n", err)
		response.Error(w, err)
		return
	}

	resp := &dtos.RetrievePasswordResponseDto{
		Key:      password.Key,
		Password: password.PasswordValue,
		MasterId: password.MasterId,
	}
	response.Json(w, http.StatusOK, resp)
}

func DeletePassword(w http.ResponseWriter, r *http.Request) {
	masterId, ok := r.Context().Value("masterId").(string)
	if !ok {
		log.Println("Master id is empty")
		response.Json(w, http.StatusForbidden, nil)
		return
	}
	key := r.URL.Query().Get("key")
	if len(key) == 0 {
		log.Println("Key is empty")
		response.Error(w, apierror.InvalidKey())
		return
	}

	deletePassword := services.NewDeletePasswordService()
	if err := deletePassword.Execute(masterId, key); err != nil {
		log.Printf("Delete Error: %v\n", err)
		response.Error(w, err)
		return
	}
	response.Json(w, http.StatusNoContent, nil)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	masterId, ok := r.Context().Value("masterId").(string)
	if !ok {
		response.Json(w, http.StatusForbidden, nil)
		return
	}
	key := r.URL.Query().Get("key")
	if len(key) == 0 {
		log.Println("Key is empty")
		response.Error(w, apierror.InvalidKey())
		return
	}

	body := &dtos.UpdatePasswordRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("Error parsing json body %v\n", err)
		response.Error(w, apierror.InvalidRequestBody(err.Error()))
		return
	}

	updatePassword := services.NewUpdatePasswordService()
	if err := updatePassword.Execute(masterId, key, body.Password); err != nil {
		log.Printf("Update Error: %v\n", err)
		response.Error(w, err)
		return
	}

	response.Json(w, http.StatusNoContent, nil)
}
