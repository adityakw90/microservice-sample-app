package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator for request validation.
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new Validator instance.
func NewValidator() *Validator {
	v := validator.New()
	return &Validator{
		validate: v,
	}
}

// ValidateStruct validates a struct and returns validation errors.
func (v *Validator) ValidateStruct(s interface{}) map[string]string {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = getErrorMsg(err)
	}
	return errors
}

// getErrorMsg returns a human-readable error message for a validation tag.
func getErrorMsg(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Should be at least " + fieldError.Param() + " characters"
	case "max":
		return "Should be at most " + fieldError.Param() + " characters"
	default:
		return "Invalid value"
	}
}
