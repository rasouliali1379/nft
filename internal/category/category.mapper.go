package category

import (
	dto "nft/internal/category/dto"
	entity "nft/internal/category/entity"
	model "nft/internal/category/model"
	usermodel "nft/internal/user/model"

	"github.com/google/uuid"
)

func mapCategoryModelToDto(catModel model.Category) dto.CategoryDto {
	return dto.CategoryDto{
		ID:            catModel.ID,
		Name:          catModel.Name,
		SubCategories: createDtoSubCategoriesList(catModel.SubCategories),
	}
}

func createDtoSubCategoriesList(categories []model.Category) []dto.CategoryDto {
	var subCategories []dto.CategoryDto

	for _, item := range categories {
		subCategories = append(subCategories, mapCategoryModelToDto(item))
	}

	return subCategories
}

func createCategoryList(cats []model.Category) dto.CategoriesListDto {
	categories := make([]dto.CategoryDto, len(cats))

	for i, item := range cats {
		categories[i] = mapCategoryModelToDto(item)
	}

	return dto.CategoriesListDto{
		Categories: categories,
	}
}

func mapCategoryDtoToModel(request dto.CategoryDto) model.Category {

	cat := model.Category{
		Name:          request.Name,
		SubCategories: createSubCategoriesList(request.SubCategories),
	}

	if request.ID != uuid.Nil {
		cat.ID = request.ID
	}

	return cat
}

func createSubCategoriesList(categories []dto.CategoryDto) []model.Category {
	var subCategories []model.Category

	for _, item := range categories {
		subCategories = append(subCategories, mapCategoryDtoToModel(item))
	}

	return subCategories
}

func mapCategoryModelToEntity(catModel model.Category) entity.Category {
	cat := entity.Category{
		Name:      catModel.Name,
		CreatedBy: catModel.CreatedBy.ID,
	}

	if catModel.ParentId != nil {
		cat.ParentId = catModel.ParentId
	}

	return cat
}

func mapCategoryEntityToModel(category *entity.Category) model.Category {
	return model.Category{
		ID:        category.ID,
		Name:      category.Name,
		CreatedBy: usermodel.User{ID: category.CreatedBy},
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func createModelCategoriesList(categories *[]entity.Category) []model.Category {
	var subCategories []model.Category

	for _, item := range *categories {
		subCategories = append(subCategories, mapCategoryEntityToModel(&item))
	}

	return subCategories
}
