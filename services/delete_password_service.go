package services

import (
	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/models"
)

type DeletePasswordService interface {
	Execute(masterId, key string) error
}

type deletePasswordService struct {
	passwordDAO models.PasswordDAO
}

func NewDeletePasswordService() DeletePasswordService {
	passwordDAO := models.NewPasswordDAODatabase()
	return &deletePasswordService{passwordDAO}
}

func (service *deletePasswordService) Execute(masterId, key string) error {
	if _, err := service.passwordDAO.FindOne(masterId, key); err != nil {
		if service.passwordDAO.NoMatchError(err) {
			return apierror.ErrPasswordNotFound()
		}
		return err
	}

	if err := service.passwordDAO.Delete(masterId, key); err != nil {
		return err
	}
	return nil
}
