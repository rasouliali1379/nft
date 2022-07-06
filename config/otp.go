package config

type Otp struct {
	Secret string `yaml:"otp.secret" required:"true"`
}
