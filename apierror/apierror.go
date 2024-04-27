package apierror

import "net/http"

type ApiError interface {
	Error() string
	Status() int
}

type apiError struct {
	err    string
	status int
}

func EmailAlreadyTaken() error {
	return &apiError{err: "Email already taken", status: http.StatusBadRequest}
}

func InvalidRequestBody(err string) error {
	return &apiError{err: err, status: http.StatusBadRequest}
}

// TODO: better name?
func InternalServerError(err error) error {
	return &apiError{err: err.Error(), status: http.StatusInternalServerError}
}

func MasterEmailNotFound() error {
	return &apiError{err: "incorrect email", status: http.StatusBadRequest}
}

func MasterIncorrectPassword() error {
	return &apiError{err: "incorrect password", status: http.StatusBadRequest}
}

func MasterNotFound() error {
	return &apiError{err: "master not found for the given token", status: http.StatusBadRequest}
}

func InvalidKey() error {
	return &apiError{err: "Invalid key", status: http.StatusBadRequest}
}

func (a *apiError) Error() string {
	return a.err
}

func (a *apiError) Status() int {
	return a.status
}
