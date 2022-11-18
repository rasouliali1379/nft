package jwt

import (
	"errors"
	"nft/contract"
	merror "nft/error"
	"nft/infra/jtrace"
	"nft/pkg/filper"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type JwtMiddleware struct {
	jwtService contract.IJwtService
}

type JwtMiddlewareParams struct {
	fx.In
	JwtService contract.IJwtService
}

func NewJwtMiddleware(params JwtMiddlewareParams) contract.IJwtMiddleware {
	return &JwtMiddleware{
		jwtService: params.JwtService,
	}
}

func (j JwtMiddleware) Handle(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "JwtMiddleware[Handle]")
	defer span.Finish()

	token := c.Get("authorization")
	if token == "" {
		return filper.GetUnAuthError(c, "no authorization token provided")
	}

	splittedToken := strings.Split(token, " ")

	userId, err := j.jwtService.Validate(ctx, splittedToken[1])
	if err != nil {
		if errors.Is(err, merror.ErrInvalidToken) {
			return filper.GetUnAuthError(c, "invalid authorization token")
		} else if errors.Is(err, merror.ErrInvalidSigningMethod) {
			return filper.GetUnAuthError(c, "invalid signing method")
		} else if errors.Is(err, merror.ErrTokenMalformed) {
			return filper.GetUnAuthError(c, "token malformed")
		} else if errors.Is(err, merror.ErrTokenExpired) {
			return filper.GetUnAuthError(c, "token expired")
		} else if errors.Is(err, merror.ErrTokenInvoked) {
			return filper.GetUnAuthError(c, "token invoked")
		}
		return filper.GetInternalError(c, "")
	}

	c.Locals("user_id", userId)

	return c.Next()
}
