package collection

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"log"
	"nft/client/jtrace"
	"nft/contract"
	apperrors "nft/error"
	"nft/pkg/filper"
	model "nft/src/collection/model"
	usermodel "nft/src/user/model"
)

type CollectionController struct {
	collectionService contract.ICollectionService
}

type CollectionControllerParams struct {
	fx.In
	CollectionService contract.ICollectionService
}

func NewCollectionController(params CollectionControllerParams) contract.ICollectionController {
	return &CollectionController{
		collectionService: params.CollectionService,
	}
}

func (co CollectionController) Add(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CollectionController[Add]")
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

	collection, errRes := mapAndValidateAddCollectionForm(form, userId)

	if len(errRes.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	collectionModel, err := co.collectionService.AddCollection(ctx, collection)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidFileExtension) {
			return filper.GetBadRequestError(c, err.Error())
		}
		log.Println(err)
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapCollectionModelToDto(collectionModel))
}

func (co CollectionController) Get(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CollectionController[Get]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	collectionId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid collection id")
	}

	collectionModel, err := co.collectionService.GetCollection(
		ctx,
		model.Collection{ID: &collectionId, User: usermodel.User{ID: userId}},
	)
	if err != nil {
		if errors.Is(err, apperrors.ErrCollectionNotFound) {
			return filper.GetNotFoundError(c, apperrors.ErrCollectionNotFound.Error())
		}
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapCollectionModelToDto(collectionModel))
}

func (co CollectionController) GetAll(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CollectionController[GetAll]")
	defer span.Finish()

	var query model.QueryCollection
	if c.Locals("user_id") != nil {
		userId := c.Locals("user_id").(uuid.UUID)
		query.UserId = &userId
	}

	//Todo add query by categories id

	collections, err := co.collectionService.GetAllCollections(ctx, query)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusOK).JSON(createCollectionListDtoFromModel(collections))
}

func (co CollectionController) Delete(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CollectionController[Delete]")
	defer span.Finish()
	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	collectionId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, apperrors.ErrInvalidNftId.Error())
	}

	err = co.collectionService.DeleteCollection(ctx, model.Collection{ID: &collectionId, User: usermodel.User{ID: userId}})
	if err != nil {
		if errors.Is(err, apperrors.ErrCollectionNotFound) {
			return filper.GetNotFoundError(c, apperrors.ErrCollectionNotFound.Error())
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "collection deleted successfully")
}
