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
		OtpCode:   emailRecord.OtpCode,
		UserId:    emailRecord.UserId,
		Verified:  emailRecord.Verified,
	}
}
