package services

import "github.com/danilomarques1/fidusserver/models"

type RetrievePasswordService interface {
	Execute(string, string) (*models.Password, error)
}

type retrievePasswordService struct {
	passwordDAO models.PasswordDAO
}

func NewRetrievePasswordService() RetrievePasswordService {
	passwordDAO := models.NewPasswordDAODatabase()
	return &retrievePasswordService{passwordDAO}
}

func (retrieveService *retrievePasswordService) Execute(masterId, key string) (*models.Password, error) {
	password, err := retrieveService.passwordDAO.FindOne(masterId, key)
	if err != nil {
		return nil, err
	}
	return password, nil
}
