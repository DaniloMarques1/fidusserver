package services

import (
	"errors"

	"github.com/danilomarques1/fidusserver/models"
)

type VerifyMasterPasswordExpirationService interface {
	Execute(masterId string) error
}

type verifyMasterPasswordExpirationService struct {
	masterDAO models.MasterDAO
}

func NewVerifyMasterPasswordExpirationService() VerifyMasterPasswordExpirationService {
	masterDAO := models.NewMasterDAODatabase()
	return &verifyMasterPasswordExpirationService{masterDAO}
}

func (v *verifyMasterPasswordExpirationService) Execute(masterId string) error {
	master, err := v.masterDAO.FindById(masterId)
	if err != nil {
		return err
	}
	if master.IsPasswordExpired() {
		// TODO better error
		return errors.New("Master password has expired and it needs to be reseted")
	}
	return nil
}
