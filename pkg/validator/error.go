package validator

type ErrorResponse struct {
	Errors  []errorDetails `json:"errors,omitempty"`
	Message string         `json:"message,omitempty"`
}

type errorDetails struct {
	Field   string `json:"field,omitempty"`
	Value   any    `json:"value,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *ErrorResponse) AddError(field string, value any, msg string) {
	e.Errors = append(e.Errors, errorDetails{
		Field:   field,
		Value:   value,
		Message: msg,
	})
}
