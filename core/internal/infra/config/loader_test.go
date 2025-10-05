package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestSetupWithExplicitPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.yaml")

	if err := os.WriteFile(path, []byte("output:\n  colors: false\n"), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	v := viper.New()
	used, err := Setup(v, path)
	if err != nil {
		t.Fatalf("Setup returned error: %v", err)
	}

	if used != path {
		t.Fatalf("expected explicit path %s, got %s", path, used)
	}

	if v.GetBool("output.colors") {
		t.Fatal("expected colors to be false from explicit config")
	}
}

func TestSetupSearchOrderPrefersHome(t *testing.T) {
	dir := t.TempDir()
	homeDir := filepath.Join(dir, "home")
	if err := os.MkdirAll(homeDir, 0o755); err != nil {
		t.Fatalf("failed to create home dir: %v", err)
	}

	origHomeFn := userHomeDir
	userHomeDir = func() (string, error) { return homeDir, nil }
	t.Cleanup(func() { userHomeDir = origHomeFn })

	homeConfig := filepath.Join(homeDir, ".bip38cli.yaml")
	if err := os.WriteFile(homeConfig, []byte("output:\n  colors: false\n"), 0o600); err != nil {
		t.Fatalf("failed to write home config: %v", err)
	}

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "bip38cli.yaml"), []byte("output:\n  colors: true\n"), 0o600); err != nil {
		t.Fatalf("failed to write local config: %v", err)
	}

	v := viper.New()
	used, err := Setup(v, "")
	if err != nil {
		t.Fatalf("Setup returned error: %v", err)
	}

	if used != homeConfig {
		t.Fatalf("expected home config %s, got %s", homeConfig, used)
	}

	if v.GetBool("output.colors") {
		t.Fatal("expected colors to be false from home config precedence")
	}
}

func TestSetupDefaultsWhenNoConfig(t *testing.T) {
	v := viper.New()
	used, err := Setup(v, "")
	if err != nil {
		t.Fatalf("Setup returned error: %v", err)
	}

	if used != "" {
		t.Fatalf("expected no config path, got %s", used)
	}

	if !v.GetBool("defaults.compressed") {
		t.Fatal("expected defaults.compressed to be true")
	}
	if v.GetString("output.format") != "text" {
		t.Fatalf("expected default output format, got %s", v.GetString("output.format"))
	}
}
