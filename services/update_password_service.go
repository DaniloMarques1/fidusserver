package services

import "github.com/danilomarques1/fidusserver/models"

type UpdatePasswordService interface {
	Execute(masterId, key, passwordValue string) error
}

type updatePasswordService struct {
	passwordDAO models.PasswordDAO
}

func NewUpdatePasswordService() UpdatePasswordService {
	passwordDAO := models.NewPasswordDAODatabase()
	return &updatePasswordService{passwordDAO}
}

func (service *updatePasswordService) Execute(masterId, key, passwordValue string) error {
	if err := service.passwordDAO.UpdatePasswordValue(masterId, key, passwordValue); err != nil {
		return err
	}
	return nil
}
