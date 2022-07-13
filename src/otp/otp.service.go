package otp

import (
	"context"
	"nft/client/jtrace"
	"nft/contract"
	merror "nft/error"
	model "nft/src/otp/model"

	"go.uber.org/fx"
)

type OtpService struct {
	otpRepository contract.IOtpRepository
}

type OtpServiceParams struct {
	fx.In
	OtpRepository contract.IOtpRepository
}

func NewOtpService(params OtpServiceParams) contract.IOtpService {
	return &OtpService{
		otpRepository: params.OtpRepository,
	}
}

func (o OtpService) NewCode(c context.Context, emailId uint) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "OtpService[NewCode]")
	defer span.Finish()

	count, err := o.otpRepository.Count(c, emailId)
	if err != nil {
		return "", err
	}

	code := o.otpRepository.Generate(c, count)

	o.otpRepository.Add(c, model.Otp{
		Code:        code,
		UserEmailId: emailId,
	})

	return code, nil
}

func (o OtpService) ValidateCode(c context.Context, code string, emailId uint) error {
	span, c := jtrace.T().SpanFromContext(c, "OtpService[ValidateCode]")
	defer span.Finish()

	count, err := o.otpRepository.Count(c, emailId)
	if err != nil {
		return err
	}

	if o.otpRepository.Validate(c, code, count) {
		return nil
	}

	return merror.ErrInvalidOtpCode
}
