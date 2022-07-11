package email

import (
	"context"
	"errors"
	"fmt"
	"maskan/client/jtrace"
	"maskan/contract"
	merror "maskan/error"
	model "maskan/src/email/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type EmailService struct {
	otpService      contract.IOtpService
	emailRepository contract.IEmailRepository
}

type EmailServiceParams struct {
	fx.In
	OtpService      contract.IOtpService
	EmailRepository contract.IEmailRepository
}

func NewEmailService(params EmailServiceParams) contract.IEmailService {
	return &EmailService{
		emailRepository: params.EmailRepository,
		otpService:      params.OtpService,
	}
}

func (e EmailService) GetUserEmail(c context.Context, userId uuid.UUID) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[GetUserEmail]")
	defer span.Finish()
	return e.emailRepository.Get(c, map[string]any{"user_id": userId})
}

func (e EmailService) AddEmail(c context.Context, userId uuid.UUID, email string) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[AddEmail]")
	defer span.Finish()

	if err := e.emailRepository.Exists(c, map[string]any{"email": email}); err != nil {
		return model.Email{}, err
	}

	emailRecord, err := e.emailRepository.Add(c, userId, email)
	if err != nil {
		return model.Email{}, err
	}

	return emailRecord, nil
}

func (e EmailService) SendOtpEmail(c context.Context, emailId uint) error {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[SendOtpEmail]")
	defer span.Finish()

	code, err := e.otpService.NewCode(c, emailId)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("The code is %s", code)

	if err := e.emailRepository.Send(c, []string{}, message); err != nil {
		return err
	}

	return nil
}

func (e EmailService) AproveEmail(c context.Context, userId uuid.UUID, email string) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[AproveEmail]")
	defer span.Finish()

	emailRecord, err := e.emailRepository.Get(c, map[string]any{"email": email})
	if err != nil {
		return model.Email{}, err
	}

	if emailRecord.UserId == userId {
		return e.emailRepository.Update(c, model.Email{
			ID:       emailRecord.ID,
			Verified: true,
		})
	}

	return model.Email{}, merror.ErrEmailDoesntBelongToUser
}

func (e EmailService) EmailExists(c context.Context, email string) error {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[EmailExists]")
	defer span.Finish()

	if err := e.emailRepository.Exists(c, map[string]any{"email": email, "verified": true}); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return merror.ErrEmailExists
}
