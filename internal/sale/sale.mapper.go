package sale

import (
	"github.com/google/uuid"
	colmapper "nft/internal/collection"
	collectiondto "nft/internal/collection/dto"
	collection "nft/internal/collection/model"
	nftmapper "nft/internal/nft"
	nftdto "nft/internal/nft/dto"
	nft "nft/internal/nft/model"
	dto "nft/internal/sale/dto"
	"nft/internal/sale/entity"
	"nft/internal/sale/model"
	usermodel "nft/internal/user/model"
)

func mapSaleModelToEntity(m model.Sale) entity.Sale {

	var assetId uuid.UUID

	if m.Nft != nil {
		assetId = *m.Nft.ID
	} else if m.Collection != nil {
		assetId = *m.Collection.ID
	}

	return entity.Sale{
		UserId:     m.User.ID,
		Expiration: m.Expiration,
		SaleType:   entity.SaleType(m.SaleType),
		AssetType:  entity.AssetType(m.AssetType),
		AssetId:    assetId,
		MinPrice:   m.MinPrice,
	}
}

func mapSaleEntityToModel(e entity.Sale) model.Sale {
	var canceledBy usermodel.User
	if e.CanceledBy != nil {
		canceledBy.ID = *e.CanceledBy
	}

	var col collection.Collection
	var nftModel nft.Nft

	switch e.AssetType {
	case entity.AssetTypeNft:
		nftModel.ID = &e.AssetId
	case entity.AssetTypeCollection:
		col.ID = &e.AssetId
	}

	return model.Sale{
		ID:         &e.ID,
		User:       usermodel.User{ID: e.UserId},
		Expiration: e.Expiration,
		CanceledBy: &canceledBy,
		CanceledAt: e.CanceledAt,
		Collection: &col,
		Nft:        &nftModel,
		MinPrice:   e.MinPrice,
		SaleType:   model.Type(e.SaleType),
		AssetType:  model.AssetType(e.AssetType),
	}
}

func mapCreateSaleDtoToModel(request dto.SaleRequest, userId uuid.UUID) model.Sale {
	return model.Sale{
		User:     usermodel.User{ID: userId},
		Nft:      &nft.Nft{ID: &request.AssetId},
		MinPrice: request.MinPrice,
		SaleType: model.Type(request.SaleType),
	}
}

func mapCreateSaleModelToDto(sale model.Sale) dto.SaleResponse {
	return dto.SaleResponse{
		SaleId:     *sale.ID,
		Expiration: sale.Expiration.Unix(),
	}
}

func createSalesListDtoFromModel(list []model.Sale) dto.SaleList {
	saleList := make([]dto.Sale, len(list))

	for i := range list {
		saleList[i] = mapSaleModelToDto(list[i])
	}

	return dto.SaleList{Sales: saleList}
}

func mapSaleModelToDto(sale model.Sale) dto.Sale {
	var col *collectiondto.Collection
	var nftDto *nftdto.Nft

	if sale.Nft != nil {
		temp := nftmapper.MapNftModelToDto(*sale.Nft)
		nftDto = &temp
	} else if sale.Collection != nil {
		temp := colmapper.MapCollectionModelToDto(*sale.Collection)
		col = &temp
	}

	return dto.Sale{
		ID:         sale.ID.String(),
		Expiration: sale.Expiration.Unix(),
		Collection: col,
		Nft:        nftDto,
		MinPrice:   sale.MinPrice,
		SaleType:   dto.SaleType(sale.SaleType),
		AssetType:  dto.AssetType(sale.AssetType),
		Status:     dto.Status(sale.Status),
	}
}

func createModelSaleListFromEntity(sales []entity.Sale) []model.Sale {
	saleList := make([]model.Sale, len(sales))

	for i := range sales {
		saleList[i] = mapSaleEntityToModel(sales[i])
	}

	return saleList
}
