package offer

import (
	"nft/internal/offer/dto"
	"nft/internal/offer/entity"
	"nft/internal/offer/model"
	userdto "nft/internal/user/dto"
	usermodel "nft/internal/user/model"
)

func mapOfferModelToEntity(m model.Offer) entity.Offer {
	return entity.Offer{
		UserId:   m.User.ID,
		SaleId:   m.SaleId,
		Price:    m.Price,
		Accepted: m.Accepted,
	}
}

func mapOfferEntityToModel(offer entity.Offer) model.Offer {
	return model.Offer{
		ID:       &offer.ID,
		User:     usermodel.User{ID: offer.UserId},
		SaleId:   offer.SaleId,
		Price:    offer.Price,
		Accepted: offer.Accepted,
	}
}

func createOfferListDtoFromModel(list []model.Offer) dto.OfferList {
	offerList := make([]dto.Offer, len(list))

	for i := range list {
		offerList[i] = mapOfferModelToDto(list[i])
	}

	return dto.OfferList{Offers: offerList}
}

func mapOfferModelToDto(offer model.Offer) dto.Offer {
	return dto.Offer{
		ID:    offer.ID.String(),
		Price: offer.Price,
		User:  userdto.User{},
	}
}
