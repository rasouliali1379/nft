package kyc

import (
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	file "nft/src/file/model"
	dto "nft/src/kyc/dto"
	entity "nft/src/kyc/entity"
	model "nft/src/kyc/model"

	"github.com/google/uuid"
)

func createKYCModel(idCard *multipart.FileHeader, portrait *multipart.FileHeader, userId uuid.UUID) (model.KYC, error) {

	idCardFile, err := idCard.Open()
	if err != nil {
		return model.KYC{}, err
	}

	idCardBytes, err := io.ReadAll(idCardFile)
	if err != nil {
		return model.KYC{}, err
	}

	portraitFile, err := idCard.Open()
	if err != nil {
		return model.KYC{}, err
	}

	portraitBytes, err := io.ReadAll(portraitFile)
	if err != nil {
		return model.KYC{}, err
	}

	return model.KYC{
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

func mapKYCModelToDto(res model.KYC) dto.KYC {
	var status dto.KYCStatus
	approved := *res.ApprovedBy != uuid.Nil
	rejected := *res.RejectedBy != uuid.Nil
	log.Println(res.ApprovedBy, approved, rejected)

	if approved {
		status = dto.KYCStatusApproved
	} else if rejected {
		status = dto.KYCStatusRejected
	} else {
		status = dto.KYCStatusUndefined
	}

	return dto.KYC{
		ID:              res.ID,
		IdCardImageUrl:  res.IdCardImage.FileUrl,
		PortraitImage:   res.PortraitImage.FileUrl,
		Status:          status,
		RejectionReason: res.RejectionReason,
	}
}

func createKYCListDtoFromModel(kycList []model.KYC) dto.KYCList {
	list := make([]dto.KYC, 0, len(kycList))

	for _, kyc := range kycList {
		list = append(list, mapKYCModelToDto(kyc))
	}

	return dto.KYCList{
		KYCList: list,
	}
}

func mapKYCModelToEntity(kyc model.KYC) entity.KYC {
	return entity.KYC{
		ID:              kyc.ID,
		ApprovedBy:      kyc.ApprovedBy,
		RejectedBy:      kyc.RejectedBy,
		UserId:          kyc.UserId,
		RejectionReason: &sql.NullString{String: kyc.RejectionReason, Valid: len(kyc.RejectionReason) > 0},
		IdCardImage:     kyc.IdCardImage.FileName,
		PortraitImage:   kyc.PortraitImage.FileName,
	}
}

func mapKYCEntityToModel(kyc *entity.KYC) model.KYC {

	var rejectionReason string
	if kyc.RejectionReason != nil {
		rejectionReason = kyc.RejectionReason.String
	}

	return model.KYC{
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

func createModelKYCList(kycs *[]entity.KYC) []model.KYC {
	var kycList []model.KYC

	for _, kyc := range *kycs {
		kycList = append(kycList, mapKYCEntityToModel(&kyc))

	}

	return kycList
}
