package config

type Talan struct {
	BaseUrl      string `yaml:"baseUrl" required:"true"`
	Address      string `yaml:"address" required:"true"`
	Generate     string `yaml:"generate" required:"true"`
	Transactions string `yaml:"transactions" required:"true"`
	Balance      string `yaml:"balance" required:"true"`
}
