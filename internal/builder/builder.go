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

// Check validates the Go UI syntax using the AST parser without building.
func Check(root string) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	f := frontend.New()
	_, err = f.BuildIR(filepath.Join(root, cfg.Source))
	if err != nil {
		return fmt.Errorf("syntax check failed: %w", err)
	}
	fmt.Println("Syntax check passed.")
	return nil
}

func prepareDynamicFiles(root string, cfg config.Config) error {
	f := frontend.New()
	prog, err := f.BuildIR(filepath.Join(root, cfg.Source))
	if err != nil {
		return fmt.Errorf("frontend parsing failed: %w", err)
	}
	if prog != nil {
		javaDir := filepath.Join(root, "android", "app", "src", "main", "java", filepath.FromSlash(strings.ReplaceAll(cfg.Package, ".", "/")))
		if err := os.MkdirAll(javaDir, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(javaDir, "Go2ApkActivity.java"), []byte(android.RenderGo2ApkActivity(cfg)), 0o644); err != nil {
			return err
		}
		if len(prog.Pages) > 0 {
			for i, page := range prog.Pages {
				activityName := page.Name
				if i == 0 {
					activityName = "MainActivity"
				}
				if err := os.WriteFile(filepath.Join(javaDir, activityName+".java"), []byte(android.RenderDynamicActivity(cfg, activityName, page.Root)), 0o644); err != nil {
					return err
				}
			}
		} else {
			if err := os.WriteFile(filepath.Join(javaDir, "MainActivity.java"), []byte(android.RenderDynamicActivity(cfg, "MainActivity", prog.UI)), 0o644); err != nil {
				return err
			}
		}
		if err := os.WriteFile(filepath.Join(javaDir, "NativeBridge.java"), []byte(android.RenderNativeBridge(cfg)), 0o644); err != nil {
			return err
		}
		// Generate broadcast receiver class if needed
		if receiverCode := android.RenderBroadcastReceiver(cfg, prog.Receivers); receiverCode != "" {
			if err := os.WriteFile(filepath.Join(javaDir, "Go2APKBroadcastReceiver.java"), []byte(receiverCode), 0o644); err != nil {
				return err
			}
		}
		// Generate VPN service class if needed
		if prog.HasVPN {
			if err := os.WriteFile(filepath.Join(javaDir, "Go2ApkVpnService.java"), []byte(android.RenderVpnService(cfg)), 0o644); err != nil {
				return err
			}
		}
		// Regenerate manifest to inject permissions and broadcast receivers
		manifestDir := filepath.Join(root, "android", "app", "src", "main")
		if err := os.MkdirAll(manifestDir, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(manifestDir, "AndroidManifest.xml"), []byte(android.RenderManifest(cfg, prog)), 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(root, cfg.Source, "events_gen.go"), []byte(android.RenderEventsGen(prog)), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// Build validates the generated project and prepares debug APK outputs.
func Build(root string, opts ...Options) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	applyOptions(&cfg, opts)
	if err := require(filepath.Join(root, "android", "app", "build.gradle")); err != nil {
		return err
	}

	if err := prepareDynamicFiles(root, cfg); err != nil {
		return err
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
		os.WriteFile(filepath.Join(out, "README.txt"), []byte(fmt.Sprintf("Debug build failed: %v\n", err)), 0o644)
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

	if err := prepareDynamicFiles(root, cfg); err != nil {
		return err
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
		os.WriteFile(filepath.Join(out, "README.txt"), []byte(fmt.Sprintf("Release build failed: %v\n", err)), 0o644)
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

// Preview generates an HTML preview of the Android application's UI structure
// without requiring an Android build. It drops a preview.html in the root.
func Preview(root string) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	f := frontend.New()
	prog, err := f.BuildIR(filepath.Join(root, cfg.Source))
	if err != nil {
		return fmt.Errorf("frontend parsing failed: %w", err)
	}
	html := android.RenderPreviewHTML(cfg, prog)
	outPath := filepath.Join(root, "preview.html")
	if err := os.WriteFile(outPath, []byte(html), 0o644); err != nil {
		return fmt.Errorf("failed to write preview.html: %w", err)
	}
	fmt.Printf("Generated UI preview at: %s\nOpen this file in your browser to view the layout.\n", outPath)
	return nil
}
