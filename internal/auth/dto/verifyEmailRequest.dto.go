package dto

type VerifyEmailRequest struct {
	Token string `json:"token"  validate:"required"`
	Code  string `json:"code"  validate:"required"`
}
