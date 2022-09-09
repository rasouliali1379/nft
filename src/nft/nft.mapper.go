package nft

import (
	"database/sql"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"nft/pkg/validator"
	category "nft/src/category/model"
	file "nft/src/file/model"
	dto "nft/src/nft/dto"
	entity "nft/src/nft/entity"
	model "nft/src/nft/model"
	user "nft/src/user/model"
	"strconv"
)

func mapAndValidateCreateNftForm(form *multipart.Form, userId uuid.UUID) (model.Nft, validator.ErrorResponse) {
	var nftModel model.Nft
	var errs validator.ErrorResponse

	draftArr, ok := form.Value["draft"]
	if !ok {
		errs.AddError("draft", nil, "This field is required")
		return model.Nft{}, errs
	}

	draft, err := strconv.ParseBool(draftArr[0])
	if err != nil {
		errs.AddError("draft", draft, "invalid value")
		return model.Nft{}, errs
	}

	if draft {
		nftModel.Status = model.NftStatusDraft
	} else {
		nftModel.Status = model.NftStatusPending
	}

	id, ok := form.Value["id"]
	if ok {
		nftId, err := uuid.Parse(id[0])
		if err != nil {
			errs.AddError("id", id[0], "invalid nft id")
		}
		nftModel.ID = &nftId
	}

	nftImage, ok := form.File["nft_image"]
	if ok {
		imageFile, err := nftImage[0].Open()
		if err != nil {
			errs.AddError("nft_image", nil, "unable to to process image file")
		}

		nftBytes, err := io.ReadAll(imageFile)
		if err != nil {
			errs.AddError("nft_image", nil, "unable to to process image file")
		}
		nftModel.NftImage = &file.Image{Content: nftBytes, FileName: nftImage[0].Filename}
	} else {
		if !draft {
			errs.AddError("nft_image", nil, "unable to get nft_image from multipart form")
		}
	}

	title, ok := form.Value["title"]
	if ok {
		nftModel.Title = title[0]
	} else {
		if !draft {
			errs.AddError("title", nil, "unable to get title from multipart form")
		}
	}

	desc, ok := form.Value["description"]
	if ok {
		nftModel.Description = desc[0]
	} else {
		if !draft {
			errs.AddError("description", nil, "unable to get description from multipart form")
		}
	}

	//collectionId, ok := form.Message["collection_id"]
	//if ok {
	//	nftModel.CollectionId = collectionId[0]
	//} else {
	//	if !draft {
	//		errs = append(errs, validator.ErrorResponse{Field: "collection_id", Message: "unable to get collection_id from multipart form"})
	//	}
	//}

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
		nftModel.Categories = categories
	} else {
		if !draft {
			errs.AddError("category_id", nil, "unable to get category_ids from multipart form")
		}
	}

	if draft && len(desc) < 1 && len(title) < 1 && nftModel.NftImage == nil {
		errs.AddError("", nil, "you need to provide a title, description or image to save draft")
		return model.Nft{}, errs
	}

	nftModel.User = user.User{ID: userId}

	return nftModel, errs
}

func mapNftModelToDto(m model.Nft) dto.Nft {
	var nftDto dto.Nft
	nftDto.ID = m.ID.String()

	if m.NftImage != nil {
		nftDto.NftImageUrl = m.NftImage.FileUrl
	}

	nftDto.Title = m.Title
	nftDto.Description = m.Description
	nftDto.Status = string(m.Status)
	nftDto.RejectionReason = m.RejectionReason

	return nftDto
}

func mapNftModelToEntity(m model.Nft) entity.Nft {
	var nftEntity entity.Nft
	nftEntity.UserId = m.User.ID

	var catIds []string
	for _, cat := range m.Categories {
		catIds = append(catIds, cat.ID.String())
	}
	nftEntity.CategoryIds = catIds

	status := false
	if m.Status == model.NftStatusDraft {
		status = true
	}
	nftEntity.Draft = status

	if m.ID != nil {
		nftEntity.ID = *m.ID
	}

	if len(m.Title) > 0 {
		nftEntity.Title = &sql.NullString{String: m.Title, Valid: true}
	}

	if len(m.Description) > 0 {
		nftEntity.Description = &sql.NullString{String: m.Description, Valid: true}
	}

	if m.NftImage != nil {
		nftEntity.NftImage = &sql.NullString{String: m.NftImage.FileName, Valid: true}
	}

	if m.ApprovedBy != nil {
		nftEntity.ApprovedBy = &m.ApprovedBy.ID
	}

	if m.RejectedBy != nil {
		nftEntity.RejectedBy = &m.RejectedBy.ID
		nftEntity.RejectionReason = &sql.NullString{String: m.RejectionReason, Valid: true}
	}

	return nftEntity
}

func mapNftEntityToModel(nft entity.Nft) model.Nft {
	var nftModel model.Nft
	var status model.NftStatus

	categories := make([]category.Category, len(nft.CategoryIds))
	for i, id := range nft.CategoryIds {
		catId, _ := uuid.Parse(id)
		categories[i] = category.Category{ID: catId}
	}

	if nft.Draft {
		status = model.NftStatusDraft
	} else if nft.ApprovedBy != nil {
		status = model.NftStatusApproved
	} else if nft.RejectedBy != nil {
		status = model.NftStatusRejected
	} else {
		status = model.NftStatusPending
	}

	if nft.NftImage != nil {
		nftModel.NftImage = &file.Image{FileName: nft.NftImage.String}
	}

	if nft.Title != nil {
		nftModel.Title = nft.Title.String
	}

	if nft.Description != nil {
		nftModel.Description = nft.Description.String
	}

	nftModel.ID = &nft.ID
	nftModel.Status = status
	nftModel.Categories = categories
	nftModel.User = user.User{ID: nft.UserId}

	return nftModel
}

func createNftListDtoFromModel(nfts []model.Nft) dto.NftList {
	nftList := make([]dto.Nft, len(nfts))

	for i, nft := range nfts {
		nftList[i] = mapNftModelToDto(nft)
	}

	return dto.NftList{Nfts: nftList}
}

func createModelNftListFromEntity(nfts []entity.Nft) []model.Nft {
	nftList := make([]model.Nft, len(nfts))
	for i, nft := range nfts {
		nftList[i] = mapNftEntityToModel(nft)
	}
	return nftList
}
