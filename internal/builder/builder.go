package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go2apk/go2apk/internal/config"
	"github.com/go2apk/go2apk/internal/gomobile"
	"github.com/go2apk/go2apk/internal/gradle"
)

// Build validates the generated project and prepares debug APK outputs.
func Build(root string) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	if _, err := gomobile.BuildAPK(gomobile.Options{Root: root, Config: cfg}); err == nil {
		return nil
	} else {
		fmt.Fprintf(os.Stderr, "go2apk: gomobile build unavailable, falling back to Gradle scaffold: %v\n", err)
	}
	if err := require(filepath.Join(root, "android", "app", "build.gradle")); err != nil {
		return err
	}
	out := filepath.Join(root, "dist", "debug")
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	if err := runGradle(root, "assembleDebug"); err != nil {
		return os.WriteFile(filepath.Join(out, "README.txt"), []byte(fmt.Sprintf("Debug build is configured, but neither gomobile nor Gradle could produce an APK in this environment: %v\n", err)), 0o644)
	}
	return copyArtifacts(filepath.Join(root, "android", "app", "build", "outputs", "apk", "debug"), out)
}

// Release prepares release APK/AAB outputs.
func Release(root string) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	if _, err := gomobile.BuildAPK(gomobile.Options{Root: root, Config: cfg, Release: true}); err == nil {
		return nil
	} else {
		fmt.Fprintf(os.Stderr, "go2apk: gomobile release unavailable, falling back to Gradle scaffold: %v\n", err)
	}
	out := filepath.Join(root, "dist", "release")
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	if err := runGradle(root, "assembleRelease", "bundleRelease"); err != nil {
		return os.WriteFile(filepath.Join(out, "README.txt"), []byte(fmt.Sprintf("Release build is configured, but neither gomobile nor Gradle could produce artifacts in this environment: %v\n", err)), 0o644)
	}
	if err := copyArtifacts(filepath.Join(root, "android", "app", "build", "outputs", "apk", "release"), out); err != nil {
		return err
	}
	return copyArtifacts(filepath.Join(root, "android", "app", "build", "outputs", "bundle", "release"), out)
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

func runGradle(root string, tasks ...string) error {
	cmdName, args := gradle.Command(root)
	args = append(args, tasks...)
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = filepath.Join(root, "android")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyArtifacts(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		in, err := os.ReadFile(filepath.Join(srcDir, entry.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dstDir, entry.Name()), in, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func loadConfig(root string) (config.Config, error) {
	path := filepath.Join(root, "go2apk.yaml")
	if err := require(path); err != nil {
		return config.Config{}, err
	}
	return config.Load(path)
}
