package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var searchOrder = []string{
	".bip38cli.yaml",
	"./bip38cli.yaml",
	"/etc/bip38cli/config.yaml",
}

var userHomeDir = os.UserHomeDir

// Setup configure viper and return config file path when found
func Setup(v *viper.Viper, explicitPath string) (string, error) {
	ApplyDefaults(v)

	if explicitPath != "" {
		v.SetConfigFile(explicitPath)
		if err := readConfig(v); err != nil {
			return "", fmt.Errorf("failed to read config %s: %w", explicitPath, err)
		}
		return explicitPath, nil
	}

	home, err := userHomeDir()
	if err != nil {
		home = ""
	}

	candidates := buildSearchList(home)
	for _, path := range candidates {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}

		v.SetConfigFile(path)
		if err := readConfig(v); err != nil {
			return "", fmt.Errorf("failed to read config %s: %w", path, err)
		}
		return path, nil
	}

	return "", nil
}

func buildSearchList(home string) []string {
	paths := make([]string, 0, len(searchOrder))

	if home != "" {
		paths = append(paths, filepath.Join(home, searchOrder[0]))
	}
	paths = append(paths, searchOrder[1:]...)
	return paths
}

// ApplyDefaults define safe defaults when config absent
func ApplyDefaults(v *viper.Viper) {
	v.SetDefault("defaults.compressed", true)
	v.SetDefault("output.format", "text")
	v.SetDefault("output.colors", true)
}

func readConfig(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return nil
		}
		return err
	}
	return nil
}
