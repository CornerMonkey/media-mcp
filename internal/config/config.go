package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Sonarr         Service
	Radarr         Service
	AllowMutations bool
}

type Service struct {
	URL    string
	APIKey string
}

func FromEnv() (Config, error) {
	cfg := Config{
		Sonarr: Service{
			URL:    normalizeURL(os.Getenv("SONARR_URL")),
			APIKey: os.Getenv("SONARR_API_KEY"),
		},
		Radarr: Service{
			URL:    normalizeURL(os.Getenv("RADARR_URL")),
			APIKey: os.Getenv("RADARR_API_KEY"),
		},
		AllowMutations: strings.EqualFold(os.Getenv("MEDIA_MCP_ALLOW_MUTATIONS"), "true"),
	}

	if err := validateService("Sonarr", cfg.Sonarr); err != nil {
		return Config{}, err
	}
	if err := validateService("Radarr", cfg.Radarr); err != nil {
		return Config{}, err
	}
	if !cfg.Sonarr.Configured() && !cfg.Radarr.Configured() {
		return Config{}, errors.New("configure at least one of Sonarr or Radarr")
	}

	return cfg, nil
}

func (s Service) Configured() bool {
	return s.URL != "" && s.APIKey != ""
}

func validateService(name string, service Service) error {
	if service.URL == "" && service.APIKey == "" {
		return nil
	}
	if service.URL == "" || service.APIKey == "" {
		return fmt.Errorf("%s requires both URL and API key", name)
	}
	return nil
}

func normalizeURL(value string) string {
	return strings.TrimRight(strings.TrimSpace(value), "/")
}
