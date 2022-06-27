package config

type Logstash struct {
	Endpoint string `yaml:"logstash.endpoint" required:"true"`
	Timeout  int    `yaml:"logstash.timeout"`
}
