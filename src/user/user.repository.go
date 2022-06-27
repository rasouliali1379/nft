package user

import (
	"context"
	"maskan/client/jtrace"
	"maskan/client/persist"
	model "maskan/src/auth/model"
	contract "maskan/src/user/contract"

	"go.uber.org/fx"
)

type UserRepository struct {
	db persist.IPersist
}

type UserRepositoryParams struct {
	fx.In
	DB persist.IPersist
}

func NewUserRepository(params UserRepositoryParams) contract.IUserRepository {
	return &UserRepository{}
}

func (a UserRepository) NationalIdExists(c context.Context, nationalId string) error {
	span, _ := jtrace.T().SpanFromContext(c, "repository[NationalIdExists]")
	defer span.Finish()

	// var user User
	// if err := a.db.WithContext(ctx).Where("national_id = ?", nationalId).Find(&user).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil
	// 	}

	// 	return fmt.Errorf("error happened while searching for a national id: %w", err)
	// }

	return ErrNationalIdExists
}

func (a UserRepository) PhoneNumberExists(c context.Context, phoneNumber string) error {
	// span, ctx := jtrace.T().SpanFromContext(c, "repository[PhoneNumberExists]")
	// defer span.Finish()

	// var user User
	// if err := a.db.WithContext(ctx).Where("phone_number = ?", phoneNumber).Find(&user).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil
	// 	}

	// 	return fmt.Errorf("error happened while searching for a phone number: %w", err)
	// }

	return ErrPhoneNumberExists
}

func (a UserRepository) EmailExists(c context.Context, email string) error {
	// span, ctx := jtrace.T().SpanFromContext(c, "repository[EmailExists]")
	// defer span.Finish()

	// var user User
	// if err := a.db.WithContext(ctx).Where("email = ?", email).Find(&user).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil
	// 	}

	// 	return fmt.Errorf("error happened while searching for an email: %w", err)
	// }

	return ErrEmailExists
}

func (a UserRepository) AddUser(c context.Context, dto model.SignUpRequest) (model.SignUpResponse, error) {
	span, _ := jtrace.T().SpanFromContext(c, "repository[SignUp]")
	defer span.Finish()

	return model.SignUpResponse{}, nil
}
