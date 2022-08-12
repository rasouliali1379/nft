package config

type File struct {
	TempDir string `yaml:"tempDir" required:"true"`
}