package user

import (
	auth "maskan/src/auth/model"
	entity "maskan/src/user/entity"
	model "maskan/src/user/model"
)

func mapSignUpRequestModelToEntity(dto auth.SignUpRequest) entity.User {
	return entity.User{
		NationalId:     dto.NationalId,
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		Email:          dto.Email,
		PhoneNumber:    dto.PhoneNumber,
		LandLineNumber: dto.LandLineNumber,
		Province:       dto.Province,
		City:           dto.City,
		Address:        dto.Address,
	}
}

func mapUserEntityToModel(e entity.User) model.User {
	return model.User{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,

		NationalId:     e.NationalId,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		Email:          e.Email,
		PhoneNumber:    e.PhoneNumber,
		LandLineNumber: e.LandLineNumber,
		Province:       e.Province,
		City:           e.City,
		Address:        e.Address,

		Password:   e.Password,
		PublicKey:  e.PublicKey,
		PrivateKey: e.PrivateKey,
	}
}
