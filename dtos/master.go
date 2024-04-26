package dtos

type CreateMasterRequestDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateMasterResponseDto struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type AuthenticateRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateResponseDto struct {
	AccessToken string `json:"access_token"`
}
