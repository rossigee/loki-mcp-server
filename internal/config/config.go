package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	LokiURL       string
	Username      string
	Password      string
	BearerToken   string
	TLSSkipVerify bool
	TenantID      string
	HTTPTimeout   time.Duration
	Transport     string
	Host          string
	Port          string
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPTimeout: 30 * time.Second,
		Transport:   "stdio",
		Host:        "127.0.0.1",
		Port:        "8000",
	}

	cfg.LokiURL = strings.TrimRight(os.Getenv("LOKI_URL"), "/")
	if cfg.LokiURL == "" {
		return nil, fmt.Errorf("LOKI_URL is required")
	}

	u, err := url.Parse(cfg.LokiURL)
	if err != nil {
		return nil, fmt.Errorf("LOKI_URL is not a valid URL: %w", err)
	}
	if !u.IsAbs() || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, fmt.Errorf("LOKI_URL must be an absolute HTTP or HTTPS URL")
	}

	cfg.Username = os.Getenv("LOKI_USERNAME")
	cfg.Password = os.Getenv("LOKI_PASSWORD")
	cfg.BearerToken = os.Getenv("LOKI_BEARER_TOKEN")

	if (cfg.Username != "" || cfg.Password != "") && cfg.BearerToken != "" {
		return nil, fmt.Errorf("LOKI_USERNAME/LOKI_PASSWORD and LOKI_BEARER_TOKEN are mutually exclusive")
	}

	if v := os.Getenv("LOKI_TLS_SKIP_VERIFY"); v == "true" || v == "1" {
		cfg.TLSSkipVerify = true
	}

	cfg.TenantID = os.Getenv("LOKI_TENANT_ID")

	if v := os.Getenv("LOKI_HTTP_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("LOKI_HTTP_TIMEOUT is not a valid duration: %w", err)
		}
		cfg.HTTPTimeout = d
	}

	if v := os.Getenv("TRANSPORT"); v != "" {
		cfg.Transport = v
	}

	if v := os.Getenv("HOST"); v != "" {
		cfg.Host = v
	}

	if v := os.Getenv("PORT"); v != "" {
		cfg.Port = v
	}

	return cfg, nil
}
