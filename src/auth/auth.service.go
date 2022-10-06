package auth

import (
	"context"
	"errors"
	"nft/client/jtrace"
	"nft/contract"
	nerror "nft/error"
	"nft/pkg/crypt"
	jwt "nft/src/jwt/model"
	user "nft/src/user/model"

	"go.uber.org/fx"
)

type AuthService struct {
	emailService contract.IEmailService
	jwtService   contract.IJwtService
	userService  contract.IUserService
	otpService   contract.IOtpService
}

type AuthServiceParams struct {
	fx.In
	EmailService contract.IEmailService
	JwtService   contract.IJwtService
	UserService  contract.IUserService
	OtpService   contract.IOtpService
}

func NewAuthService(params AuthServiceParams) contract.IAuthService {
	return &AuthService{
		emailService: params.EmailService,
		jwtService:   params.JwtService,
		userService:  params.UserService,
		otpService:   params.OtpService,
	}
}

func (a AuthService) SignUp(c context.Context, userModel user.User) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "AuthService[SignUp]")
	defer span.Finish()

	createdUser, err := a.userService.AddUser(c, userModel)
	if err != nil {
		return "", err
	}

	userEmail, err := a.emailService.GetUserEmail(c, createdUser.ID)
	if err != nil {
		return "", err
	}

	err = a.emailService.SendOtpEmail(c, userEmail.ID)
	if err != nil {
		return "", err
	}

	token, err := a.jwtService.GenereteOtpToken(c, createdUser.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a AuthService) Login(c context.Context, email string, password string) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "AuthService[Login]")
	defer span.Finish()

	userEmail, err := a.emailService.GetEmail(c, email)
	if err != nil {
		if errors.Is(err, nerror.ErrEmailNotFound) {
			return jwt.Jwt{}, nerror.ErrInvalidCredentials
		}
		return jwt.Jwt{}, err
	}

	userModel, err := a.userService.GetUser(c, map[string]any{"id": userEmail.UserId})
	if err != nil {
		return jwt.Jwt{}, err
	}

	if !crypt.CompareHash(password, userModel.Password) {
		return jwt.Jwt{}, nerror.ErrInvalidCredentials
	}

	token, err := a.jwtService.Generate(c, userModel.ID.String())
	if err != nil {
		return jwt.Jwt{}, err
	}

	return token, nil
}

func (a AuthService) VerifyEmail(c context.Context, token string, code string) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "AuthService[VerifyEmail]")
	defer span.Finish()

	userId, err := a.jwtService.Validate(c, token)
	if err != nil {
		return jwt.Jwt{}, err
	}

	emailModel, err := a.emailService.GetUserEmail(c, userId)
	if err != nil {
		return jwt.Jwt{}, err
	}

	if err = a.otpService.ValidateCode(c, code, emailModel.ID); err != nil {
		return jwt.Jwt{}, err
	}

	if err := a.emailService.ApproveEmail(c, userId, emailModel.Email); err != nil {
		return jwt.Jwt{}, err
	}

	jwtToken, err := a.jwtService.Generate(c, userId.String())
	if err != nil {
		return jwt.Jwt{}, err
	}

	return jwtToken, nil
}

func (a AuthService) ResendVerificationEmail(c context.Context, token string) (string, error) {
	span, c := jtrace.T().SpanFromContext(c, "AuthService[ResendVerificationEmail]")
	defer span.Finish()

	userId, err := a.jwtService.Validate(c, token)
	if err != nil {
		return "", err
	}

	emailModel, err := a.emailService.GetUserEmail(c, userId)
	if err != nil {
		return "", err
	}

	err = a.emailService.SendOtpEmail(c, emailModel.ID)
	if err != nil {
		return "", err
	}

	return token, err
}
