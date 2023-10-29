package common_validators

import (
	"github.com/go-playground/validator/v10"
	"net/mail"
)

// NewEmailValidator func for create a new validator for model fields.
func NewEmailValidator() *validator.Validate {
	// Create a new validator
	validate := validator.New()

	// Custom validation for Emails fields.
	_ = validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()

		if _, err := mail.ParseAddress(field); err != nil {
			return true
		}

		return false
	})

	return validate
}

// EmailValidatorErrors func for show validation errors for each invalid fields.
func EmailValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
