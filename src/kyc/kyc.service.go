package kyc

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"nft/client/jtrace"
	persist "nft/client/persist/model"
	"nft/contract"
	apperrors "nft/error"
	model "nft/src/kyc/model"

	"go.uber.org/fx"
)

type KYCService struct {
	fileService   contract.IFileService
	kycRepository contract.IKYCRepository
}

type KYCServiceParams struct {
	fx.In
	FileService   contract.IFileService
	KYCRepository contract.IKYCRepository
}

func NewKYCService(params KYCServiceParams) contract.IKYCService {
	return KYCService{
		fileService:   params.FileService,
		kycRepository: params.KYCRepository,
	}
}

func (k KYCService) Appeal(c context.Context, m model.KYC) (model.KYC, error) {
	span, c := jtrace.T().SpanFromContext(c, "KYCService[Appeal]")
	defer span.Finish()

	idCardFileName, err := k.fileService.UploadKYCImage(c, m.IdCardImage)
	if err != nil {
		return model.KYC{}, err
	}

	portraitFileName, err := k.fileService.UploadKYCImage(c, m.PortraitImage)
	if err != nil {
		return model.KYC{}, err
	}
	m.IdCardImage.FileName = idCardFileName
	m.PortraitImage.FileName = portraitFileName

	kyc, err := k.kycRepository.Add(c, m)
	if err != nil {
		return model.KYC{}, err
	}

	idCardUrl, err := k.fileService.GetKYCImageUrl(c, m.IdCardImage.FileName)
	if err != nil {
		return model.KYC{}, err
	}

	portraitUrl, err := k.fileService.GetKYCImageUrl(c, m.PortraitImage.FileName)
	if err != nil {
		return model.KYC{}, err
	}

	kyc.IdCardImage.FileUrl = idCardUrl
	kyc.PortraitImage.FileUrl = portraitUrl

	return kyc, nil
}

func (k KYCService) Approve(c context.Context, m model.KYC) error {
	span, c := jtrace.T().SpanFromContext(c, "KYCService[Approve]")
	defer span.Finish()

	kycModel, err := k.kycRepository.Get(c, persist.Conds{"id": m.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrAppealNotFoundError
		}
		return err
	}
	kycModel.ApprovedBy = m.ApprovedBy
	kycModel.RejectedBy = &uuid.Nil
	kycModel.RejectionReason = ""

	if _, err := k.kycRepository.Update(c, kycModel); err != nil {
		return err
	}

	return nil
}

func (k KYCService) Reject(c context.Context, m model.KYC) error {
	span, c := jtrace.T().SpanFromContext(c, "KYCService[Reject]")
	defer span.Finish()

	kycModel, err := k.kycRepository.Get(c, persist.Conds{"id": m.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrAppealNotFoundError
		}
		return err
	}
	kycModel.RejectedBy = m.RejectedBy
	kycModel.RejectionReason = m.RejectionReason
	kycModel.ApprovedBy = &uuid.Nil
	if _, err := k.kycRepository.Update(c, kycModel); err != nil {
		return err
	}

	return nil
}

func (k KYCService) GetAppeal(c context.Context, m model.KYC) (model.KYC, error) {
	span, c := jtrace.T().SpanFromContext(c, "KYCService[GetAppeal]")
	defer span.Finish()

	appeal, err := k.kycRepository.Get(c, persist.Conds{"id": m.ID, "user_id": m.UserId})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return model.KYC{}, apperrors.ErrAppealNotFoundError
		}
		return model.KYC{}, err
	}

	idCardUrl, err := k.fileService.GetKYCImageUrl(c, appeal.IdCardImage.FileName)
	if err != nil {
		return model.KYC{}, err
	}

	portraitUrl, err := k.fileService.GetKYCImageUrl(c, appeal.PortraitImage.FileName)
	if err != nil {
		return model.KYC{}, err
	}

	appeal.IdCardImage.FileUrl = idCardUrl
	appeal.PortraitImage.FileUrl = portraitUrl

	return appeal, nil
}

func (k KYCService) GetAllAppeals(c context.Context, m model.KYC) ([]model.KYC, error) {
	span, c := jtrace.T().SpanFromContext(c, "KYCService[GetAllAppeals]")
	defer span.Finish()

	appeals, err := k.kycRepository.GetAll(c, persist.Conds{})
	if err != nil {
		return nil, err
	}

	for i, appeal := range appeals {
		idCardUrl, err := k.fileService.GetKYCImageUrl(c, appeal.IdCardImage.FileName)
		if err != nil {
			return nil, err
		}

		portraitUrl, err := k.fileService.GetKYCImageUrl(c, appeal.PortraitImage.FileName)
		if err != nil {
			return nil, err
		}

		appeals[i].IdCardImage.FileUrl = idCardUrl
		appeals[i].PortraitImage.FileUrl = portraitUrl
	}

	return appeals, nil
}
