package contract

import (
	"context"
	model "maskan/src/jwt/model"
	"time"

	"github.com/google/uuid"
)

type IJwtRepository interface {
	Generate(c context.Context, userId string, expirationTime time.Time) (string, error)
	Validate(c context.Context, token string) (uuid.UUID, error)
	Add(c context.Context, token string, userId string) error
	Update(c context.Context, id uint, token string) error
	Get(c context.Context, token string) (model.RefreshToken, error)
}

type IJwtService interface {
	Generate(c context.Context, userId string) (model.Jwt, error)
	Validate(c context.Context, accessToken string) (uuid.UUID, error)
	Refresh(c context.Context, refreshToken string) (model.Jwt, error)
	GenereteOtpToken(c context.Context, userId string) (string, error)
}
