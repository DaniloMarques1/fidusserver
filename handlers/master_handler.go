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

func CreateMaster(w http.ResponseWriter, r *http.Request) {
	body := &dtos.CreateMasterRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("Invalid request %v\n", err)
		response.Error(w, apierror.InvalidRequestBody(err.Error()))
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

	response.Success(w, respBody, http.StatusCreated)
}

func AuthenticateMaster(w http.ResponseWriter, r *http.Request) {
	body := &dtos.AuthenticateRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("Invalid request %v\n", err)
		response.Error(w, apierror.InvalidRequestBody(err.Error()))
		return
	}

	authenticateservice := services.NewAuthenticateMasterService()
	accessToken, err := authenticateservice.Execute(body)
	if err != nil {
		log.Printf("Error %v\n", err)
		response.Error(w, err)
		return
	}

	resp := &dtos.AuthenticateResponseDto{AccessToken: accessToken}

	response.Success(w, resp, http.StatusOK)
}
