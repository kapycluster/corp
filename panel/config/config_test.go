package config

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name: "Basic server config",
			envVars: map[string]string{
				"PANEL_SERVER_HOST":    "localhost",
				"PANEL_SERVER_PORT":    "8080",
				"PANEL_SERVER_BASEURL": "http://localhost:8080",
			},
			expected: Config{
				Server: ServerConfig{
					ListenHost: "localhost",
					ListenPort: 8080,
					BaseURL:    "http://localhost:8080",
				},
			},
		},
		{
			name: "OAuth config",
			envVars: map[string]string{
				"PANEL_OAUTH_GITHUB_KEY":    "github-key",
				"PANEL_OAUTH_GITHUB_SECRET": "github-secret",
				"PANEL_OAUTH_GOOGLE_KEY":    "google-key",
				"PANEL_OAUTH_GOOGLE_SECRET": "google-secret",
			},
			expected: Config{
				OAuth: OAuthConfig{
					GitHub: GitHubConfig{
						Key:    "github-key",
						Secret: "github-secret",
					},
					Google: GoogleConfig{
						Key:    "google-key",
						Secret: "google-secret",
					},
				},
			},
		},
		{
			name: "Session config",
			envVars: map[string]string{
				"PANEL_SESSION_SECRET":   "session-secret",
				"PANEL_SESSION_MAXAGE":   "3600",
				"PANEL_SESSION_SECURE":   "true",
				"PANEL_SESSION_HTTPONLY": "true",
			},
			expected: Config{
				Session: SessionConfig{
					Secret:   "session-secret",
					MaxAge:   3600,
					Secure:   true,
					HttpOnly: true,
				},
			},
		},
		{
			name: "DNS config",
			envVars: map[string]string{
				"PANEL_DNS_CLOUDFLARE_APITOKEN": "cf-token",
				"PANEL_DNS_CLOUDFLARE_ZONEID":   "zone-123",
			},
			expected: Config{
				DNS: DNSConfig{
					Cloudflare: CloudflareConfig{
						APIToken: "cf-token",
						ZoneID:   "zone-123",
					},
				},
			},
		},
		{
			name: "Database config",
			envVars: map[string]string{
				"PANEL_DATABASE_URL": "postgres://localhost:5432/db",
			},
			expected: Config{
				Database: DatabaseConfig{
					URL: "postgres://localhost:5432/db",
				},
			},
		},
		{
			name: "Kubernetes config",
			envVars: map[string]string{
				"PANEL_KUBERNETES_KUBECONFIGS": "/path/to/kubeconfig",
			},
			expected: Config{
				Kubernetes: KubernetesConfig{
					KubeconfigsDir: "/path/to/kubeconfig",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg := NewConfig()

			// Server config checks
			if cfg.Server.ListenHost != tt.expected.Server.ListenHost {
				t.Errorf("ListenHost = %v, want %v", cfg.Server.ListenHost, tt.expected.Server.ListenHost)
			}
			if cfg.Server.ListenPort != tt.expected.Server.ListenPort {
				t.Errorf("ListenPort = %v, want %v", cfg.Server.ListenPort, tt.expected.Server.ListenPort)
			}

			// OAuth config checks
			if cfg.OAuth.GitHub.Key != tt.expected.OAuth.GitHub.Key {
				t.Errorf("GitHub Key = %v, want %v", cfg.OAuth.GitHub.Key, tt.expected.OAuth.GitHub.Key)
			}
			if cfg.OAuth.Google.Key != tt.expected.OAuth.Google.Key {
				t.Errorf("Google Key = %v, want %v", cfg.OAuth.Google.Key, tt.expected.OAuth.Google.Key)
			}

			// Session config checks
			if cfg.Session.Secret != tt.expected.Session.Secret {
				t.Errorf("Session Secret = %v, want %v", cfg.Session.Secret, tt.expected.Session.Secret)
			}
			if cfg.Session.MaxAge != tt.expected.Session.MaxAge {
				t.Errorf("Session MaxAge = %v, want %v", cfg.Session.MaxAge, tt.expected.Session.MaxAge)
			}

			// DNS config checks
			if cfg.DNS.Cloudflare.APIToken != tt.expected.DNS.Cloudflare.APIToken {
				t.Errorf("Cloudflare APIToken = %v, want %v", cfg.DNS.Cloudflare.APIToken, tt.expected.DNS.Cloudflare.APIToken)
			}

			// Database config checks
			if cfg.Database.URL != tt.expected.Database.URL {
				t.Errorf("Database URL = %v, want %v", cfg.Database.URL, tt.expected.Database.URL)
			}

			// Kubernetes config checks
			if cfg.Kubernetes.KubeconfigsDir != tt.expected.Kubernetes.KubeconfigsDir {
				t.Errorf("Kubernetes KubeconfigsDir = %v, want %v", cfg.Kubernetes.KubeconfigsDir, tt.expected.Kubernetes.KubeconfigsDir)
			}
		})
	}
}

func TestAsEnv(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple key",
			input:    "host",
			expected: "PANEL_HOST",
		},
		{
			name:     "Nested key",
			input:    "server.host",
			expected: "PANEL_SERVER_HOST",
		},
		{
			name:     "Multiple nested levels",
			input:    "oauth.github.key",
			expected: "PANEL_OAUTH_GITHUB_KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AsEnv(tt.input)
			if result != tt.expected {
				t.Errorf("AsEnv(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
