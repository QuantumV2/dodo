package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func DefaultLibDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("LOCALAPPDATA")
		if base == "" {
			return "", fmt.Errorf("LOCALAPPDATA is not set")
		}
		return filepath.Join(base, "dodo", "lib"), nil

	case "linux":
		base := os.Getenv("XDG_DATA_HOME")
		if base == "" {
			base = filepath.Join(home, ".local", "share")
		}
		return filepath.Join(base, "dodo", "lib"), nil

	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "dodo", "lib"), nil

	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
func GetLibrary(name string) ([]byte, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		deflib, deflibErr := DefaultLibDir()
		path := filepath.Join(deflib, name)
		if _, exists := os.Stat(path); exists == nil && deflibErr == nil {
			content, err = os.ReadFile(path)
		}
		return content, err
	}
	return content, err
}
