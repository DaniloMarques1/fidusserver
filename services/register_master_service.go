package services

import (
	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/models"
)

type RegisterService interface {
	Execute(*dtos.CreateMasterRequestDto) (*models.Master, error)
}

type registerService struct {
	dao models.MasterDAO
}

func NewRegisterService() RegisterService {
	dao := models.NewMasterDAODatabase()
	return &registerService{dao}
}

func (service *registerService) Execute(createMasterDto *dtos.CreateMasterRequestDto) (*models.Master, error) {
	m, err := service.dao.FindByEmail(createMasterDto.Email)
	if !service.dao.NoMatchError(err) {
		if m != nil {
			return nil, apierror.ErrEmailAlreadyTaken()
		}
		return nil, err
	}

	master, err := models.NewMaster(createMasterDto.Name, createMasterDto.Email, createMasterDto.Password)
	if err != nil {
		return nil, err
	}

	if err := service.dao.Save(master); err != nil {
		return nil, apierror.ErrInternalServerError(err.Error())
	}

	return master, nil
}
