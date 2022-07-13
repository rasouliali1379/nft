package email

import (
	entity "nft/src/email/entity"
	model "nft/src/email/model"
)

func mapEmailEntityToModel(emailRecord *entity.Email) model.Email {
	return model.Email{
		ID:        emailRecord.ID,
		CreatedAt: emailRecord.CreatedAt,
		UpdatedAt: emailRecord.UpdatedAt,
		Email:     emailRecord.Email,
		UserId:    emailRecord.UserId,
		Verified:  emailRecord.Verified,
	}
}

func createMapFromEmailModel(emailModel model.Email) entity.Email {
	return entity.Email{
		Email:    emailModel.Email,
		Verified: emailModel.Verified,
		UserId:   emailModel.UserId,
	}
}
