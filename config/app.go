package config

// Service details
type App struct {
	Name      string `yaml:"app.name" required:"true"`
	BaseURL   string `yaml:"app.baseURL"`
	DebugMode bool   `yaml:"app.debugMode"`
	Http      Http   `yaml:"app.http"`
}

type Http struct {
	Cors string `yaml:"app.http.cors" required:"true"`
	Port string `yaml:"app.http.port" required:"true"`
	Host  string `yaml:"app.http.host" required:"true"`
}
