package jwt

import (
	entity "maskan/src/jwt/entity"
	model "maskan/src/jwt/model"
)

func mapJwtEntityToRefreshTokenModel(refresh entity.Jwt) model.RefreshToken {
	return model.RefreshToken{
		Id:        refresh.ID,
		Token:     refresh.Token,
		Invoked:   refresh.Invoked,
		UserId:    refresh.UserId,
		CreatedAt: refresh.CreatedAt,
		UpdatedAt: refresh.UpdatedAt,
	}
}
