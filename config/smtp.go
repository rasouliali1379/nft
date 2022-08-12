package config

type Smtp struct {
	Host     string `yaml:"smtp.host"`
	Port     string `yaml:"smtp.port"`
	From     string `yaml:"smtp.from"`
	Password string `yaml:"smtp.password"`
}
