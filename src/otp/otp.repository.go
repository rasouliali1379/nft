package otp

import (
	"maskan/config"
	"maskan/contract"

	"github.com/xlzd/gotp"
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

func (o OtpRepository) Generate(index int) string {
	return gotp.NewDefaultHOTP(config.C().Otp.Secret).At(index)
}

func (o OtpRepository) Validate(code string, index int) bool {
	return gotp.NewDefaultHOTP(config.C().Otp.Secret).Verify(code, index)
}
