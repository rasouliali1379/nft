package card

import (
	"context"
	"errors"
	contract "nft/contract"
	merror "nft/error"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
	entity "nft/internal/card/entity"
	model "nft/internal/card/model"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type CardRepository struct {
	db contract.IPersist
}

type CardRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewCardRepository(params CardRepositoryParams) contract.ICardRepository {
	return &CardRepository{
		db: params.DB,
	}
}

func (r CardRepository) Exists(c context.Context, conditions persist.D) (bool, error) {
	span, c := jtrace.T().SpanFromContext(c, "CardRepository[Exists]")
	defer span.Finish()

	if _, err := r.db.Get(c, &entity.Card{}, conditions); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r CardRepository) Add(c context.Context, cardModel model.Card) (model.Card, error) {
	span, c := jtrace.T().SpanFromContext(c, "CardRepository[Add]")
	defer span.Finish()

	catEntity := mapCardModelToEntity(cardModel)
	catEntity.ID = uuid.New()

	newCat, err := r.db.Create(c, &catEntity)
	if err != nil {
		return model.Card{}, err
	}

	return mapCardEntityToModel(newCat.(*entity.Card)), nil
}

func (r CardRepository) Update(c context.Context, cardModel model.Card) (model.Card, error) {
	span, c := jtrace.T().SpanFromContext(c, "CardRepository[Update]")
	defer span.Finish()

	data := mapCardModelToEntity(cardModel)

	updatedCat, err := r.db.Update(c, &entity.Card{ID: cardModel.ID}, data)
	if err != nil {
		return model.Card{}, err
	}

	return mapCardEntityToModel(updatedCat.(*entity.Card)), nil
}

func (r CardRepository) Delete(c context.Context, cardId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "CardRepository[Delete]")
	defer span.Finish()

	if _, err := r.db.Update(c, &entity.Card{ID: cardId}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return err
	}

	return nil
}

func (r CardRepository) Get(c context.Context, conditions persist.D) (model.Card, error) {
	span, c := jtrace.T().SpanFromContext(c, "CardRepository[Get]")
	defer span.Finish()

	category, err := r.db.Get(c, &entity.Card{}, conditions)
	if err != nil {
		return model.Card{}, err
	}

	return mapCardEntityToModel(category.(*entity.Card)), nil
}

func (r CardRepository) GetAll(c context.Context, conditions persist.D) ([]model.Card, error) {
	span, c := jtrace.T().SpanFromContext(c, "CardRepository[GetAll]")
	defer span.Finish()

	catList, err := r.db.GetAll(c, &[]entity.Card{}, conditions)
	if err != nil {
		return nil, err
	}

	return createCardsListModel(catList.(*[]entity.Card)), nil
}
