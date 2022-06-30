package validator

import "github.com/go-playground/validator"

func Validate(dto any) []ErrorResponse {
	validate := validator.New()

	var errors []ErrorResponse
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)
		}
	}

	return errors
}
