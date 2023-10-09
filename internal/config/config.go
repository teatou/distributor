package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger LoggerConfig `yaml:"logger"`
	Dist   DistConfig   `yaml:"dist"`
	Target TargetConfig `yaml:"target"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type DistConfig struct {
}

type TargetConfig struct {
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	cfg := Config{}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file: %w", err)
	}
	return &cfg, nil
}
