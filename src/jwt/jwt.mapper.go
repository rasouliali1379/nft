package jwt

import (
	entity "nft/src/jwt/entity"
	model "nft/src/jwt/model"
)

func mapJwtEntityToRefreshTokenModel(refresh *entity.Jwt) model.RefreshToken {
	return model.RefreshToken{
		Id:        refresh.ID,
		Token:     refresh.Token,
		Invoked:   refresh.Invoked,
		UserId:    refresh.UserId,
		CreatedAt: refresh.CreatedAt,
		UpdatedAt: refresh.UpdatedAt,
	}
}

func mapRefreshTokenModelToJwtEntity(data model.RefreshToken) entity.Jwt {
	return entity.Jwt{
		ID:      data.Id,
		Token:   data.Token,
		Invoked: data.Invoked,
		UserId:  data.UserId,
	}
}
