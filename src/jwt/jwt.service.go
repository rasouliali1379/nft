package jwt

import (
	"context"
	"log"
	"maskan/client/jtrace"
	"maskan/contract"
	jerror "maskan/error"
	jwt "maskan/src/jwt/model"
	"time"

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
	span, c := jtrace.T().SpanFromContext(c, "service[Generate]")
	defer span.Finish()

	accessToken, err := j.jwtRepository.GenerateToken(c, userId, time.Now().Add(time.Duration(time.Minute*30)))
	if err != nil {
		log.Println(err)
		return jwt.Jwt{}, err
	}

	refreshToken, err := j.jwtRepository.GenerateToken(c, userId, time.Now().Add(time.Duration(time.Hour*720)))
	if err != nil {
		log.Println(err)
		return jwt.Jwt{}, err
	}

	if err := j.jwtRepository.SaveToken(c, refreshToken, userId); err != nil {
		log.Println(err)
		return jwt.Jwt{}, err
	}

	return jwt.Jwt{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j JwtService) Validate(c context.Context, accessToken string) error {
	span, c := jtrace.T().SpanFromContext(c, "service[Validate]")
	defer span.Finish()

	_, err := j.jwtRepository.Validate(c, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (j JwtService) Refresh(c context.Context, refreshToken string) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[Refresh]")
	defer span.Finish()

	token, err := j.jwtRepository.RetrieveToken(c, refreshToken)
	if err != nil {
		return jwt.Jwt{}, err
	}

	if token.Invoked {
		return jwt.Jwt{}, jerror.ErrTokenInvoked
	}

	accessToken, err := j.jwtRepository.GenerateToken(c, token.UserId, time.Now().Add(time.Duration(time.Minute*30)))
	if err != nil {
		return jwt.Jwt{}, err
	}

	refToken, err := j.jwtRepository.GenerateToken(c, token.UserId, time.Now().Add(time.Duration(time.Hour*720)))
	if err != nil {
		log.Println(err)
		return jwt.Jwt{}, err
	}

	if err := j.jwtRepository.UpdateToken(c, token.Id, refToken); err != nil {
		return jwt.Jwt{}, err
	}

	return jwt.Jwt{
		AccessToken:  accessToken,
		RefreshToken: refToken,
	}, nil
}
