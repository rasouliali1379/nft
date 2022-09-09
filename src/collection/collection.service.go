package collection

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"nft/client/jtrace"
	persist "nft/client/persist/model"
	"nft/config"
	"nft/contract"
	apperrors "nft/error"
	model "nft/src/collection/model"
)

type CollectionService struct {
	fileService          contract.IFileService
	collectionRepository contract.ICollectionRepository
}

type CollectionServiceParams struct {
	fx.In
	CollectionRepository contract.ICollectionRepository
	FileService          contract.IFileService
}

func NewCollectionService(params CollectionServiceParams) contract.ICollectionService {
	return CollectionService{
		collectionRepository: params.CollectionRepository,
		fileService:          params.FileService,
	}
}

func (cs CollectionService) GetCollection(c context.Context, m model.Collection) (model.Collection, error) {
	span, c := jtrace.T().SpanFromContext(c, "CollectionService[GetCollection]")
	defer span.Finish()

	collection, err := cs.collectionRepository.Get(c, persist.Conds{"id": m.ID, "user_id": m.User.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return model.Collection{}, apperrors.ErrCollectionNotFound
		}
		return model.Collection{}, err
	}

	if collection.Status == model.CollectionStatusDraft {
		if collection.HeaderImage == nil {
			return collection, nil
		}
	}

	collection.HeaderImage.Bucket = config.C().Storage.Buckets.Collection
	nftImageUrl, err := cs.fileService.GetImageUrl(c, *collection.HeaderImage)
	if err != nil {
		return model.Collection{}, err
	}
	collection.HeaderImage.FileUrl = nftImageUrl

	return collection, nil
}

func (cs CollectionService) GetAllCollections(c context.Context, query model.QueryCollection) ([]model.Collection, error) {
	span, c := jtrace.T().SpanFromContext(c, "CollectionService[GetAllCollections]")
	defer span.Finish()

	collections, err := cs.collectionRepository.GetAll(c, persist.Conds{"user_id": query.UserId})
	if err != nil {
		return nil, err
	}

	for i, collection := range collections {
		if collection.HeaderImage == nil {
			continue
		}

		collection.HeaderImage.Bucket = config.C().Storage.Buckets.NFT
		nftUrl, err := cs.fileService.GetImageUrl(c, *collection.HeaderImage)
		if err != nil {
			return nil, err
		}

		collections[i].HeaderImage.FileUrl = nftUrl
	}

	return collections, nil
}

func (cs CollectionService) AddCollection(c context.Context, m model.Collection) (model.Collection, error) {
	span, c := jtrace.T().SpanFromContext(c, "CollectionService[AddCollection]")
	defer span.Finish()

	if m.Status == model.CollectionStatusDraft {
		if m.ID != nil {
			nftModel, err := cs.collectionRepository.Get(c, persist.Conds{"id": m.ID.String()})
			if err != nil {
				return model.Collection{}, apperrors.ErrCollectionDraftNotFound
			}

			if nftModel.Status != model.CollectionStatusDraft {
				return model.Collection{}, apperrors.ErrCollectionIsNotDraft
			}

			err = cs.collectionRepository.HardDelete(c, *m.ID)
			if err != nil {
				return model.Collection{}, err
			}
		}
	}

	if m.HeaderImage != nil {
		m.HeaderImage.Bucket = config.C().Storage.Buckets.Collection
		nftFileName, err := cs.fileService.UploadImage(c, *m.HeaderImage)
		if err != nil {
			return model.Collection{}, err
		}
		m.HeaderImage.FileName = nftFileName
	}

	nftModel, err := cs.collectionRepository.Add(c, m)
	if err != nil {
		return model.Collection{}, err
	}

	return cs.GetCollection(c, model.Collection{ID: nftModel.ID, User: m.User})
}

func (cs CollectionService) DeleteCollection(c context.Context, m model.Collection) error {
	span, c := jtrace.T().SpanFromContext(c, "CollectionService[DeleteCollection]")
	defer span.Finish()

	_, err := cs.GetCollection(c, m)
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrCollectionNotFound
		}
		return err
	}

	return cs.collectionRepository.Delete(c, m)
}
