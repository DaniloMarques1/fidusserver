package validate

import (
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate
var once sync.Once

func validateMasterPassword(field validator.FieldLevel) bool {
	value := field.Field().String()
	rgString := regexp.MustCompile(`[a-z].*[A-Z]|[A-Z].*[a-z]`)
	rgNumber := regexp.MustCompile(`\d`)
	rgSymbol := regexp.MustCompile(`[^a-zA-Z0-9]`)

	return rgString.MatchString(value) && rgNumber.MatchString(value) && rgSymbol.MatchString(value) && len(value) >= 8

}

func Validate() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterValidation("master_password", validateMasterPassword)
	})

	return validate
}
