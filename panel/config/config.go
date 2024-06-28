package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Dirs struct {
		Templates string `yaml:"templates"`
		Static    string `yaml:"static"`
	} `yaml:"dirs"`
}

// Read reads a YAML config file and returns a Config struct
func Read(s string) (*Config, error) {
	// Read config file
	f, err := os.ReadFile(s)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	// Unmarshal config file
	var c Config
	if err := yaml.Unmarshal(f, &c); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Get absolute paths for templates and static
	if c.Dirs.Templates, err = filepath.Abs(c.Dirs.Templates); err != nil {
		return nil, fmt.Errorf("getting absolute path for templates: %w", err)
	}
	if c.Dirs.Static, err = filepath.Abs(c.Dirs.Static); err != nil {
		return nil, fmt.Errorf("getting absolute path for static: %w", err)
	}

	return &c, nil
}
