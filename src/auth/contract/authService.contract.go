package auth

import (
	"context"
	model "maskan/src/auth/model"
)

type IAuthService interface {
	SignUp(context.Context, model.SignUpRequest) (model.SignUpResponse, error)
}
