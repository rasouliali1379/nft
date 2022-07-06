package contract

import (
	"context"
	"maskan/src/email/model"
)

type IEmailRepository interface {
	GetEmail(c context.Context, email string) (email.Email, error)
	GetEmailByUserId(c context.Context, email string) (email.Email, error)
	AddEmail(c context.Context, userId string, email string) (email.Email, error)
	UpdateEmail(c context.Context, id uint) (email.Email, error)
}

type IEmailService interface {
	// GetUserEmail(c context.Context, userId string) (email.Email, error)
	// AddEmail(c context.Context, userId string, email string) (email.Email, error)
	// AproveEmail(c context.Context, userId string, email string) (email.Email, error)
}
