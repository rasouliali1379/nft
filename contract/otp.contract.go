package contract

import (
	"context"
	model "nft/src/otp/model"
)

type IOtpRepository interface {
	Generate(c context.Context, index int) string
	Validate(c context.Context, code string, index int) bool
	Add(c context.Context, otpModel model.Otp) (model.Otp, error)
	GetByEmailId(c context.Context, emailId uint) (model.Otp, error)
	Count(c context.Context, emailId uint) (int, error)
}

type IOtpService interface {
	NewCode(c context.Context, emailId uint) (string, error)
	ValidateCode(c context.Context, code string, emailId uint) error
}
