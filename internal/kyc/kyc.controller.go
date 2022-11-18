package kyc

import (
	"errors"
	"log"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	dto "nft/internal/kyc/dto"
	kyc "nft/internal/kyc/model"
	"nft/pkg/filper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

type KycController struct {
	kycService contract.IKycService
}

type KycControllerParams struct {
	fx.In
	KycService contract.IKycService
}

func NewKycController(params KycControllerParams) contract.IKycController {
	return &KycController{
		kycService: params.KycService,
	}
}

// Appeal godoc
// @Summary  appeal for Kyc
// @Tags     kyc
// @Accept   multipart/form-data
// @Produce  json
// @Param    id_card   formData  file  true  "Image of user's id card"
// @Param    portrait  formData  file  true  "Image of user holding his id card are other things request by business"
// @Router   /v1/kyc [post]
func (k KycController) Appeal(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KycController[Appeal]")
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

	idCard := form.File["id_card"]
	if len(idCard) < 1 {
		return filper.GetBadRequestError(c, "you need to provide id card image")
	}

	portrait := form.File["portrait"]
	if len(portrait) < 1 {
		return filper.GetBadRequestError(c, "you need to provide an image of user holding his id card")
	}

	kycModel, err := createKycModel(idCard[0], portrait[0], userId)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	appeal, err := k.kycService.Appeal(ctx, kycModel)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidFileExtension) {
			return filper.GetBadRequestError(c, err.Error())
		}
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapKycModelToDto(appeal))
}

// Approve godoc
// @Summary  approve Kyc appeal
// @Tags     kyc
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "appeal id that will be approved"
// @Success  200  {string}  string  "appeal approved successfully"
// @Router   /v1/kyc/{id}/approve [post]
func (k KycController) Approve(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KycController[Approve]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	appealId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid appeal id")
	}

	err = k.kycService.Approve(ctx, kyc.Kyc{ID: appealId, ApprovedBy: &userId})
	if err != nil {
		if errors.Is(err, apperrors.ErrAppealNotFound) {
			return filper.GetNotFoundError(c, "appeal not found")
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "appeal approved successfully")
}

// Reject godoc
// @Summary  reject Kyc appeal
// @Tags     kyc
// @Accept   json
// @Produce  json
// @Param    id       path      int               true  "appeal id that will be rejected"
// @Param    message  body      dto.RejectAppeal  true  "optional rejection message"
// @Success  200      {string}  string            "appeal rejected successfully"
// @Router   /v1/kyc/{id}/reject [post]
func (k KycController) Reject(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KycController[Reject]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	appealId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, apperrors.ErrInvalidAppealId.Error())
	}

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.RejectAppeal
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	err = k.kycService.Reject(ctx, kyc.Kyc{ID: appealId, RejectedBy: &userId, RejectionReason: request.Message})
	if err != nil {
		if errors.Is(err, apperrors.ErrAppealNotFound) {
			return filper.GetNotFoundError(c, "appeal not found")
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "appeal rejected successfully")
}

// GetAppeal godoc
// @Summary  get Kyc appeal
// @Tags     kyc
// @Accept   json
// @Produce  json
// @Param    id   path      int  true  "appeal id that will be retrieved"
// @Success  200  {object}  dto.Kyc
// @Router   /v1/kyc/{id} [get]
func (k KycController) GetAppeal(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KycController[GetAppeal]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	appealId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid appeal id")
	}

	appeal, err := k.kycService.GetAppeal(ctx, kyc.Kyc{ID: appealId, UserId: userId})
	if err != nil {
		log.Println(err)
		if errors.Is(err, apperrors.ErrAppealNotFound) {
			return filper.GetNotFoundError(c, "appeal not found")
		}

		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusOK).JSON(mapKycModelToDto(appeal))
}

// GetAllAppeals godoc
// @Summary  get all Kyc appeals
// @Tags     kyc
// @Accept   json
// @Produce  json
// @Success  200  {object}  dto.KycList
// @Router   /v1/kyc [get]
func (k KycController) GetAllAppeals(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KycController[GetAllAppeals]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	appeal, err := k.kycService.GetAllAppeals(ctx, kyc.Kyc{UserId: userId})
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusOK).JSON(createKycListDtoFromModel(appeal))
}
