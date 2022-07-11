package email

import (
	entity "maskan/src/email/entity"
	model "maskan/src/email/model"
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

func createMapFromEmailModel(emailModel model.Email) map[string]any {
	return map[string]any{
		"email":    emailModel.Email,
		"verified": emailModel.Verified,
		"user_id":  emailModel.UserId,
	}
}
