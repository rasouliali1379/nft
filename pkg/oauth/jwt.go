package oauth

import (
	"context"
	"errors"
	"fmt"
	"maskan/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	JWT oautInterface = &AccessDetails{}
)

type oautInterface interface {
	VerifyToken(ctx context.Context, request string) (*AccessDetails, error)
}

// AccessDetails struct
type AccessDetails struct {
	UserID string
	Dialer string
}

func (a *AccessDetails) VerifyToken(ctx context.Context, request string) (*AccessDetails, error) {
	token, err := jwt.Parse(request, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(config.C().JWT.HMACSecret), nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return jwt.ParseRSAPublicKeyFromPEM([]byte(config.C().JWT.RSASecret))
		}
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			return jwt.ParseECPublicKeyFromPEM([]byte(config.C().JWT.ECDSASecret))
		}

		zap.L().Error(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		return nil, errors.New("new errors")
	})
	if err != nil {
		return nil, errors.New("new errors")
	}
	if token.Claims == nil || !token.Valid {
		return nil, errors.New("new errors")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			zap.L().Error(fmt.Sprintf("error in get user sub from token: %v", claims))
			return nil, errors.New("new errors")
		}

		if _, err := uuid.Parse(userID); err != nil {
			zap.L().Error(fmt.Sprintf("error in sub typo: %v", err))
			return nil, errors.New("new errors")
		}

		dialer, ok := claims["mobile"].(string)
		if !ok {
			zap.L().Error(fmt.Sprintf("error in get user mobile from token: %v", claims))
			return nil, errors.New("new errors")
		}

		return &AccessDetails{
			UserID: userID,
			Dialer: dialer,
		}, nil
	}
	return nil, err
}
