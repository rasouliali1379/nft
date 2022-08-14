package kyc

import (
	"errors"
	"log"
	"nft/client/jtrace"
	"nft/contract"
	apperrors "nft/error"
	"nft/pkg/filper"
	dto "nft/src/kyc/dto"
	kyc "nft/src/kyc/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

type KYCController struct {
	kycService contract.IKYCService
}

type KYCControllerParams struct {
	fx.In
	KYCService contract.IKYCService
}

func NewKYCController(params KYCControllerParams) contract.IKYCController {
	return &KYCController{
		kycService: params.KYCService,
	}
}

// Appeal godoc
// @Summary  appeal for KYC
// @Tags     kyc
// @Accept   multipart/form-data
// @Produce  json
// @Param    id_card   formData  file  true  "Image of user's id card"
// @Param    portrait  formData  file  true  "Image of user holding his id card are other things request by business"
// @Router   /v1/kyc [post]
func (k KYCController) Appeal(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KYCController[Appeal]")
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

	kycModel, err := createKYCModel(idCard[0], portrait[0], userId)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	appeal, err := k.kycService.Appeal(ctx, kycModel)
	if err != nil {
		log.Println(err)
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapKYCModelToDto(appeal))
}

// Approve godoc
// @Summary  approve KYC appeal
// @Tags     kyc
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "appeal id that will be approved"
// @Success  200  {string}  string  "appeal approved successfully"
// @Router   /v1/kyc/{id}/approve [post]
func (k KYCController) Approve(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KYCController[Approve]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	appealId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid appeal id")
	}

	err = k.kycService.Approve(ctx, kyc.KYC{ID: appealId, ApprovedBy: &userId})
	if err != nil {
		if errors.Is(err, apperrors.ErrAppealNotFoundError) {
			return filper.GetNotFoundError(c, "appeal not found")
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "appeal approved successfully")
}

// Reject godoc
// @Summary  reject KYC appeal
// @Tags     kyc
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "appeal id that will be rejected"
// @Param    message  body      dto.KYC  true  "optional rejection message"
// @Success  200  {string}  string  "appeal rejected successfully"
// @Router   /v1/kyc/{id}/reject [post]
func (k KYCController) Reject(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KYCController[Reject]")
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

	err = k.kycService.Reject(ctx, kyc.KYC{ID: appealId, RejectedBy: &userId, RejectionReason: request.Message})
	if err != nil {
		if errors.Is(err, apperrors.ErrAppealNotFoundError) {
			return filper.GetNotFoundError(c, "appeal not found")
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "appeal rejected successfully")
}

// GetAppeal godoc
// @Summary  get KYC appeal
// @Tags     kyc
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "appeal id that will be retrieved"
// @Success  200  {object}  dto.KYC
// @Router   /v1/kyc/{id} [get]
func (k KYCController) GetAppeal(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KYCController[GetAppeal]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	appealId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid appeal id")
	}

	appeal, err := k.kycService.GetAppeal(ctx, kyc.KYC{ID: appealId, UserId: userId})
	if err != nil {
		log.Println(err)
		if errors.Is(err, apperrors.ErrAppealNotFoundError) {
			return filper.GetNotFoundError(c, "appeal not found")
		}

		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusOK).JSON(mapKYCModelToDto(appeal))
}

// GetAllAppeals godoc
// @Summary  get all KYC appeals
// @Tags     kyc
// @Accept   json
// @Produce  json
// @Success  200  {object}  dto.KYCList
// @Router   /v1/kyc [get]
func (k KYCController) GetAllAppeals(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "KYCController[GetAllAppeals]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	appeal, err := k.kycService.GetAllAppeals(ctx, kyc.KYC{UserId: userId})
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusOK).JSON(createKYCListDtoFromModel(appeal))
}
