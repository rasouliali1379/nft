package otp

import (
	"context"
	"maskan/client/jtrace"
	"maskan/config"
	"maskan/contract"
	entity "maskan/src/otp/entity"
	model "maskan/src/otp/model"

	"github.com/xlzd/gotp"
	"go.uber.org/fx"
)

type OtpRepository struct {
	db contract.IPersist
}

type OtpRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewOtpRepository(params OtpRepositoryParams) contract.IOtpRepository {
	return &OtpRepository{
		db: params.DB,
	}
}

func (o OtpRepository) Generate(c context.Context, index int) string {
	span, c := jtrace.T().SpanFromContext(c, "OtpRepository[Generate]")
	defer span.Finish()
	return gotp.NewDefaultHOTP(config.C().Otp.Secret).At(index)
}

func (o OtpRepository) Validate(c context.Context, code string, index int) bool {
	span, c := jtrace.T().SpanFromContext(c, "OtpRepository[Validate]")
	defer span.Finish()

	if config.C().Env == config.TEST || config.C().Env == config.DEVELOPMENT {
		if code == "111111" {
			return true
		}
		return false
	}
	return gotp.NewDefaultHOTP(config.C().Otp.Secret).Verify(code, index)
}

func (o OtpRepository) Add(c context.Context, otpModel model.Otp) (model.Otp, error) {
	span, c := jtrace.T().SpanFromContext(c, "OtpRepository[Add]")
	defer span.Finish()

	otpEntity, err := o.db.Create(c, &entity.Otp{Code: otpModel.Code, UserEmailId: otpModel.UserEmailId})
	if err != nil {
		return model.Otp{}, err
	}

	return mapOtpEntityToModel(otpEntity.(*entity.Otp)), nil
}

func (o OtpRepository) GetByEmailId(c context.Context, emailId uint) (model.Otp, error) {
	span, c := jtrace.T().SpanFromContext(c, "OtpRepository[Validate]")
	defer span.Finish()

	otpEntity, err := o.db.Get(c, &entity.Otp{}, map[string]any{"user_email_id": emailId})
	if err != nil {
		return model.Otp{}, err
	}

	return mapOtpEntityToModel(otpEntity.(*entity.Otp)), nil
}

func (o OtpRepository) Count(c context.Context, emailId uint) (int, error) {
	span, c := jtrace.T().SpanFromContext(c, "OtpRepository[Count]")
	defer span.Finish()

	count, err := o.db.Count(c, &entity.Otp{}, map[string]any{"user_email_id": emailId})
	if err != nil {
		return 0, err
	}

	return count, nil
}
