package services

import (
	"github.com/danilomarques1/fidusserver/apierror"
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
		if v.masterDAO.NoMatchError(err) {
			return apierror.ErrMasterNotFound()
		}
		return err
	}
	if master.IsPasswordExpired() {
		return apierror.ErrPasswordExpired()
	}
	return nil
}
