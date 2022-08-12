package contract

import (
	"context"
	persist "nft/client/persist/model"
	model "nft/src/kyc/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IKYCController interface {
	Appeal(c *fiber.Ctx) error
	Approve(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error
	GetAppeal(c *fiber.Ctx) error
	GetAllAppeals(c *fiber.Ctx) error
}

type IKYCService interface {
	Appeal(c context.Context, m model.KYC) (model.KYC, error)
	Approve(c context.Context, m model.KYC) error
	Reject(c context.Context, m model.KYC) error
	GetAppeal(c context.Context, m model.KYC) (model.KYC, error)
	GetAllAppeals(c context.Context, m model.KYC) ([]model.KYC, error)
}

type IKYCRepository interface {
	Exists(c context.Context, conditions persist.Conds) error
	Add(c context.Context, kyc model.KYC) (model.KYC, error)
	Update(c context.Context, kyc model.KYC) (model.KYC, error)
	Delete(c context.Context, userId uuid.UUID) error
	Get(c context.Context, conditions persist.Conds) (model.KYC, error)
	GetAll(c context.Context, conditions persist.Conds) ([]model.KYC, error)
}
