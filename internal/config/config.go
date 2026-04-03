package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Proxy         string `yaml:"proxy"`
	Language      string `yaml:"language"`
	Timeout       int    `yaml:"timeout"`
	OutputFormat  string `yaml:"output_format"`
	Delay         int    `yaml:"delay"`
	Cookie        string `yaml:"cookie"`
	CrossRefMailto string `yaml:"crossref_mailto"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	if cfg.Language == "" {
		cfg.Language = "zh-CN"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "bibtex"
	}

	return &cfg, nil
}
