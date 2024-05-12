package services

import (
	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/models"
)

type RetrieveKeys interface {
	Execute(masterId string) ([]string, error)
}

type retrieveKeys struct {
	passwordDAO models.PasswordDAO
	masterDAO   models.MasterDAO
}

func NewRetrieveKeys() RetrieveKeys {
	passwordDAO := models.NewPasswordDAODatabase()
	masterDAO := models.NewMasterDAODatabase()
	return &retrieveKeys{passwordDAO, masterDAO}
}

func (rk *retrieveKeys) Execute(masterId string) ([]string, error) {
	if _, err := rk.masterDAO.FindById(masterId); err != nil {
		if rk.masterDAO.NoMatchError(err) {
			return nil, apierror.ErrMasterNotFound()
		}
		return nil, err
	}

	keys, err := rk.passwordDAO.Keys(masterId)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
