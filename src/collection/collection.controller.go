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
	dto "nft/src/collection/dto"
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

// Add godoc
// @Summary  add new collection
// @Tags     collection
// @Accept   multipart/form-data
// @Produce  json
// @Router   /v1/collection [post]
// @Param    id            formData  string   false  "Collection id. Required for updating draft"
// @Param    title         formData  string   false  "Collection title. Not required for draft"
// @Param    description   formData  string   false  "Collection description. Not required for draft"
// @Param    draft         formData  boolean  true   "Collection submission type. If it's true it will be saved as draft. If it's false it will be submitted to be processed."
// @Param    category_id   formData  array    false  "Collection category or sub category id."
// @Param    header_image  formData  file     true   "Collection header image"
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

// Get godoc
// @Summary  get collection
// @Tags     collection
// @Accept   json
// @Produce  json
// @Param    id   path      string  true  "collection id that will be retrieved"
// @Success  200  {object}  dto.Collection
// @Router   /v1/collection/{id} [get]
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

	var collectionDto dto.Collection = mapCollectionModelToDto(collectionModel)

	return c.Status(fiber.StatusCreated).JSON(collectionDto)
}

// GetAll godoc
// @Summary  get all collections
// @Tags     collection
// @Accept   json
// @Produce  json
// @Success  200  {object}  dto.CollectionList
// @Router   /v1/collection [get]
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

// Delete godoc
// @Summary  delete draft collection
// @Tags     collection
// @Accept   json
// @Produce  json
// @Param    id   path      string  true  "collection id that will be deleted"
// @Success  200  {string}  string  "collection deleted successfully"
// @Router   /v1/collection/{id} [delete]
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
