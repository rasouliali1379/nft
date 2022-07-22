package email

import (
	"context"
	"errors"
	"fmt"
	"nft/client/jtrace"
	"nft/client/persist/model"
	"nft/contract"
	merror "nft/error"
	entity "nft/src/email/entity"
	model "nft/src/email/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type EmailRepository struct {
	db contract.IPersist
}

type EmailRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewEmailRepository(params EmailRepositoryParams) contract.IEmailRepository {
	return &EmailRepository{
		db: params.DB,
	}
}

func (e EmailRepository) Get(c context.Context, conditions persist.Conds) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailRepository[Get]")
	defer span.Finish()

	emailRecord, err := e.db.Get(c, &entity.Email{}, conditions)
	if err != nil {
		return model.Email{}, err
	}

	return mapEmailEntityToModel(emailRecord.(*entity.Email)), nil
}

func (e EmailRepository) Last(c context.Context, conditions persist.Conds) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailRepository[Last]")
	defer span.Finish()

	emailRecord, err := e.db.Last(c, &entity.Email{}, conditions)
	if err != nil {
		return model.Email{}, err
	}

	return mapEmailEntityToModel(emailRecord.(*entity.Email)), nil
}

func (e EmailRepository) Add(c context.Context, userId uuid.UUID, email string) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailRepository[Add]")
	defer span.Finish()

	emailEntity, err := e.db.Create(c, &entity.Email{UserId: userId, Email: email})
	if err != nil {
		return model.Email{}, err
	}

	return mapEmailEntityToModel(emailEntity.(*entity.Email)), nil
}

func (e EmailRepository) Update(c context.Context, emailModel model.Email) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailRepository[Update]")
	defer span.Finish()

	data := createMapFromEmailModel(emailModel)

	updatedEmail, err := e.db.Update(c, &entity.Email{ID: emailModel.ID}, data)
	if err != nil {
		return model.Email{}, fmt.Errorf("error happened while updating email: %w", err)
	}

	return mapEmailEntityToModel(updatedEmail.(*entity.Email)), nil
}

func (e EmailRepository) Send(c context.Context, receivers []string, message string) error {
	span, c := jtrace.T().SpanFromContext(c, "EmailRepository[Send]")
	defer span.Finish()

	// from := config.C().Smtp.From
	// password := config.C().Smtp.Password
	// smtpHost := config.C().Smtp.Host
	// smtpPort := config.C().Smtp.Port

	// auth := smtp.PlainAuth("", from, password, smtpHost)

	// return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, receivers, []byte(message))

	return nil
}

func (e EmailRepository) Exists(c context.Context, conditions persist.Conds) (bool, error) {
	span, c := jtrace.T().SpanFromContext(c, "EmailRepository[Exists]")
	defer span.Finish()

	if _, err := e.db.Get(c, &entity.Email{}, conditions); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
