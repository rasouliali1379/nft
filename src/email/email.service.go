package email

import (
	"maskan/contract"

	"go.uber.org/fx"
)

type EmailService struct {
}

type EmailServiceParams struct {
	fx.In
}

func NewEmailService(params EmailServiceParams) contract.IEmailService {
	return &EmailService{}
}
