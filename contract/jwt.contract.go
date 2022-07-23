package contract

import (
	"context"
	"nft/client/persist/model"
	model "nft/src/jwt/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IJwtMiddleware interface {
	Handle(c *fiber.Ctx) error
}
type IJwtRepository interface {
	Generate(c context.Context, userId string, expirationTime time.Time) (string, error)
	Validate(c context.Context, token string) (uuid.UUID, error)
	Add(c context.Context, token string, userId string) error
	Update(c context.Context, data model.RefreshToken) error
	Get(c context.Context, conditions persist.Conds) (model.RefreshToken, error)
}

type IJwtService interface {
	Generate(c context.Context, userId string) (model.Jwt, error)
	Validate(c context.Context, token string) (uuid.UUID, error)
	Refresh(c context.Context, refreshToken string) (model.Jwt, error)
	GenereteOtpToken(c context.Context, userId string) (string, error)
	InvokeRefreshToken(c context.Context, refreshToken string) error
	GetToken(c context.Context, refreshToken string) (model.RefreshToken, error)
}
