package services

import (
	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/models"
)

type StorePasswordService interface {
	Execute(string, *dtos.StorePasswordRequestDto) error
}

type storePasswordService struct {
	masterDAO   models.MasterDAO
	passwordDAO models.PasswordDAO
}

func NewStorePasswordService() StorePasswordService {
	masterDAO := models.NewMasterDAODatabase()
	passwordDAO := models.NewPasswordDAODatabase()
	return &storePasswordService{masterDAO, passwordDAO}
}

func (passwordService *storePasswordService) Execute(masterId string, req *dtos.StorePasswordRequestDto) error {
	pwdKeyUsed, err := passwordService.passwordDAO.FindOne(masterId, req.Key)
	if err != nil && !passwordService.passwordDAO.NoMatchError(err) {
		return err
	}
	if pwdKeyUsed != nil {
		return apierror.ErrKeyAlreadyUsed()
	}

	if _, err := passwordService.masterDAO.FindById(masterId); err != nil {
		if passwordService.masterDAO.NoMatchError(err) {
			return apierror.ErrMasterNotFound()
		}
		return err
	}
	password := &models.Password{Key: req.Key, MasterId: masterId, PasswordValue: req.Password}
	if err := passwordService.passwordDAO.Save(password); err != nil {
		return err
	}

	return nil
}
