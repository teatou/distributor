package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Balancer BalancerConfig `yaml:"balancer"`
	Cluster  ClusterConfig  `yaml:"cluster"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type BalancerConfig struct {
	Port int `yaml:"port"`
}

type ClusterConfig struct {
	Ports []int `yaml:"ports"`
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
