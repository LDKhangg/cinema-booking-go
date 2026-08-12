package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadUsesDotEnvValues(t *testing.T) {
	resetConfigEnv(t)
	writeEnvFile(t, "PORT=9090\nDB_DSN=postgres://from-dot-env\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != ":9090" {
		t.Fatalf("expected port %q, got %q", ":9090", cfg.Port)
	}

	if cfg.DatabaseDSN != "postgres://from-dot-env" {
		t.Fatalf("expected db dsn %q, got %q", "postgres://from-dot-env", cfg.DatabaseDSN)
	}
}

func TestLoadPrefersExistingEnvironmentOverDotEnv(t *testing.T) {
	resetConfigEnv(t)
	writeEnvFile(t, "PORT=9090\nDB_DSN=postgres://from-dot-env\n")

	t.Setenv("PORT", "7070")
	t.Setenv("DB_DSN", "postgres://from-env")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != ":7070" {
		t.Fatalf("expected port %q, got %q", ":7070", cfg.Port)
	}

	if cfg.DatabaseDSN != "postgres://from-env" {
		t.Fatalf("expected db dsn %q, got %q", "postgres://from-env", cfg.DatabaseDSN)
	}
}

func TestLoadDefaultsPortWhenNotProvided(t *testing.T) {
	resetConfigEnv(t)
	writeEnvFile(t, "DB_DSN=postgres://from-dot-env\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Port != ":8080" {
		t.Fatalf("expected port %q, got %q", ":8080", cfg.Port)
	}
}

func TestLoadReturnsErrorWhenDatabaseDSNMissing(t *testing.T) {
	resetConfigEnv(t)
	writeEnvFile(t, "PORT=8080\n")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "DB_DSN is required") {
		t.Fatalf("expected DB_DSN required error, got %v", err)
	}
}

func TestNormalizePort(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty uses default", input: "", want: ":8080"},
		{name: "number only gets colon", input: "8080", want: ":8080"},
		{name: "existing colon stays same", input: ":9090", want: ":9090"},
		{name: "trim spaces before normalize", input: " 3000 ", want: ":3000"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizePort(tc.input)
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func resetConfigEnv(t *testing.T) {
	t.Helper()
	t.Setenv("PORT", "")
	t.Setenv("DB_DSN", "")
}

func writeEnvFile(t *testing.T, content string) {
	t.Helper()

	configDir, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("resolve config dir: %v", err)
	}

	envPath := filepath.Join(configDir, ".env")
	originalContent, readErr := os.ReadFile(envPath)
	hadOriginal := readErr == nil
	if readErr != nil && !os.IsNotExist(readErr) {
		t.Fatalf("read existing .env: %v", readErr)
	}

	t.Cleanup(func() {
		if hadOriginal {
			if err := os.WriteFile(envPath, originalContent, 0o644); err != nil {
				t.Fatalf("restore .env: %v", err)
			}
			return
		}

		if err := os.Remove(envPath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("remove test .env: %v", err)
		}
	})

	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
}
