package card

import (
	dto "nft/internal/card/dto"
	entity "nft/internal/card/entity"
	model "nft/internal/card/model"
)

func mapCardModelToDTo(card model.Card) dto.Card {
	return dto.Card{
		ID:         card.ID,
		CardNumber: card.CardNumber,
		IBAN:       card.CardNumber,
		Approved:   card.ApprovedBy != nil,
	}
}

func mapCardListModelToDto(cards []model.Card) dto.CardList {
	cardsList := make([]dto.Card, 0, len(cards))

	for _, card := range cards {
		cardsList = append(cardsList, mapCardModelToDTo(card))
	}

	return dto.CardList{
		Cards: cardsList,
	}
}

func mapCardModelToEntity(cardModel model.Card) entity.Card {
	return entity.Card{
		CardNumber: cardModel.CardNumber,
		IBAN:       cardModel.IBAN,
		UserId:     cardModel.UserId,
		ApprovedBy: cardModel.ApprovedBy,
	}
}

func mapCardEntityToModel(card *entity.Card) model.Card {
	return model.Card{
		ID:         card.ID,
		CardNumber: card.CardNumber,
		IBAN:       card.IBAN,
		UserId:     card.UserId,
		ApprovedBy: card.ApprovedBy,
	}
}

func createCardsListModel(cards *[]entity.Card) []model.Card {
	var cardsList []model.Card

	for _, card := range *cards {
		cardsList = append(cardsList, mapCardEntityToModel(&card))
	}

	return cardsList
}
