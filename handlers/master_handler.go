package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/response"
	"github.com/danilomarques1/fidusserver/services"
	"github.com/danilomarques1/fidusserver/validate"
)

func CreateMaster(w http.ResponseWriter, r *http.Request) {
	body := &dtos.CreateMasterRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("Invalid request %v\n", err)
		response.Error(w, apierror.ErrInvalidRequest(err.Error()))
		return
	}

	validate := validate.Validate()
	if err := validate.Struct(body); err != nil {
		response.Error(w, apierror.ErrInvalidRequest("Invalid parameters"))
		return
	}

	registerService := services.NewRegisterService()
	master, err := registerService.Execute(body)
	if err != nil {
		log.Printf("Error %v\n", err)
		response.Error(w, err)
		return
	}

	respBody := &dtos.CreateMasterResponseDto{
		ID:           master.ID,
		Name:         master.Name,
		Email:        master.Email,
		PasswordHash: master.PasswordHash,
	}

	response.Json(w, http.StatusCreated, respBody)
}

func AuthenticateMaster(w http.ResponseWriter, r *http.Request) {
	body := &dtos.AuthenticateRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("Invalid request %v\n", err)
		response.Error(w, apierror.ErrInvalidRequest(err.Error()))
		return
	}
	validate := validate.Validate()
	if err := validate.Struct(body); err != nil {
		response.Error(w, apierror.ErrInvalidRequest("Invalid parameters"))
		return
	}

	authenticateservice := services.NewAuthenticateMasterService()
	accessToken, expiresAt, err := authenticateservice.Execute(body)
	if err != nil {
		log.Printf("Error %v\n", err)
		response.Error(w, err)
		return
	}

	resp := &dtos.AuthenticateResponseDto{AccessToken: accessToken, ExpiresAt: expiresAt}

	response.Json(w, http.StatusOK, resp)
}
