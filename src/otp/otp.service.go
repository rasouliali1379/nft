package otp

import (
	"maskan/contract"

	"go.uber.org/fx"
)

type OtpService struct {
}

type OtpServiceParams struct {
	fx.In
}

func NewOtpService(params OtpServiceParams) contract.IOtpService {
	return &OtpService{}
}

