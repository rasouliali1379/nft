package category

import (
	"context"
	"errors"
	"nft/client/jtrace"
	"nft/contract"
	merror "nft/error"
	model "nft/src/category/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type CategoryService struct {
	categoryRepository contract.ICategoryRepository
	userService        contract.IUserService
}

type CategoryServiceParams struct {
	fx.In
	CategoryRepository contract.ICategoryRepository
	UserService        contract.IUserService
}

func NewCategoryService(params CategoryServiceParams) contract.ICategoryService {
	return CategoryService{
		categoryRepository: params.CategoryRepository,
		userService:        params.UserService,
	}
}

func (cat CategoryService) GetCategory(c context.Context, id uuid.UUID) (model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryService[GetCategory]")
	defer span.Finish()

	catModel, err := cat.categoryRepository.Get(c, map[string]any{"id": id})
	if err != nil {
		return model.Category{}, err
	}

	createdBy, err := cat.userService.GetUser(c, map[string]any{"id": catModel.CreatedBy.ID})
	if err != nil {
		return model.Category{}, err
	}
	catModel.CreatedBy = createdBy

	catList, err := cat.GetSubCategories(c, catModel.ID)
	if err != nil {
		return model.Category{}, err
	}
	catModel.SubCategories = catList

	return catModel, nil
}

func (cat CategoryService) GetSubCategories(c context.Context, id uuid.UUID) ([]model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryService[GetSubCategories]")
	defer span.Finish()

	catModels, err := cat.categoryRepository.GetAll(c, map[string]any{"parent_id": id})
	if err != nil {
		return nil, err
	}

	for i, catModel := range catModels {
		catList, err := cat.GetSubCategories(c, catModel.ID)
		if err != nil {
			return nil, err
		}
		catModels[i].SubCategories = catList
	}
	return catModels, nil
}

func (cat CategoryService) GetAllCategories(c context.Context) ([]model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryService[GetAllCategories]")
	defer span.Finish()

	catList, err := cat.categoryRepository.GetAll(c, map[string]any{"parent_id": uuid.Nil})
	if err != nil {
		return nil, err
	}

	for i, item := range catList {
		createdBy, err := cat.userService.GetUser(c, map[string]any{"id": item.CreatedBy.ID})
		if err != nil {
			return nil, err
		}
		item.CreatedBy = createdBy

		subCatList, err := cat.GetSubCategories(c, item.ID)
		if err != nil {
			return nil, err
		}
		item.SubCategories = subCatList

		catList[i] = item
	}

	return catList, nil
}

func (cat CategoryService) AddCategory(c context.Context, category model.Category) (model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryService[AddCategory]")
	defer span.Finish()

	if category.ParentId != nil && *category.ParentId != uuid.Nil {
		if err := cat.categoryRepository.Exists(c, map[string]any{"id": category.ParentId}); err != nil {
			if errors.Is(err, merror.ErrRecordNotFound) {
				return model.Category{}, merror.ErrParentCategoryNotFound
			}
			return model.Category{}, err
		}
	}

	return cat.categoryRepository.Add(c, category)
}

func (cat CategoryService) UpdateCategory(c context.Context, category model.Category) (model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryService[UpdateCategory]")
	defer span.Finish()

	catModel, err := cat.GetCategory(c, category.ID)
	if err != nil {
		return model.Category{}, err
	}

	if category.ParentId != nil && catModel.ParentId != category.ParentId && *category.ParentId != uuid.Nil {
		if err := cat.categoryRepository.Exists(c, map[string]any{"id": category.ParentId}); err != nil {
			if errors.Is(err, merror.ErrRecordNotFound) {
				return model.Category{}, merror.ErrParentCategoryNotFound
			}
			return model.Category{}, err
		}
	}

	return cat.categoryRepository.Update(c, category)
}

func (cat CategoryService) DeleteCategory(c context.Context, catId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "CategoryService[DeleteCategory]")
	defer span.Finish()
	if err := cat.categoryRepository.Exists(c, map[string]any{"id": catId}); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return merror.ErrCategoryNotFound
		}
		return err
	}

	return cat.categoryRepository.Delete(c, catId)
}
