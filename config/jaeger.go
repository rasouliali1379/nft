package config

// Jaeger tracer
type Jaeger struct {
	HostPort string `yaml:"jaeger.hostPort" required:"true"`
	LogSpans bool   `yaml:"jaeger.logSpans"`
}
