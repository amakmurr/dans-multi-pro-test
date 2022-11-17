package internal

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Source string `yaml:"source"`
	} `yaml:"database"`
	Dans DansConfig `yaml:"dans"`
	JWT  JWTConfig  `yaml:"jwt"`
}

type DansConfig struct {
	BaseURL string `yaml:"baseURL"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	d := yaml.NewDecoder(file)
	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
