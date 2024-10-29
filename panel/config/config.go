package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

const envPrefix = "PANEL_"

type ServerConfig struct {
	ListenHost string `koanf:"host"`
	ListenPort int    `koanf:"port"`
	BaseURL    string `koanf:"baseurl"`
}

type GitHubConfig struct {
	Key    string `koanf:"key"`
	Secret string `koanf:"secret"`
}

type GoogleConfig struct {
	Key    string `koanf:"key"`
	Secret string `koanf:"secret"`
}

type OAuthConfig struct {
	GitHub GitHubConfig `koanf:"github"`
	Google GoogleConfig `koanf:"google"`
}

type SessionConfig struct {
	Secret   string `koanf:"secret"`
	MaxAge   int    `koanf:"maxage"`
	Secure   bool   `koanf:"secure"`
	HttpOnly bool   `koanf:"httponly"`
}

type Config struct {
	Server  ServerConfig  `koanf:"server"`
	OAuth   OAuthConfig   `koanf:"oauth"`
	Session SessionConfig `koanf:"session"`
}

func NewConfig() *Config {
	k := koanf.New(".")
	k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)

	// Setup additional koanf providers here.

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return &cfg
}
