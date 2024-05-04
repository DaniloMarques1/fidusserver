package services

import (
	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/models"
	"github.com/danilomarques1/fidusserver/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateMasterService interface {
	Execute(*dtos.AuthenticateRequestDto) (string, int64, error)
}

type authenticateMasterService struct {
	masterDAO models.MasterDAO
}

func NewAuthenticateMasterService() AuthenticateMasterService {
	masterDAO := models.NewMasterDAODatabase()
	return &authenticateMasterService{masterDAO}
}

func (authService *authenticateMasterService) Execute(req *dtos.AuthenticateRequestDto) (string, int64, error) {
	master, err := authService.masterDAO.FindByEmail(req.Email)
	if err != nil {
		if authService.masterDAO.NoMatchError(err) {
			return "", 0, apierror.ErrIncorrectCredentials()
		}
		return "", 0, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(master.PasswordHash), []byte(req.Password)); err != nil {
		return "", 0, apierror.ErrIncorrectCredentials()
	}

	signedToken, expiresAt, err := token.GenerateToken(master.ID, master.Email)
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresAt, nil
}
