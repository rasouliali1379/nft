package jwt

import (
	"context"
	"errors"
	"fmt"
	"maskan/client/jtrace"
	"maskan/config"
	"maskan/contract"
	jerror "maskan/error"
	jwt "maskan/src/jwt/entity"
	model "maskan/src/jwt/model"
	"time"

	jwtlib "github.com/golang-jwt/jwt"
	"go.uber.org/fx"
)

type JwtRepository struct {
	db contract.IPersist
}

type JwtRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewJwtRepository(params JwtRepositoryParams) contract.IJwtRepository {
	return &JwtRepository{
		db: params.DB,
	}
}

func (j JwtRepository) GenerateToken(c context.Context, userId string, expirationTime time.Time) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[GenerateToken]")
	defer span.Finish()

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.StandardClaims{
		Id:        userId,
		ExpiresAt: expirationTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.C().JWT.HMACSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j JwtRepository) Validate(c context.Context, token string) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[Validate]")
	defer span.Finish()

	parsedToken, err := jwtlib.Parse(token, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, jerror.ErrInvalidSigningMethod
		}

		return []byte(config.C().JWT.HMACSecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwtlib.ValidationError); ok {
			if ve.Errors&jwtlib.ValidationErrorMalformed != 0 {
				return "", jerror.ErrTokenMalformed
			} else if ve.Errors&(jwtlib.ValidationErrorExpired|jwtlib.ValidationErrorNotValidYet) != 0 {
				return "", jerror.ErrTokenExpired
			}
		}
		return "", fmt.Errorf("error happened while parsing token: %w", err)
	}

	if !parsedToken.Valid {
		return "", jerror.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwtlib.StandardClaims)

	if !ok {
		return "", errors.New("error while casting claims")
	}

	return claims.Id, nil
}

func (j JwtRepository) SaveToken(c context.Context, token string, userId string) error {
	span, c := jtrace.T().SpanFromContext(c, "repository[SaveToken]")
	defer span.Finish()

	err := j.db.SaveToken(c, jwt.Jwt{
		Token:  token,
		UserId: userId,
	})

	if err != nil {
		return fmt.Errorf("error happened while saving token in database: %w", err)
	}

	return nil
}

func (j JwtRepository) RetrieveToken(c context.Context, token string) (model.RefreshToken, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[RetrieveToken]")
	defer span.Finish()

	refresh, err := j.db.RetrieveToken(c, token)
	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("error happened while retrieving token from database: %w", err)
	}

	return mapJwtEntityToRefreshTokenModel(refresh), nil
}

func (j JwtRepository) UpdateToken(c context.Context, id uint, token string) error {
	span, c := jtrace.T().SpanFromContext(c, "repository[UpdateToken]")
	defer span.Finish()

	if err := j.db.UpdateToken(c, id, token); err != nil {
		return fmt.Errorf("error happened while updating jwt: %w", err)
	}

	return nil
}
