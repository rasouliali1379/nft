package validator

import "github.com/go-playground/validator"

func Validate(dto any) ErrorResponse {
	validate := validator.New()

	var errRes ErrorResponse
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errRes.AddError(err.StructNamespace(), nil, err.Param())
		}
	}

	return errRes
}
