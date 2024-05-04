package apierror

import "net/http"

type ApiError interface {
	Error() string
	Status() int
}

type apiError struct {
	errorMessage string
	status       int
}

func (a *apiError) Error() string {
	return a.errorMessage
}

func (a *apiError) Status() int {
	return a.status
}

func ErrEmailAlreadyTaken() error {
	return &apiError{errorMessage: "Email already taken", status: http.StatusBadRequest}
}

func ErrInvalidRequest(errorMessage string) error {
	return &apiError{errorMessage: errorMessage, status: http.StatusBadRequest}
}

func ErrInternalServerError(errorMessage string) error {
	return &apiError{errorMessage: errorMessage, status: http.StatusInternalServerError}
}

func ErrIncorrectCredentials() error {
	return &apiError{errorMessage: "Incorrect credentials", status: http.StatusUnauthorized}
}

func ErrMasterNotFound() error {
	return &apiError{errorMessage: "Master not found for the given token", status: http.StatusBadRequest}
}

func ErrPasswordNotFound() error {
	return &apiError{errorMessage: "Password not found", status: http.StatusNotFound}
}

func ErrInvalidKey() error {
	return &apiError{errorMessage: "Invalid key", status: http.StatusBadRequest}
}

func ErrKeyAlreadyUsed() error {
	return &apiError{errorMessage: "Key already in use", status: http.StatusBadRequest}
}
