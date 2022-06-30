package contract

import (
	"context"
	model "maskan/src/jwt/model"
	"time"
)

type IJwtRepository interface {
	GenerateToken(c context.Context, userId string, expirationTime time.Time) (string, error)
	Validate(c context.Context, token string) (string, error)
	SaveToken(c context.Context, token string, userId string) error
	UpdateToken(c context.Context, id uint, token string) error
	RetrieveToken(c context.Context, token string) (model.RefreshToken, error)
}

type IJwtService interface {
	Generate(c context.Context, userId string) (model.Jwt, error)
	Validate(c context.Context, accessToken string) error
	Refresh(c context.Context, refreshToken string) (model.Jwt, error)
}
