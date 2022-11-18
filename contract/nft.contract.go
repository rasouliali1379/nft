package contract

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"nft/infra/persist/type"
	model "nft/internal/nft/model"
)

type INftController interface {
	Create(c *fiber.Ctx) error
	GetNft(c *fiber.Ctx) error
	GetNftList(c *fiber.Ctx) error
	Approve(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error
	DeleteDraft(c *fiber.Ctx) error
}

type INftService interface {
	Create(c context.Context, m model.Nft) (model.Nft, error)
	Approve(c context.Context, m model.Nft) error
	Reject(c context.Context, m model.Nft) error
	GetNft(c context.Context, m model.Nft) (model.Nft, error)
	GetOwnedNft(c context.Context, m model.Nft) (model.Nft, error)
	GetAllNfts(c context.Context, userId uuid.UUID) ([]model.Nft, error)
	DeleteDraft(c context.Context, m model.Nft) error
}

type INftRepository interface {
	Exists(c context.Context, conditions persist.D) error
	Add(c context.Context, kyc model.Nft) (model.Nft, error)
	Update(c context.Context, kyc model.Nft) (model.Nft, error)
	Delete(c context.Context, userId uuid.UUID) error
	HardDelete(c context.Context, id uuid.UUID) error
	Get(c context.Context, conditions persist.D) (model.Nft, error)
	GetAll(c context.Context, conditions persist.D) ([]model.Nft, error)
}
