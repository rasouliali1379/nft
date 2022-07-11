package auth

import (
	"context"
	"maskan/client/jtrace"
	"maskan/contract"
	merrors "maskan/error"
	"maskan/pkg/crypt"
	jwt "maskan/src/jwt/model"
	user "maskan/src/user/model"

	"go.uber.org/fx"
)

type AuthService struct {
	jwtService   contract.IJwtService
	userService  contract.IUserService
	emailService contract.IEmailService
	otpService   contract.IOtpService
}

type AuthServiceParams struct {
	fx.In
	JwtService   contract.IJwtService
	UserService  contract.IUserService
	EmailService contract.IEmailService
	OtpService   contract.IOtpService
}

func NewAuthService(params AuthServiceParams) contract.IAuthService {
	return AuthService{
		jwtService:   params.JwtService,
		userService:  params.UserService,
		emailService: params.EmailService,
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

	addedEmail, err := a.emailService.AddEmail(c, createdUser.ID, userModel.Email)
	if err != nil {
		return "", err
	}

	err = a.emailService.SendOtpEmail(c, addedEmail.ID)
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

	userModel, err := a.userService.GetUser(c, map[string]any{"email": email})
	if err != nil {
		return jwt.Jwt{}, err
	}

	if !crypt.CompareHash(password, userModel.Password) {
		return jwt.Jwt{}, merrors.ErrInvalidCredentials
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

	jwtToken, err := a.jwtService.Generate(c, userId.String())
	if err != nil {
		return jwt.Jwt{}, err
	}

	return jwtToken, nil
}
