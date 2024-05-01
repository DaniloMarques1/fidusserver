package dtos

type StorePasswordRequestDto struct {
	Key      string `json:"key" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RetrievePasswordResponseDto struct {
	MasterId string `json:"master_id"`
	Key      string `json:"key"`
	Password string `json:"password"`
}

type UpdatePasswordRequestDto struct {
	Password string `json:"password" validate:"required"`
}
