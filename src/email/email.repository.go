package email

import (
	"context"
	"go.uber.org/fx"
	"maskan/client/jtrace"
	"maskan/contract"
	entity "maskan/src/email/entity"
	model "maskan/src/email/model"
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

func (e EmailRepository) GetEmail(c context.Context, email string) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[GetEmail]")
	defer span.Finish()

	emailRecord, err := e.db.Get(c, &entity.Email{}, map[string]any{"email": email})
	if err != nil {
		return model.Email{}, err
	}

	return mapEmailEntityToModel(emailRecord.(*entity.Email)), nil
}

func (e EmailRepository) GetEmailByUserId(c context.Context, id string) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[GetEmailByUserId]")
	defer span.Finish()

	emailRecord, err := e.db.Get(c, &entity.Email{}, map[string]any{"id": id})
	if err != nil {
		return model.Email{}, err
	}

	return mapEmailEntityToModel(emailRecord.(*entity.Email)), nil
}

func (e EmailRepository) AddEmail(c context.Context, userId string, email string) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[AddEmail]")
	defer span.Finish()
	return model.Email{}, nil
}

func (e EmailRepository) UpdateEmail(c context.Context, id uint) (model.Email, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[UpdateEmail]")
	defer span.Finish()
	return model.Email{}, nil
}
