package category

import (
	"context"
	"nft/infra/persist/type"
	"time"

	"nft/contract"
	"nft/infra/jtrace"
	entity "nft/internal/category/entity"
	model "nft/internal/category/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type CategoryRepository struct {
	db contract.IPersist
}

type CategoryRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewCategoryRepository(params CategoryRepositoryParams) contract.ICategoryRepository {
	return &CategoryRepository{
		db: params.DB,
	}
}

func (cat CategoryRepository) Exists(c context.Context, conditions persist.D) error {
	span, c := jtrace.T().SpanFromContext(c, "CategoryRepository[Exists]")
	defer span.Finish()

	if _, err := cat.db.Get(c, &entity.Category{}, conditions); err != nil {
		return err
	}

	return nil
}

func (cat CategoryRepository) Add(c context.Context, category model.Category) (model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryRepository[Add]")
	defer span.Finish()

	catEntity := mapCategoryModelToEntity(category)
	catEntity.ID = uuid.New()

	newCat, err := cat.db.Create(c, &catEntity)
	if err != nil {
		return model.Category{}, err
	}

	return mapCategoryEntityToModel(newCat.(*entity.Category)), nil
}

func (cat CategoryRepository) Update(c context.Context, catModel model.Category) (model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryRepository[Update]")
	defer span.Finish()

	data := mapCategoryModelToEntity(catModel)

	updatedCat, err := cat.db.Update(c, &entity.Category{ID: catModel.ID}, data)
	if err != nil {
		return model.Category{}, err
	}

	return mapCategoryEntityToModel(updatedCat.(*entity.Category)), nil
}

func (cat CategoryRepository) Delete(c context.Context, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "CategoryRepository[Delete]")
	defer span.Finish()

	if _, err := cat.db.Update(c, &entity.Category{ID: userId}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return err
	}
	return nil
}

func (cat CategoryRepository) Get(c context.Context, conditions persist.D) (model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryRepository[Get]")
	defer span.Finish()

	category, err := cat.db.Get(c, &entity.Category{}, conditions)
	if err != nil {
		return model.Category{}, err
	}
	return mapCategoryEntityToModel(category.(*entity.Category)), nil
}

func (cat CategoryRepository) GetAll(c context.Context, conditions persist.D) ([]model.Category, error) {
	span, c := jtrace.T().SpanFromContext(c, "CategoryRepository[GetAll]")
	defer span.Finish()

	catList, err := cat.db.GetAll(c, &[]entity.Category{}, conditions)
	if err != nil {
		return nil, err
	}

	return createModelCategoriesList(catList.(*[]entity.Category)), nil
}
