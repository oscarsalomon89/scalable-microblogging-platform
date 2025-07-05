// Package validator contains a helper to centralize the golang struct validator global reference
package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ValidationErrors = validator.ValidationErrors

var v = validator.New()

func Validate(entity any) error {
	return v.Struct(entity)
}

func validateUUIDFormat(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())
	return err == nil
}

var validationMap = map[string]func(validator.FieldLevel) bool{
	"validUUIDFormat": validateUUIDFormat,
}

var registered bool

func RegisterValidation() (err error) {
	if !registered {
		err = registerValidationMap(validationMap)
		registered = true
	}
	return err
}

func registerValidationMap(m map[string]func(validator.FieldLevel) bool) error {
	for key, value := range m {
		if err := v.RegisterValidation(key, value); err != nil {
			return err
		}
	}

	return nil
}
