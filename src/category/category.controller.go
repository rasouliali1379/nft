package category

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/client/jtrace"
	"nft/contract"
	nerror "nft/error"
	"nft/pkg/filper"
	"nft/pkg/validator"
	dto "nft/src/category/dto"
	category "nft/src/category/model"
	user "nft/src/user/model"
)

type CategoryController struct {
	categoryService contract.ICategoryService
}

type CategoryControllerParams struct {
	fx.In
	CategoryService contract.ICategoryService
}

func NewCategoryController(params CategoryControllerParams) contract.ICategoryController {
	return &CategoryController{
		categoryService: params.CategoryService,
	}
}

// GetCategory godoc
// @Summary  get single category
// @Tags     category
// @Accept   json
// @Produce  json
// @Param    id   path      int  true  "category id that will be retrieved"
// @Success  200  {object}  dto.CategoryDto
// @Router   /v1/category/{id} [get]
func (cat CategoryController) GetCategory(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CategoryController[GetCategory]")
	defer span.Finish()

	catId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid category id")
	}

	catModel, err := cat.categoryService.GetCategory(ctx, catId)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.JSON(mapCategoryModelToDto(catModel))
}

// GetAllCategories godoc
// @Summary  get categories list
// @Tags     category
// @Accept   json
// @Produce  json
// @Success  200  {object}  dto.CategoriesListDto
// @Router   /v1/category [get]
func (cat CategoryController) GetAllCategories(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CategoryController[GetAllCategories]")
	defer span.Finish()

	cats, err := cat.categoryService.GetAllCategories(ctx)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.JSON(createCategoryList(cats))
}

// AddCategory godoc
// @Summary  add category
// @Tags     category
// @Accept   json
// @Produce  json
// @Param    message  body      dto.AddCategoryRequest  true  "Add category request body. Not providing paent_id means the it's a main category"
// @Success  200      {object}  dto.CategoryDto
// @Router   /v1/category [post]
func (cat CategoryController) AddCategory(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CategoryController[AddCategory]")
	defer span.Finish()

	userId := c.Locals("user_id").(uuid.UUID)

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.AddCategoryRequest
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(request)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)

	}

	catModel, err := cat.categoryService.AddCategory(ctx, category.Category{
		Name:     request.Name,
		ParentId: &request.ParentId,
		CreatedBy: user.User{
			ID: userId,
		},
	})

	if err != nil {
		if errors.Is(err, nerror.ErrParentCategoryNotFound) {
			return filper.GetBadRequestError(c, "parent category not found")
		}
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapCategoryModelToDto(catModel))
}

// UpdateCategory godoc
// @Summary  update existing category
// @Tags     category
// @Accept   json
// @Produce  json
// @Param    id       path  int                     true      "category id that will be updated"
// @Param    message  body  dto.AddCategoryRequest  true      "update category request body"
// @Success  200                                    {object}  dto.CategoryDto
// @Router   /v1/category/{id} [patch]
func (cat CategoryController) UpdateCategory(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CategoryController[UpdateCategory]")
	defer span.Finish()

	userId := c.Locals("user_id").(uuid.UUID)

	catId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid category id")
	}

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.AddCategoryRequest
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	catModel, err := cat.categoryService.UpdateCategory(ctx, category.Category{
		ID:       catId,
		Name:     request.Name,
		ParentId: &request.ParentId,
		CreatedBy: user.User{
			ID: userId,
		},
	})

	if err != nil {
		if errors.Is(err, nerror.ErrParentCategoryNotFound) {
			return filper.GetBadRequestError(c, "parent category not found")
		}
		return filper.GetInternalError(c, "")
	}

	return c.JSON(mapCategoryModelToDto(catModel))
}

// DeleteCategory godoc
// @Summary  delete existing category
// @Tags     category
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "category id that will be deleted"
// @Success  200  {string}  string  "category deleted successfully"
// @Router   /v1/category/{id} [delete]
func (cat CategoryController) DeleteCategory(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CategoryController[DeleteCategory]")
	defer span.Finish()

	catId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid category id")
	}

	err = cat.categoryService.DeleteCategory(ctx, catId)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "category deleted successfully")
}
