package services

import (
	"database/sql"
	"errors"

	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/models"
	"github.com/danilomarques1/fidusserver/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateMasterService interface {
	Execute(*dtos.AuthenticateRequestDto) (string, error)
}

type authenticateMasterService struct {
	masterDAO models.MasterDAO
}

func NewAuthenticateMasterService() AuthenticateMasterService {
	masterDAO := models.NewMasterDAODatabase()
	return &authenticateMasterService{masterDAO}
}

func (authService *authenticateMasterService) Execute(req *dtos.AuthenticateRequestDto) (string, error) {
	master, err := authService.masterDAO.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apierror.MasterEmailNotFound()
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(master.PasswordHash), []byte(req.Password)); err != nil {
		return "", apierror.MasterIncorrectPassword()
	}

	token, err := token.GenerateToken(master.ID, master.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
