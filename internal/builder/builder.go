package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Baba01hacker666/Go2APK/frontend"
	"github.com/Baba01hacker666/Go2APK/internal/android"
	"github.com/Baba01hacker666/Go2APK/internal/config"
	"github.com/Baba01hacker666/Go2APK/internal/gradle"
)

// Options controls build-time behavior.
type Options struct {
	Obfuscate bool
}

// Build validates the generated project and prepares debug APK outputs.
func Build(root string, opts ...Options) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	applyOptions(&cfg, opts)
	// TODO: Replace Gradle scaffold with internal build pipeline
	if err := require(filepath.Join(root, "android", "app", "build.gradle")); err != nil {
		return err
	}

	// Dynamic AST Transpilation
	f := frontend.New()
	prog, err := f.BuildIR(filepath.Join(root, cfg.Source))
	if err == nil && prog != nil {
		javaDir := filepath.Join(root, "android", "app", "src", "main", "java", filepath.FromSlash(strings.ReplaceAll(cfg.Package, ".", "/")))
		os.MkdirAll(javaDir, 0o755)
		os.WriteFile(filepath.Join(javaDir, "MainActivity.java"), []byte(android.RenderDynamicMainActivity(cfg, prog)), 0o644)
		os.WriteFile(filepath.Join(javaDir, "NativeBridge.java"), []byte(android.RenderNativeBridge(cfg)), 0o644)
	}

	if cfg.Obfuscate {
		if err := writeObfuscatedGradle(root, cfg); err != nil {
			return err
		}
	}
	out := filepath.Join(root, "dist", "debug")
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	if err := runGradle(root, "assembleDebug"); err != nil {
		os.WriteFile(filepath.Join(out, "README.txt"), []byte(fmt.Sprintf("Debug build is configured, but neither gomobile nor Gradle could produce an APK in this environment: %v\n", err)), 0o644)
		return fmt.Errorf("gradle build failed: %w", err)
	}
	return copyArtifacts(filepath.Join(root, "android", "app", "build", "outputs", "apk", "debug"), out)
}

// Release prepares release APK/AAB outputs.
func Release(root string, opts ...Options) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	applyOptions(&cfg, opts)
	if err := require(filepath.Join(root, "android", "app", "build.gradle")); err != nil {
		return err
	}

	// Dynamic AST Transpilation
	f := frontend.New()
	prog, err := f.BuildIR(filepath.Join(root, cfg.Source))
	if err == nil && prog != nil {
		javaDir := filepath.Join(root, "android", "app", "src", "main", "java", filepath.FromSlash(strings.ReplaceAll(cfg.Package, ".", "/")))
		os.MkdirAll(javaDir, 0o755)
		os.WriteFile(filepath.Join(javaDir, "MainActivity.java"), []byte(android.RenderDynamicMainActivity(cfg, prog)), 0o644)
		os.WriteFile(filepath.Join(javaDir, "NativeBridge.java"), []byte(android.RenderNativeBridge(cfg)), 0o644)
	}

	if cfg.Obfuscate {
		if err := writeObfuscatedGradle(root, cfg); err != nil {
			return err
		}
	}
	out := filepath.Join(root, "dist", "release")
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	if err := runGradle(root, "assembleRelease", "bundleRelease"); err != nil {
		os.WriteFile(filepath.Join(out, "README.txt"), []byte(fmt.Sprintf("Release build is configured, but neither gomobile nor Gradle could produce artifacts in this environment: %v\n", err)), 0o644)
		return fmt.Errorf("gradle build failed: %w", err)
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

func applyOptions(cfg *config.Config, opts []Options) {
	for _, opt := range opts {
		if opt.Obfuscate {
			cfg.Obfuscate = true
		}
	}
}

func writeObfuscatedGradle(root string, cfg config.Config) error {
	appDir := filepath.Join(root, "android", "app")
	if err := os.WriteFile(filepath.Join(appDir, "build.gradle"), []byte(android.RenderBuildGradle(cfg)), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(appDir, "proguard-rules.pro"), []byte(android.RenderProguardRules()), 0o644)
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
