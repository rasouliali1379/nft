package otp

import (
	entity "maskan/src/otp/entity"
	model "maskan/src/otp/model"
)

func mapOtpEntityToModel(otp *entity.Otp) model.Otp {
	return model.Otp{
		Id:          otp.Id,
		CreatedAt:   otp.CreatedAt,
		Code:        otp.Code,
		UserEmailId: otp.UserEmailId,
	}
}
