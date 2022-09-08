package nft

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"nft/client/jtrace"
	"nft/contract"
	apperrors "nft/error"
	"nft/pkg/filper"
	dto "nft/src/nft/dto"
	model "nft/src/nft/model"
	usermodel "nft/src/user/model"
)

type NftController struct {
	nftService contract.INftService
}

type NftControllerParams struct {
	fx.In
	NftService contract.INftService
}

func NewNftController(params NftControllerParams) contract.INftController {
	return &NftController{
		nftService: params.NftService,
	}
}

// Create godoc
// @Summary  create new nft
// @Tags     nft
// @Accept   multipart/form-data
// @Produce  json
// @Router   /v1/nft [post]
// @Param    id             formData  string   false  "Nft id. Required for updating draft"
// @Param    title          formData  string   false  "Nft title. Not required for draft"
// @Param    description    formData  string   false  "Nft description. Not required for draft"
// @Param    draft          formData  boolean  true   "Nft submission type. If it's true it will be saved as draft. If it's false it will be submitted to be processed."
// @Param    category_id    formData  array    false  "Nft category or sub category id."
// @Param    collection_id  formData  string   false  "Nft related collection id"
func (n NftController) Create(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "NftController[Create]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	form, err := c.MultipartForm()
	if err != nil {
		if errors.Is(err, fasthttp.ErrNoMultipartForm) {
			return filper.GetBadRequestError(c, "you to provide multipart form request body")
		}
		return filper.GetInternalError(c, "")
	}

	nftModel, errRes := mapAndValidateCreateNftForm(form, userId)

	if len(errRes.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	createdNft, err := n.nftService.Create(ctx, nftModel)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidFileExtension) {
			return filper.GetBadRequestError(c, err.Error())
		}
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapNftModelToDto(createdNft))
}

// GetNft godoc
// @Summary  get nft
// @Tags     nft
// @Accept   json
// @Produce  json
// @Param    id   path      int  true  "nft id that will be retrieved"
// @Success  200  {object}  dto.Nft
// @Router   /v1/nft/{id} [get]
func (n NftController) GetNft(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "NftController[GetNft]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	nftId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid nft id")
	}

	nftModel, err := n.nftService.GetNft(ctx, model.Nft{
		ID:   &nftId,
		User: usermodel.User{ID: userId},
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrNftNotFound) {
			return filper.GetNotFoundError(c, apperrors.ErrNftNotFound.Error())
		}

		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusOK).JSON(mapNftModelToDto(nftModel))
}

// GetNftList godoc
// @Summary  get all nfts
// @Tags     nft
// @Accept   json
// @Produce  json
// @Success  200  {object}  dto.NftList
// @Router   /v1/nft [get]
func (n NftController) GetNftList(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "NftController[GetNftList]")
	defer span.Finish()
	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	nfts, err := n.nftService.GetAllNfts(ctx, userId)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusOK).JSON(createNftListDtoFromModel(nfts))
}

// Approve godoc
// @Summary  approve nft
// @Tags     nft
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "nft id that will be approved"
// @Success  200  {string}  string  "nft approved successfully"
// @Router   /v1/nft/{id}/approve [post]
func (n NftController) Approve(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "NftController[Approve]")
	defer span.Finish()
	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	nftId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, apperrors.ErrInvalidNftId.Error())
	}

	err = n.nftService.Approve(ctx, model.Nft{ID: &nftId, ApprovedBy: &usermodel.User{ID: userId}})
	if err != nil {
		if errors.Is(err, apperrors.ErrNftNotFound) {
			return filper.GetNotFoundError(c, apperrors.ErrNftNotFound.Error())
		} else if errors.Is(err, apperrors.ErrNftNotSubmittedForReview) {
			return filper.GetBadRequestError(c, apperrors.ErrNftNotSubmittedForReview.Error())
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "nft approved successfully")
}

// Reject godoc
// @Summary  reject nft
// @Tags     nft
// @Accept   json
// @Produce  json
// @Param    id       path      int            true  "nft id that will be rejected"
// @Param    message  body      dto.RejectNft  true  "optional rejection message"
// @Success  200      {string}  string         "nft rejected successfully"
// @Router   /v1/nft/{id}/reject [post]
func (n NftController) Reject(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "NftController[Reject]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	nftId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, apperrors.ErrInvalidNftId.Error())
	}

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.RejectNft
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	err = n.nftService.Reject(ctx, model.Nft{ID: &nftId, RejectedBy: &usermodel.User{ID: userId}, RejectionReason: request.Message})
	if err != nil {
		if errors.Is(err, apperrors.ErrNftNotFound) {
			return filper.GetNotFoundError(c, apperrors.ErrNftNotFound.Error())
		}
		return filper.GetInternalError(c, "")
	} else if errors.Is(err, apperrors.ErrNftNotSubmittedForReview) {
		return filper.GetBadRequestError(c, apperrors.ErrNftNotSubmittedForReview.Error())
	}

	return filper.GetSuccessResponse(c, "nft rejected successfully")
}

// DeleteDraft godoc
// @Summary  delete draft nft
// @Tags     nft
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "nft id that will be deleted"
// @Success  200  {string}  string  "draft deleted successfully"
// @Router   /v1/nft/{id} [delete]
func (n NftController) DeleteDraft(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "NftController[DeleteDraft]")
	defer span.Finish()
	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	nftId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, apperrors.ErrInvalidNftId.Error())
	}

	err = n.nftService.DeleteDraft(ctx, model.Nft{ID: &nftId, User: usermodel.User{ID: userId}})
	if err != nil {
		if errors.Is(err, apperrors.ErrNftNotFound) {
			return filper.GetNotFoundError(c, apperrors.ErrNftNotFound.Error())
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "draft deleted successfully")
}
