package kyc

import (
	"database/sql"
	"io"
	"mime/multipart"
	file "nft/src/file/model"
	dto "nft/src/kyc/dto"
	entity "nft/src/kyc/entity"
	model "nft/src/kyc/model"

	"github.com/google/uuid"
)

func createKycModel(idCard *multipart.FileHeader, portrait *multipart.FileHeader, userId uuid.UUID) (model.Kyc, error) {

	idCardFile, err := idCard.Open()
	if err != nil {
		return model.Kyc{}, err
	}

	idCardBytes, err := io.ReadAll(idCardFile)
	if err != nil {
		return model.Kyc{}, err
	}

	portraitFile, err := idCard.Open()
	if err != nil {
		return model.Kyc{}, err
	}

	portraitBytes, err := io.ReadAll(portraitFile)
	if err != nil {
		return model.Kyc{}, err
	}

	return model.Kyc{
		IdCardImage: file.Image{
			FileName: idCard.Filename,
			Content:  idCardBytes,
		},
		PortraitImage: file.Image{
			FileName: portrait.Filename,
			Content:  portraitBytes,
		},
		UserId: userId,
	}, nil
}

func mapKycModelToDto(res model.Kyc) dto.Kyc {
	var status dto.KycStatus
	approved, rejected := false, false

	if res.ApprovedBy != nil {
		approved = *res.ApprovedBy != uuid.Nil
	}

	if res.RejectedBy != nil {
		rejected = *res.RejectedBy != uuid.Nil
	}

	if approved {
		status = dto.KYCStatusApproved
	} else if rejected {
		status = dto.KYCStatusRejected
	} else {
		status = dto.KYCStatusUndefined
	}

	return dto.Kyc{
		ID:              res.ID,
		IdCardImageUrl:  res.IdCardImage.FileUrl,
		PortraitImage:   res.PortraitImage.FileUrl,
		Status:          status,
		RejectionReason: res.RejectionReason,
	}
}

func createKycListDtoFromModel(kycList []model.Kyc) dto.KycList {
	list := make([]dto.Kyc, 0, len(kycList))

	for _, kyc := range kycList {
		list = append(list, mapKycModelToDto(kyc))
	}

	return dto.KycList{
		KYCList: list,
	}
}

func mapKycModelToEntity(kyc model.Kyc) entity.Kyc {
	return entity.Kyc{
		ID:              kyc.ID,
		ApprovedBy:      kyc.ApprovedBy,
		RejectedBy:      kyc.RejectedBy,
		UserId:          kyc.UserId,
		RejectionReason: &sql.NullString{String: kyc.RejectionReason, Valid: len(kyc.RejectionReason) > 0},
		IdCardImage:     kyc.IdCardImage.FileName,
		PortraitImage:   kyc.PortraitImage.FileName,
	}
}

func mapKycEntityToModel(kyc *entity.Kyc) model.Kyc {

	var rejectionReason string
	if kyc.RejectionReason != nil {
		rejectionReason = kyc.RejectionReason.String
	}

	return model.Kyc{
		ID:              kyc.ID,
		ApprovedBy:      kyc.ApprovedBy,
		RejectedBy:      kyc.RejectedBy,
		UserId:          kyc.UserId,
		RejectionReason: rejectionReason,
		IdCardImage: file.Image{
			FileName: kyc.IdCardImage,
		},
		PortraitImage: file.Image{
			FileName: kyc.PortraitImage,
		},
	}
}

func createModelKycList(kycs *[]entity.Kyc) []model.Kyc {
	var kycList []model.Kyc

	for _, kyc := range *kycs {
		kycList = append(kycList, mapKycEntityToModel(&kyc))
	}

	return kycList
}
