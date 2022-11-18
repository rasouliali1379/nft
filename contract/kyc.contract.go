package contract

import (
	"context"
	"nft/infra/persist/type"
	model "nft/internal/kyc/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IKycController interface {
	Appeal(c *fiber.Ctx) error
	Approve(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error
	GetAppeal(c *fiber.Ctx) error
	GetAllAppeals(c *fiber.Ctx) error
}

type IKycService interface {
	Appeal(c context.Context, m model.Kyc) (model.Kyc, error)
	Approve(c context.Context, m model.Kyc) error
	Reject(c context.Context, m model.Kyc) error
	GetAppeal(c context.Context, m model.Kyc) (model.Kyc, error)
	GetAllAppeals(c context.Context, m model.Kyc) ([]model.Kyc, error)
}

type IKycRepository interface {
	Exists(c context.Context, conditions persist.D) error
	Add(c context.Context, kyc model.Kyc) (model.Kyc, error)
	Update(c context.Context, kyc model.Kyc) (model.Kyc, error)
	Delete(c context.Context, userId uuid.UUID) error
	Get(c context.Context, conditions persist.D) (model.Kyc, error)
	GetAll(c context.Context, conditions persist.D) ([]model.Kyc, error)
}
