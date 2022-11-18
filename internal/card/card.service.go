package card

import (
	"context"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
	model "nft/internal/card/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type CardService struct {
	CardRepository contract.ICardRepository
}

type CardServiceParams struct {
	fx.In
	CardRepository contract.ICardRepository
}

func NewCardService(params CardServiceParams) contract.ICardService {
	return CardService{
		CardRepository: params.CardRepository,
	}
}

func (s CardService) GetAllCards(c context.Context, userId uuid.UUID) ([]model.Card, error) {
	span, c := jtrace.T().SpanFromContext(c, "CardService[GetAllCards]")
	defer span.Finish()
	return s.CardRepository.GetAll(c, persist.D{"user_id": userId})
}

func (s CardService) GetCard(c context.Context, id uuid.UUID, userId uuid.UUID) (model.Card, error) {
	span, c := jtrace.T().SpanFromContext(c, "CardService[GetCard]")
	defer span.Finish()

	cardModel, err := s.CardRepository.Get(c, persist.D{"user_id": userId, "id": id})
	if err != nil {
		return model.Card{}, err
	}

	if cardModel.UserId != userId {
		return model.Card{}, apperrors.ErrCardDoesntBelongToUser
	}

	return cardModel, nil
}

func (s CardService) AddCard(c context.Context, cardModel model.Card) (model.Card, error) {
	span, c := jtrace.T().SpanFromContext(c, "CardService[AddCard]")
	defer span.Finish()

	cardExists, err := s.CardRepository.Exists(c, persist.D{
		"user_id":     cardModel.UserId,
		"card_number": cardModel.CardNumber,
	})

	if err != nil {
		return model.Card{}, err
	}

	if cardExists {
		return model.Card{}, apperrors.ErrCardAlreadyExistsForUser
	}

	return s.CardRepository.Add(c, cardModel)
}

func (s CardService) ApproveCard(c context.Context, id uuid.UUID, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "CardService[UpdateCard]")
	defer span.Finish()

	if _, err := s.CardRepository.Exists(c, persist.D{"id": id}); err != nil {
		return err
	}

	if _, err := s.CardRepository.Update(c, model.Card{ID: id, ApprovedBy: &userId}); err != nil {
		return err
	}

	return nil
}

func (s CardService) DeleteCard(c context.Context, id uuid.UUID, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "CardService[DeleteCard]")
	defer span.Finish()

	if _, err := s.CardRepository.Exists(c, persist.D{"id": id, "user_id": userId}); err != nil {
		return err
	}

	return s.CardRepository.Delete(c, id)
}
