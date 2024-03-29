package jwt

import (
	"context"
	"errors"
	"fmt"
	"nft/config"
	"nft/contract"
	nerror "nft/error"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
	jwt "nft/internal/jwt/entity"
	model "nft/internal/jwt/model"
	"time"

	jwtlib "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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

func (j JwtRepository) Generate(c context.Context, userId string, expirationTime time.Time) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "JwtRepository[Generate]")
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

func (j JwtRepository) Validate(c context.Context, token string) (uuid.UUID, error) {
	span, c := jtrace.T().SpanFromContext(c, "JwtRepository[Validate]")
	defer span.Finish()

	parsedToken, err := jwtlib.Parse(token, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, nerror.ErrInvalidSigningMethod
		}

		return []byte(config.C().JWT.HMACSecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwtlib.ValidationError); ok {
			if ve.Errors&jwtlib.ValidationErrorMalformed != 0 {
				return uuid.UUID{}, nerror.ErrTokenMalformed
			} else if ve.Errors&(jwtlib.ValidationErrorExpired|jwtlib.ValidationErrorNotValidYet) != 0 {
				return uuid.UUID{}, nerror.ErrTokenExpired
			}
		}
		return uuid.UUID{}, fmt.Errorf("error happened while parsing token: %w", err)
	}

	if !parsedToken.Valid {
		return uuid.UUID{}, nerror.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwtlib.MapClaims)
	if !ok {
		return uuid.UUID{}, errors.New("error while casting claims")
	}

	userId, err := uuid.Parse(claims["jti"].(string))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error happened while parsing user id: %w", err)
	}

	return userId, nil
}

func (j JwtRepository) Add(c context.Context, token string, userId string) error {
	span, c := jtrace.T().SpanFromContext(c, "JwtRepository[Add]")
	defer span.Finish()

	_, err := j.db.Create(c, &jwt.Jwt{
		Token:  token,
		UserId: userId,
	})

	if err != nil {
		return fmt.Errorf("error happened while saving token in database: %w", err)
	}

	return nil
}

func (j JwtRepository) Get(c context.Context, conditions persist.D) (model.RefreshToken, error) {
	span, c := jtrace.T().SpanFromContext(c, "JwtRepository[Get]")
	defer span.Finish()

	refresh, err := j.db.Get(c, &jwt.Jwt{}, conditions)
	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("error happened while retrieving token from database: %w", err)
	}

	return mapJwtEntityToRefreshTokenModel(refresh.(*jwt.Jwt)), nil
}

func (j JwtRepository) Update(c context.Context, data model.RefreshToken) error {
	span, c := jtrace.T().SpanFromContext(c, "JwtRepository[Update]")
	defer span.Finish()

	if _, err := j.db.Update(c, &jwt.Jwt{ID: data.Id}, mapRefreshTokenModelToJwtEntity(data)); err != nil {
		return fmt.Errorf("error happened while updating jwt: %w", err)
	}

	return nil
}
