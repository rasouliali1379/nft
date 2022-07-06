package contract

type IOtpRepository interface {
	Generate(index int) string
	Validate(code string, index int) bool
}

type IOtpService interface {
}
