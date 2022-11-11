package collection

import (
	"database/sql"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	category "nft/internal/category/model"
	dto "nft/internal/collection/dto"
	entity "nft/internal/collection/entity"
	model "nft/internal/collection/model"
	file "nft/internal/file/model"
	user "nft/internal/user/model"
	"nft/pkg/validator"
	"strconv"
)

func mapAndValidateAddCollectionForm(form *multipart.Form, userId uuid.UUID) (model.Collection, validator.ErrorResponse) {
	var collectionModel model.Collection
	var errs validator.ErrorResponse

	draftArr, ok := form.Value["draft"]
	if !ok {
		errs.AddError("draft", nil, "This field is required")
		return model.Collection{}, errs
	}

	draft, err := strconv.ParseBool(draftArr[0])
	if err != nil {
		errs.AddError("draft", draft, "invalid value")
		return model.Collection{}, errs
	}

	if draft {
		collectionModel.Status = model.CollectionStatusDraft
	} else {
		collectionModel.Status = model.CollectionStatusSaved
	}

	id, ok := form.Value["id"]
	if ok {
		nftId, err := uuid.Parse(id[0])
		if err != nil {
			errs.AddError("id", id[0], "invalid nft id")
		}
		collectionModel.ID = &nftId
	}

	nftImage, ok := form.File["header_image"]
	if ok {
		imageFile, err := nftImage[0].Open()
		if err != nil {
			errs.AddError("header_image", nil, "unable to to process image file")
		}

		nftBytes, err := io.ReadAll(imageFile)
		if err != nil {
			errs.AddError("header_image", nil, "unable to to process image file")
		}
		collectionModel.HeaderImage = &file.Image{Content: nftBytes, FileName: nftImage[0].Filename}
	} else {
		if !draft {
			errs.AddError("header_image", nil, "unable to get header_image from multipart form")
		}
	}

	title, ok := form.Value["title"]
	if ok {
		collectionModel.Title = title[0]
	} else {
		if !draft {
			errs.AddError("title", nil, "unable to get title from multipart form")
		}
	}

	desc, ok := form.Value["description"]
	if ok {
		collectionModel.Description = desc[0]
	} else {
		if !draft {
			errs.AddError("description", nil, "unable to get description from multipart form")
		}
	}

	categoryIds, ok := form.Value["category_id"]
	if ok {
		var categories []category.Category
		for _, c := range categoryIds {
			catId, err := uuid.Parse(c)
			if err != nil {
				errs.AddError("category_id", c, "invalid category id")
			}
			categories = append(categories, category.Category{ID: catId})
		}
		collectionModel.Categories = categories
	} else {
		if !draft {
			errs.AddError("category_id", nil, "unable to get category_ids from multipart form")
		}
	}

	if draft && len(desc) < 1 && len(title) < 1 && collectionModel.HeaderImage == nil {
		errs.AddError("", nil, "you need to provide a title, description or header image to save draft")
		return model.Collection{}, errs
	}

	collectionModel.User = user.User{ID: userId}

	return collectionModel, errs
}

func mapCollectionModelToDto(m model.Collection) dto.Collection {
	var collectionDto dto.Collection
	collectionDto.ID = m.ID.String()

	if m.HeaderImage != nil {
		collectionDto.HeaderImage = m.HeaderImage.FileUrl
	}

	collectionDto.Title = m.Title
	collectionDto.Description = m.Description
	collectionDto.Status = string(m.Status)

	return collectionDto
}

func createCollectionListDtoFromModel(collections []model.Collection) dto.CollectionList {
	collectionList := make([]dto.Collection, len(collections))

	for i, collection := range collections {
		collectionList[i] = mapCollectionModelToDto(collection)
	}

	return dto.CollectionList{Collections: collectionList}
}

func mapCollectionEntityToModel(collection entity.Collection) model.Collection {

	var m model.Collection

	m.ID = &collection.ID

	if collection.Title != nil {
		m.Title = collection.Title.String
	}

	if collection.Description != nil {
		m.Description = collection.Description.String
	}

	categories := make([]category.Category, len(collection.CategoryIds))
	for i, id := range collection.CategoryIds {
		catId, _ := uuid.Parse(id)
		categories[i] = category.Category{ID: catId}
	}

	m.Categories = categories

	if collection.Draft {
		m.Status = model.CollectionStatusDraft
	} else {
		m.Status = model.CollectionStatusSaved
	}

	if collection.HeaderImage != nil {
		m.HeaderImage = &file.Image{FileName: collection.HeaderImage.String}
	}

	return m
}

func mapCollectionModelToEntity(m model.Collection) entity.Collection {
	var collection entity.Collection
	collection.UserId = m.User.ID

	var catIds []string
	for _, cat := range m.Categories {
		catIds = append(catIds, cat.ID.String())
	}
	collection.CategoryIds = catIds

	status := false
	if m.Status == model.CollectionStatusDraft {
		status = true
	}
	collection.Draft = status

	if m.ID != nil {
		collection.ID = *m.ID
	}

	if len(m.Title) > 0 {
		collection.Title = &sql.NullString{String: m.Title, Valid: true}
	}

	if len(m.Description) > 0 {
		collection.Description = &sql.NullString{String: m.Description, Valid: true}
	}

	if m.HeaderImage != nil {
		collection.HeaderImage = &sql.NullString{String: m.HeaderImage.FileName, Valid: true}
	}

	return collection
}

func createModelCollectionListFromEntity(collections []entity.Collection) []model.Collection {
	nftList := make([]model.Collection, len(collections))
	for i, nft := range collections {
		nftList[i] = mapCollectionEntityToModel(nft)
	}
	return nftList
}
