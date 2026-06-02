package config

import "testing"

func TestFromEnvLoadsArrInstances(t *testing.T) {
	t.Setenv("SONARR_URL", "http://sonarr.local:8989/")
	t.Setenv("SONARR_API_KEY", "sonarr-key")
	t.Setenv("RADARR_URL", "http://radarr.local:7878")
	t.Setenv("RADARR_API_KEY", "radarr-key")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv() error = %v", err)
	}

	if cfg.Sonarr.URL != "http://sonarr.local:8989" {
		t.Fatalf("Sonarr.URL = %q", cfg.Sonarr.URL)
	}
	if cfg.Sonarr.APIKey != "sonarr-key" {
		t.Fatalf("Sonarr.APIKey = %q", cfg.Sonarr.APIKey)
	}
	if cfg.Radarr.URL != "http://radarr.local:7878" {
		t.Fatalf("Radarr.URL = %q", cfg.Radarr.URL)
	}
	if cfg.Radarr.APIKey != "radarr-key" {
		t.Fatalf("Radarr.APIKey = %q", cfg.Radarr.APIKey)
	}
	if cfg.AllowMutations {
		t.Fatal("AllowMutations defaulted to true")
	}
}

func TestFromEnvRequiresAtLeastOneConfiguredService(t *testing.T) {
	_, err := FromEnv()
	if err == nil {
		t.Fatal("FromEnv() error = nil")
	}
}

func TestFromEnvRejectsPartialServiceConfig(t *testing.T) {
	t.Setenv("SONARR_URL", "http://sonarr.local:8989")

	_, err := FromEnv()
	if err == nil {
		t.Fatal("FromEnv() error = nil")
	}
}
