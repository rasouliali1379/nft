package email

import (
	"context"
	"errors"
	"fmt"
	"nft/client/jtrace"
	"nft/contract"
	"nft/error"
	model "nft/src/email/model"

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
	return EmailService{
		emailRepository: params.EmailRepository,
		otpService:      params.OtpService,
	}
}

func (e EmailService) GetEmail(c context.Context, email string) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[GetEmail]")
	defer span.Finish()

	emailModel, err := e.emailRepository.Get(c, map[string]any{"email": email})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return model.Email{}, apperrors.ErrEmailNotFound
		}
		return model.Email{}, err
	}

	return emailModel, nil
}

func (e EmailService) GetUserEmail(c context.Context, userId uuid.UUID) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[GetUserEmail]")
	defer span.Finish()
	return e.emailRepository.Get(c, map[string]any{"user_id": userId})
}

func (e EmailService) AddEmail(c context.Context, userId uuid.UUID, email string) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[AddEmail]")
	defer span.Finish()

	emailExists, err := e.EmailExists(c, email)
	if err != nil {
		return model.Email{}, err
	}

	if emailExists {
		return model.Email{}, apperrors.ErrEmailExists
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

func (e EmailService) ApproveEmail(c context.Context, userId uuid.UUID, email string) error {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[ApproveEmail]")
	defer span.Finish()

	emailRecord, err := e.emailRepository.Get(c, map[string]any{"email": email})
	if err != nil {
		return err
	}

	if emailRecord.UserId == userId {
		if _, err := e.emailRepository.Update(c, model.Email{ID: emailRecord.ID, Verified: true}); err != nil {
			return err
		}
		return nil
	}

	return apperrors.ErrEmailDoesntBelongToUser
}

func (e EmailService) EmailExists(c context.Context, email string) (bool, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[EmailExists]")
	defer span.Finish()
	return e.emailRepository.Exists(c, map[string]any{"email": email, "verified": true})
}

func (e EmailService) GetLastVerifiedEmail(c context.Context, userId uuid.UUID) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailService[GetLastVerifiedEmail]")
	defer span.Finish()
	return e.emailRepository.Last(c, map[string]any{"user_id": userId, "verified": true})
}
