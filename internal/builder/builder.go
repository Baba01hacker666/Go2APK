package builder

import (
	"fmt"
	"os"
	"path/filepath"
)

// Build validates the generated project and prepares the debug output directory.
func Build(root string) error {
	if err := require(filepath.Join(root, "go2apk.yaml")); err != nil {
		return err
	}
	out := filepath.Join(root, "dist", "debug")
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(out, "README.txt"), []byte("Debug APK artifacts will be written here once Android SDK integration is configured.\n"), 0o644)
}

// Release prepares the release output directory.
func Release(root string) error {
	if err := require(filepath.Join(root, "go2apk.yaml")); err != nil {
		return err
	}
	out := filepath.Join(root, "dist", "release")
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(out, "README.txt"), []byte("Release APK, AAB, mapping files, and checksums will be written here.\n"), 0o644)
}

// Clean removes generated build artifacts.
func Clean(root string) error {
	return os.RemoveAll(filepath.Join(root, "dist"))
}

func require(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("required file %s is missing; run go2apk init first", path)
	}
	return nil
}
