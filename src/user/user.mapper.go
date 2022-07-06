package user

import (
	"github.com/google/uuid"
	auth "maskan/src/auth/dto"
	dto "maskan/src/user/dto"
	entity "maskan/src/user/entity"
	model "maskan/src/user/model"
)

func MapSignUpDtoToUserModel(dto auth.SignUpRequest, userId uuid.UUID) model.User {
	return model.User{
		ID:             userId,
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		NationalId:     dto.NationalId,
		Email:          dto.Email,
		PhoneNumber:    dto.PhoneNumber,
		Password:       dto.Password,
		LandLineNumber: dto.LandLineNumber,
		Province:       dto.Province,
		City:           dto.City,
		Address:        dto.Address,
	}
}

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

func createMapFromUserModel(userModel model.User) map[string]any {
	return map[string]any{
		"id":               userModel.ID,
		"national_id":      userModel.NationalId,
		"first_name":       userModel.FirstName,
		"last_name":        userModel.LastName,
		"email":            userModel.Email,
		"phone_number":     userModel.PhoneNumber,
		"land_line_number": userModel.LandLineNumber,
		"province":         userModel.Province,
		"city":             userModel.City,
		"address":          userModel.Address,
	}
}

func mapUserEntityToModel(e *entity.User) model.User {
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

func mapUserModelToResponse(userModel model.User) dto.UserDto {
	return dto.UserDto{
		ID:             userModel.ID.String(),
		NationalId:     userModel.NationalId,
		FirstName:      userModel.FirstName,
		LastName:       userModel.LastName,
		Email:          userModel.Email,
		PhoneNumber:    userModel.PhoneNumber,
		LandLineNumber: userModel.LandLineNumber,
		Province:       userModel.Province,
		City:           userModel.City,
		Address:        userModel.Address,
		PublicKey:      userModel.PublicKey,
	}
}

func createUserList(users []model.User) dto.UserListDto {
	userList := make([]dto.UserDto, len(users))
	for i, userModel := range users {
		userList[i] = dto.UserDto{
			ID:             userModel.ID.String(),
			NationalId:     userModel.NationalId,
			FirstName:      userModel.FirstName,
			LastName:       userModel.LastName,
			Email:          userModel.Email,
			PhoneNumber:    userModel.PhoneNumber,
			LandLineNumber: userModel.LandLineNumber,
			Province:       userModel.Province,
			City:           userModel.City,
			Address:        userModel.Address,
			PublicKey:      userModel.PublicKey,
		}
	}

	return dto.UserListDto{
		Users: userList,
	}
}

func createUserModelList(users *[]entity.User) []model.User {
	var userList []model.User
	for _, userModel := range *users {
		userList = append(userList, model.User{
			ID:             userModel.ID,
			NationalId:     userModel.NationalId,
			FirstName:      userModel.FirstName,
			LastName:       userModel.LastName,
			Email:          userModel.Email,
			PhoneNumber:    userModel.PhoneNumber,
			LandLineNumber: userModel.LandLineNumber,
			Province:       userModel.Province,
			City:           userModel.City,
			Address:        userModel.Address,
			PublicKey:      userModel.PublicKey,
		})
	}

	return userList
}

func mapUserModelToEntity(userModel model.User) entity.User {
	return entity.User{
		ID:             userModel.ID,
		NationalId:     userModel.NationalId,
		FirstName:      userModel.FirstName,
		LastName:       userModel.LastName,
		Email:          userModel.Email,
		PhoneNumber:    userModel.PhoneNumber,
		LandLineNumber: userModel.LandLineNumber,
		Province:       userModel.Province,
		City:           userModel.City,
		Address:        userModel.Address,
	}
}
