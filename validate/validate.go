package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
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
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	})

	return validate
}

func GetValidationErrorMessage(err error) string {
	vErr, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	errMessage := "Validation error on field(s): %s"
	fields := ""
	for idx, fErr := range vErr {
		if idx != len(vErr)-1 {
			fields += fmt.Sprintf("%s, ", fErr.Field())
		} else {
			fields += fmt.Sprintf("%s", fErr.Field())
		}
	}

	return fmt.Sprintf(errMessage, fields)
}
