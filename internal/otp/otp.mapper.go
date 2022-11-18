package otp

import (
	entity "nft/internal/otp/entity"
	model "nft/internal/otp/model"
)

func mapOtpEntityToModel(otp *entity.Otp) model.Otp {
	return model.Otp{
		Id:          otp.Id,
		CreatedAt:   otp.CreatedAt,
		Code:        otp.Code,
		UserEmailId: otp.UserEmailId,
	}
}
