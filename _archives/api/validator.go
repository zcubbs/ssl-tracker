package api

import (
	"github.com/go-playground/validator/v10"
	validator2 "github.com/zcubbs/tlz/pkg/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type XValidator struct {
	validator *validator.Validate
}

type ValidationErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

var validDomainName validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if domainName, ok := fieldLevel.Field().Interface().(string); ok {
		return validator2.IsDomaineNameValid(domainName)
	}
	return false
}

func (v XValidator) Validate(data interface{}) []ValidationErrorResponse {
	var validationErrors []ValidationErrorResponse
	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ValidationErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
