package jwt

import (
	"context"
	"errors"
	"nft/client/jtrace"
	"nft/config"
	"nft/contract"
	merror "nft/error"
	jwt "nft/src/jwt/model"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type JwtService struct {
	jwtRepository contract.IJwtRepository
}

type JwtServiceParams struct {
	fx.In
	JwtRepository contract.IJwtRepository
}

func NewJwtService(params JwtServiceParams) contract.IJwtService {
	return &JwtService{
		jwtRepository: params.JwtRepository,
	}
}

func (j JwtService) Generate(c context.Context, userId string) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "JwtService[Generate]")
	defer span.Finish()

	accessToken, err := j.jwtRepository.Generate(c,
		userId,
		time.Now().Add(time.Duration(time.Minute*time.Duration(config.C().JWT.AccExpInMin))))
	if err != nil {
		return jwt.Jwt{}, err
	}

	refreshToken, err := j.jwtRepository.Generate(c,
		userId,
		time.Now().Add(time.Duration(time.Hour*time.Duration(config.C().JWT.RefExpInHour))))
	if err != nil {
		return jwt.Jwt{}, err
	}

	if err := j.jwtRepository.Add(c, refreshToken, userId); err != nil {
		return jwt.Jwt{}, err
	}

	return jwt.Jwt{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j JwtService) Validate(c context.Context, token string) (uuid.UUID, error) {
	span, c := jtrace.T().SpanFromContext(c, "JwtService[Validate]")
	defer span.Finish()

	userId, err := j.jwtRepository.Validate(c, token)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userId, nil
}

func (j JwtService) Refresh(c context.Context, refreshToken string) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "JwtService[Refresh]")
	defer span.Finish()

	token, err := j.jwtRepository.Get(c, refreshToken)
	if err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return jwt.Jwt{}, merror.ErrTokenInvoked
		}
		return jwt.Jwt{}, err
	}

	exp := time.Hour * time.Duration(config.C().JWT.RefExpInHour)

	if token.UpdatedAt != nil {
		if time.Now().After(token.UpdatedAt.Add(exp)) {
			return jwt.Jwt{}, merror.ErrTokenExpired
		}
	} else {
		if time.Now().After(token.CreatedAt.Add(exp)) {
			return jwt.Jwt{}, merror.ErrTokenExpired
		}
	}

	if token.Invoked {
		return jwt.Jwt{}, merror.ErrTokenInvoked
	}

	accessToken, err := j.jwtRepository.Generate(c,
		token.UserId,
		time.Now().Add(time.Duration(time.Minute*time.Duration(config.C().JWT.AccExpInMin))))
	if err != nil {
		return jwt.Jwt{}, err
	}

	refToken, err := j.jwtRepository.Generate(c,
		token.UserId,
		time.Now().Add(time.Duration(exp)))
	if err != nil {
		return jwt.Jwt{}, err
	}

	if err := j.jwtRepository.Update(c, token.Id, refToken); err != nil {
		return jwt.Jwt{}, err
	}

	return jwt.Jwt{
		AccessToken:  accessToken,
		RefreshToken: refToken,
	}, nil
}

func (j JwtService) GenereteOtpToken(c context.Context, userId string) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "JwtService[GenereteOtpToken]")
	defer span.Finish()

	otpToken, err := j.jwtRepository.Generate(
		c,
		userId,
		time.Now().Add(time.Duration(time.Minute*time.Duration(config.C().Otp.TokenExpInMin))))
	if err != nil {
		return "", err
	}

	return otpToken, nil
}
