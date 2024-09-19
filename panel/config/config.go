package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

const envPrefix = "PANEL_"

type Config struct {
	Server struct {
		ListenHost string `koanf:"host"`
		ListenPort int    `koanf:"port"`
		BaseURL    string `koanf:"baseurl"`
	} `koanf:"server"`
	OAuth struct {
		GitHub struct {
			Key    string `koanf:"key"`
			Secret string `koanf:"secret"`
		} `koanf:"github"`
	} `koanf:"oauth"`
	Session struct {
		Secret   string `koanf:"secret"`
		MaxAge   int    `koanf:"maxage"`
		Secure   bool   `koanf:"secure"`
		HttpOnly bool   `koanf:"httponly"`
	} `koanf:"session"`
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
