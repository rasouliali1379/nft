package sale

import (
	"github.com/google/uuid"
	dto "nft/internal/sale/dto"
	"nft/internal/sale/entity"
	"nft/internal/sale/model"
	usermodel "nft/internal/user/model"
)

func mapSaleModelToEntity(m model.Sale) entity.Sale {
	return entity.Sale{
		UserId:     m.User.ID,
		Expiration: m.Expiration,
		SaleType:   entity.SaleType(m.SaleType),
		AssetType:  entity.AssetType(m.AssetType),
		AssetId:    m.AssetId,
		MinPrice:   m.MinPrice,
	}
}

func mapSaleEntityToModel(e entity.Sale) model.Sale {
	var canceledBy usermodel.User
	if e.CanceledBy != nil {
		canceledBy.ID = *e.CanceledBy
	}

	return model.Sale{
		ID:         &e.ID,
		User:       usermodel.User{ID: e.UserId},
		Expiration: e.Expiration,
		CanceledBy: &canceledBy,
		CanceledAt: e.CanceledAt,
		AssetId:    e.AssetId,
		MinPrice:   e.MinPrice,
		SaleType:   model.SaleType(e.SaleType),
		AssetType:  model.AssetType(e.AssetType),
	}
}

func mapCreateSaleDtoToModel(request dto.SaleRequest, userId uuid.UUID) model.Sale {
	return model.Sale{
		User:     usermodel.User{ID: userId},
		AssetId:  request.NftId,
		MinPrice: request.MinPrice,
	}
}

func mapCreateSaleModelToDto(sale model.Sale) dto.SaleResponse {
	return dto.SaleResponse{
		SaleId:     *sale.ID,
		Expiration: sale.Expiration.Unix(),
	}
}
