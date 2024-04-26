package services

import (
	"database/sql"
	"errors"

	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	if !errors.Is(err, sql.ErrNoRows) {
		if m != nil {
			return nil, apierror.EmailAlreadyTaken()
		}
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(createMasterDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	master := &models.Master{
		ID:           uuid.NewString(),
		Name:         createMasterDto.Name,
		Email:        createMasterDto.Email,
		PasswordHash: string(hashed),
	}

	if err := service.dao.Save(master); err != nil {
		return nil, apierror.InternalServerError(err)
	}

	return master, nil
}
