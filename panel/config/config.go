package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

const envPrefix = "PANEL_"

type ServerConfig struct {
	// ListenHost is the host address on which the server listens.
	ListenHost string `koanf:"host"`
	// ListenPort is the port on which the server listens.
	ListenPort int `koanf:"port"`
	// BaseURL is the base URL for the server.
	BaseURL string `koanf:"baseurl"`
	// PullToken is the token used to pull images from the registry.
	// For now, this is the kapyserver image.
	PullToken string `koanf:"pulltoken"`

	// ControlPlaneBaseURL is the base URL for the control plane.
	ControlPlaneBaseURL string `koanf:"controlplanebaseurl"`

	// LocalDev indicates if the server is running in local development mode.
	LocalDev bool `koanf:"localdev"`
}

type GitHubConfig struct {
	// Key is the GitHub OAuth key.
	Key string `koanf:"key"`
	// Secret is the GitHub OAuth secret.
	Secret string `koanf:"secret"`
}

type GoogleConfig struct {
	// Key is the Google OAuth key.
	Key string `koanf:"key"`
	// Secret is the Google OAuth secret.
	Secret string `koanf:"secret"`
}

type OAuthConfig struct {
	// GitHub contains the GitHub OAuth configuration.
	GitHub GitHubConfig `koanf:"github"`
	// Google contains the Google OAuth configuration.
	Google GoogleConfig `koanf:"google"`
}

type SessionConfig struct {
	// Secret is the secret key used for session encryption.
	Secret string `koanf:"secret"`
	// MaxAge is the maximum age of the session in seconds.
	MaxAge int `koanf:"maxage"`
	// Secure indicates if the session cookie should be secure.
	Secure bool `koanf:"secure"`
	// HttpOnly indicates if the session cookie should be HTTP only.
	HttpOnly bool `koanf:"httponly"`
}

type CloudflareConfig struct {
	// APIToken is the Cloudflare API token.
	APIToken string `koanf:"apitoken"`
	// ZoneID is the Cloudflare zone ID.
	ZoneID string `koanf:"zoneid"`
}

type DNSConfig struct {
	// Cloudflare contains the Cloudflare DNS configuration.
	Cloudflare CloudflareConfig `koanf:"cloudflare"`
}

type DatabaseConfig struct {
	URL string `koanf:"url"`
}

type KubernetesConfig struct {
	KubeconfigsDir string `koanf:"kubeconfigs"`
}

type Config struct {
	// Server contains the server configuration.
	Server ServerConfig `koanf:"server"`
	// OAuth contains the OAuth configuration.
	OAuth OAuthConfig `koanf:"oauth"`
	// Session contains the session configuration.
	Session SessionConfig `koanf:"session"`
	// DNS contains the DNS configuration.
	DNS DNSConfig `koanf:"dns"`
	// Database contains the database configuration.
	Database DatabaseConfig `koanf:"database"`
	// Kubernetes contains Kubernetes-specific configuration.
	Kubernetes KubernetesConfig `koanf:"kubernetes"`

	// Koanf is the underlying koanf instance.
	Koanf *koanf.Koanf
}

func NewConfig() *Config {
	k := koanf.New(".")
	k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)

	// TODO: check if values are set and fail if not
	// Setup additional koanf providers here.

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	cfg.Koanf = k

	return &cfg
}

func AsEnv(key string) string {
	return envPrefix + strings.ToUpper(strings.Replace(key, ".", "_", -1))
}
