package contract

import (
	"context"
	"nft/client/persist/model"
	model "nft/src/email/model"

	"github.com/google/uuid"
)

type IEmailRepository interface {
	Get(c context.Context, conditions persist.Conds) (model.Email, error)
	Last(c context.Context, conditions persist.Conds) (model.Email, error)
	Add(c context.Context, userId uuid.UUID, email string) (model.Email, error)
	Update(c context.Context, emailModel model.Email) (model.Email, error)
	Send(c context.Context, receivers []string, message string) error
	Exists(c context.Context, conditions persist.Conds) (bool, error)
}

type IEmailService interface {
	EmailExists(c context.Context, email string) (bool, error)
	GetEmail(c context.Context, email string) (model.Email, error)
	GetUserEmail(c context.Context, userId uuid.UUID) (model.Email, error)
	AddEmail(c context.Context, userId uuid.UUID, email string) (model.Email, error)
	ApproveEmail(c context.Context, userId uuid.UUID, email string) error
	SendOtpEmail(c context.Context, emailId uint) error
	GetLastVerifiedEmail(c context.Context, userId uuid.UUID) (model.Email, error)
}
