package nftasset

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/client/jtrace"
	persist "nft/client/persist/model"
	"nft/contract"
	apperrors "nft/error"
	model "nft/src/nftasset/model"
)

type NftService struct {
	fileService   contract.IFileService
	nftRepository contract.INftRepository
}

type NftServiceParams struct {
	fx.In
	FileService   contract.IFileService
	NftRepository contract.INftRepository
}

func NewNftService(params NftServiceParams) contract.INftService {
	return NftService{
		fileService:   params.FileService,
		nftRepository: params.NftRepository,
	}
}

func (n NftService) Create(c context.Context, m model.Nft) (model.Nft, error) {
	span, c := jtrace.T().SpanFromContext(c, "NftService[Create]")
	defer span.Finish()

	if m.Status == model.NftStatusDraft {
		if m.ID != nil {
			nftModel, err := n.nftRepository.Get(c, persist.Conds{"id": m.ID.String()})
			if err != nil {
				return model.Nft{}, apperrors.ErrNftDraftNotFound
			}

			if nftModel.Status != model.NftStatusDraft {
				return model.Nft{}, apperrors.ErrNftIsNotDraft
			}

			err = n.nftRepository.HardDelete(c, *m.ID)
			if err != nil {
				return model.Nft{}, err
			}
		}
	}

	if m.NftImage != nil {
		nftFileName, err := n.fileService.UploadNftImage(c, *m.NftImage)
		if err != nil {
			return model.Nft{}, err
		}
		m.NftImage.FileName = nftFileName
	}

	nftModel, err := n.nftRepository.Add(c, m)
	if err != nil {
		return model.Nft{}, err
	}

	return n.GetNft(c, model.Nft{ID: nftModel.ID, User: nftModel.User})
}

func (n NftService) Approve(c context.Context, m model.Nft) error {
	span, c := jtrace.T().SpanFromContext(c, "NftService[Approve]")
	defer span.Finish()

	nftModel, err := n.nftRepository.Get(c, persist.Conds{"id": m.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrNftNotFound
		}
		return err
	}

	if nftModel.Status != model.NftStatusPending {
		return apperrors.ErrNftNotSubmittedForReview
	}

	nftModel.ApprovedBy = m.ApprovedBy
	nftModel.RejectedBy = nil
	nftModel.RejectionReason = ""

	if _, err := n.nftRepository.Update(c, nftModel); err != nil {
		return err
	}

	return nil
}

func (n NftService) Reject(c context.Context, m model.Nft) error {
	span, c := jtrace.T().SpanFromContext(c, "NftService[Reject]")
	defer span.Finish()

	nftModel, err := n.nftRepository.Get(c, persist.Conds{"id": m.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrNftNotFound
		}
		return err
	}

	if nftModel.Status != model.NftStatusPending {
		return apperrors.ErrNftNotSubmittedForReview
	}

	nftModel.RejectedBy = m.RejectedBy
	nftModel.RejectionReason = m.RejectionReason
	nftModel.ApprovedBy = nil
	if _, err := n.nftRepository.Update(c, nftModel); err != nil {
		return err
	}
	return nil
}

func (n NftService) GetNft(c context.Context, m model.Nft) (model.Nft, error) {
	span, c := jtrace.T().SpanFromContext(c, "NftService[GetNft]")
	defer span.Finish()

	nftModel, err := n.nftRepository.Get(c, persist.Conds{"id": m.ID, "user_id": m.User.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return model.Nft{}, apperrors.ErrNftNotFound
		}
		return model.Nft{}, err
	}

	if nftModel.Status == model.NftStatusDraft {
		if nftModel.NftImage == nil {
			return nftModel, nil
		}
	}

	nftImageUrl, err := n.fileService.GetNftImageUrl(c, nftModel.NftImage.FileName)
	if err != nil {
		return model.Nft{}, err
	}
	nftModel.NftImage.FileUrl = nftImageUrl

	return nftModel, nil
}

func (n NftService) GetAllNfts(c context.Context, userId uuid.UUID) ([]model.Nft, error) {
	span, c := jtrace.T().SpanFromContext(c, "NftService[GetAllNfts]")
	defer span.Finish()

	nfts, err := n.nftRepository.GetAll(c, persist.Conds{"user_id": userId})
	if err != nil {
		return nil, err
	}

	for i, nft := range nfts {
		if nft.NftImage == nil {
			continue
		}
		nftUrl, err := n.fileService.GetNftImageUrl(c, nft.NftImage.FileName)
		if err != nil {
			return nil, err
		}

		nfts[i].NftImage.FileUrl = nftUrl
	}

	return nfts, nil
}

func (n NftService) DeleteDraft(c context.Context, m model.Nft) error {
	span, c := jtrace.T().SpanFromContext(c, "NftService[DeleteDraft]")
	defer span.Finish()

	_, err := n.nftRepository.Get(c, persist.Conds{"id": m.ID, "user_id": m.User.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrNftNotFound
		}
		return err
	}

	return n.nftRepository.Delete(c, *m.ID)
}
