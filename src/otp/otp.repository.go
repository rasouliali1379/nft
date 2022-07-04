package otp

import (
	"maskan/contract"

	"go.uber.org/fx"
)

type OtpRepository struct {
}

type OtpRepositoryParams struct {
	fx.In
}

func NewOtpRepository(params OtpRepositoryParams) contract.IOtpRepository {
	return &OtpRepository{}
}
