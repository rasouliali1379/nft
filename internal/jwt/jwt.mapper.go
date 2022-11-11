package jwt

import (
	entity "nft/internal/jwt/entity"
	model "nft/internal/jwt/model"
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
