package services

import (
	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/models"
)

type ResetMasterPasswordService interface {
	Execute(body *dtos.ResetMasterPasswordRequestDto) error
}

type resetMasterPasswordService struct {
	masterDAO models.MasterDAO
}

func NewResetMasterPasswordService() ResetMasterPasswordService {
	masterDAO := models.NewMasterDAODatabase()
	return &resetMasterPasswordService{masterDAO}
}

func (r *resetMasterPasswordService) Execute(body *dtos.ResetMasterPasswordRequestDto) error {
	master, err := r.masterDAO.FindByEmail(body.Email)
	if err != nil {
		if r.masterDAO.NoMatchError(err) {
			return apierror.ErrMasterNotFound()
		}
		return err
	}
	if err := master.ComparePassword(body.OldPassword); err != nil {
		return apierror.ErrIncorrectCredentials()
	}
	if err := master.HashPassword(body.NewPassword); err != nil {
		return err
	}
	master.SetPasswordExpiration()

	if err := r.masterDAO.ResetPassword(master.ID, master.PasswordHash, master.PasswordExpirationDate); err != nil {
		return err
	}
	return nil
}
