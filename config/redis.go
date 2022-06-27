package config

// redis struct
type Redis struct {
	Username string `yaml:"redis.username"`
	Password string `yaml:"redis.password"`
	DB       int    `yaml:"redis.db"`
	Host     string `yaml:"redis.host" required:"true"`
	Logger   bool   `yaml:"redis.logger"`
}
