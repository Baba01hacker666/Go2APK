// Package builder turns a Go2APK project into Android sources and, when the
// toolchain is available, into APK artifacts.
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
	"github.com/Baba01hacker666/Go2APK/ir"
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
	if _, err := parseProgram(root, cfg); err != nil {
		return fmt.Errorf("syntax check failed: %w", err)
	}
	fmt.Println("Syntax check passed.")
	return nil
}

// Generate regenerates the Android project from the app's Go source without
// invoking Gradle.
func Generate(root string) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	if _, err := prepareDynamicFiles(root, cfg); err != nil {
		return err
	}
	fmt.Println("Generated Android sources.")
	return nil
}

// Build validates the generated project and produces a debug APK.
func Build(root string, opts ...Options) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	applyOptions(&cfg, opts)
	if err := require(filepath.Join(root, "android", "app", "build.gradle")); err != nil {
		return err
	}
	if _, err := prepareDynamicFiles(root, cfg); err != nil {
		return err
	}

	out := filepath.Join(root, "dist", "debug")
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	if err := runGradle(root, "assembleDebug"); err != nil {
		return fmt.Errorf("debug build failed: %w", err)
	}
	return copyArtifacts(filepath.Join(root, "android", "app", "build", "outputs", "apk", "debug"), out)
}

// Release produces release APK and AAB artifacts.
func Release(root string, opts ...Options) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	applyOptions(&cfg, opts)
	if err := require(filepath.Join(root, "android", "app", "build.gradle")); err != nil {
		return err
	}
	if _, err := prepareDynamicFiles(root, cfg); err != nil {
		return err
	}

	if !cfg.Signed() {
		fmt.Fprintln(os.Stderr, "warning: no keystore configured; the release APK will be signed with "+
			"the debug key and cannot be published")
	}

	out := filepath.Join(root, "dist", "release")
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	if err := runGradle(root, "assembleRelease", "bundleRelease"); err != nil {
		return fmt.Errorf("release build failed: %w", err)
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

// Preview renders an HTML approximation of the app's UI so a layout can be
// checked in a browser without the Android SDK.
func Preview(root string) error {
	cfg, err := loadConfig(root)
	if err != nil {
		return err
	}
	prog, err := parseProgram(root, cfg)
	if err != nil {
		return fmt.Errorf("frontend parsing failed: %w", err)
	}
	outPath := filepath.Join(root, "preview.html")
	if err := os.WriteFile(outPath, []byte(android.RenderPreviewHTML(cfg, prog)), 0o644); err != nil {
		return fmt.Errorf("failed to write preview.html: %w", err)
	}
	fmt.Printf("Wrote UI preview to %s\nOpen it in a browser to check the layout.\n", outPath)
	return nil
}

// parseProgram reads the app's Go source into the IR.
func parseProgram(root string, cfg config.Config) (*ir.Program, error) {
	prog, err := frontend.New().BuildIRDirs(resolveSources(root, cfg))
	if err != nil {
		return nil, err
	}
	ReportWarnings(os.Stderr, LintProgram(prog))
	return prog, nil
}

// prepareDynamicFiles regenerates every file that is derived from the app's
// source: the Java activities, the manifest, the Gradle script and the two Go
// files that wire up JNI callbacks.
//
// It returns the extra Android dependencies the generated code needs.
func prepareDynamicFiles(root string, cfg config.Config) ([]string, error) {
	prog, err := parseProgram(root, cfg)
	if err != nil {
		return nil, fmt.Errorf("frontend parsing failed: %w", err)
	}

	extraDeps := android.RequiredDependencies(prog)

	javaDir := filepath.Join(root, "android", "app", "src", "main", "java", filepath.FromSlash(strings.ReplaceAll(cfg.Package, ".", "/")))
	if err := os.MkdirAll(javaDir, 0o755); err != nil {
		return nil, err
	}
	if err := write(filepath.Join(javaDir, "Go2ApkActivity.java"), android.RenderGo2ApkActivity(cfg)); err != nil {
		return nil, err
	}
	if err := write(filepath.Join(javaDir, "NativeBridge.java"), android.RenderNativeBridge(cfg)); err != nil {
		return nil, err
	}

	switch {
	case len(prog.Pages) > 0:
		for i, page := range prog.Pages {
			activityName := page.Name
			if i == 0 {
				activityName = "MainActivity"
			}
			body := android.RenderDynamicActivity(cfg, activityName, page.Root)
			if err := write(filepath.Join(javaDir, activityName+".java"), body); err != nil {
				return nil, err
			}
		}
	default:
		body := android.RenderDynamicActivity(cfg, "MainActivity", prog.UI)
		if err := write(filepath.Join(javaDir, "MainActivity.java"), body); err != nil {
			return nil, err
		}
	}

	// Optional platform classes, only emitted when the app declares them.
	if code := android.RenderBroadcastReceiver(cfg, prog.Receivers); code != "" {
		if err := write(filepath.Join(javaDir, "Go2APKBroadcastReceiver.java"), code); err != nil {
			return nil, err
		}
	}
	if prog.HasVPN {
		if err := write(filepath.Join(javaDir, "Go2ApkVpnService.java"), android.RenderVpnService(cfg)); err != nil {
			return nil, err
		}
	}

	manifestDir := filepath.Join(root, "android", "app", "src", "main")
	if err := os.MkdirAll(manifestDir, 0o755); err != nil {
		return nil, err
	}
	if err := write(filepath.Join(manifestDir, "AndroidManifest.xml"), android.RenderManifest(cfg, prog)); err != nil {
		return nil, err
	}

	// The Gradle script is derived too: it has to carry the dependencies the
	// generated Java needs, so regenerating it keeps the two in step.
	if err := write(filepath.Join(root, "android", "app", "build.gradle"), android.RenderBuildGradle(cfg, extraDeps...)); err != nil {
		return nil, err
	}
	if err := write(filepath.Join(root, "android", "app", "proguard-rules.pro"), android.RenderProguardRules()); err != nil {
		return nil, err
	}

	if err := writeGeneratedGo(filepath.Join(root, cfg.Source, "events_gen.go"),
		len(prog.Events) > 0, android.RenderEventsGen(prog)); err != nil {
		return nil, err
	}
	if err := writeGeneratedGo(filepath.Join(root, cfg.Source, "jni_exports_gen.go"),
		true, android.RenderJNIExports(cfg.Package, prog.AndroidImport)); err != nil {
		return nil, err
	}

	// Copy user assets and resources into the Android project.
	if err := copyIfDir(filepath.Join(root, cfg.Source, "assets"), filepath.Join(root, "android", "app", "src", "main", "assets")); err != nil {
		return nil, fmt.Errorf("failed to copy assets: %w", err)
	}
	if err := copyIfDir(filepath.Join(root, cfg.Source, "res"), filepath.Join(root, "android", "app", "src", "main", "res")); err != nil {
		return nil, fmt.Errorf("failed to copy res: %w", err)
	}
	return extraDeps, nil
}

// writeGeneratedGo writes a Go file the tool owns. When the file is no longer
// needed it is removed, but only if it still carries the generated marker, so a
// file the user wrote is never deleted.
func writeGeneratedGo(path string, needed bool, contents string) error {
	if !needed {
		if existing, err := os.ReadFile(path); err == nil && strings.Contains(string(existing), android.GeneratedFileHeader) {
			return os.Remove(path)
		}
		return nil
	}
	return write(path, contents)
}

func write(path string, contents string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(contents), 0o644)
}

func copyIfDir(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil || !stat.IsDir() {
		return nil
	}
	return copyDir(src, dst)
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

func runGradle(root string, tasks ...string) error {
	cmdName, args, err := gradle.Command(root)
	if err != nil {
		return err
	}
	args = append(args, tasks...)
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = filepath.Join(root, "android")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running %s %s: %w", cmdName, strings.Join(tasks, " "), err)
	}
	return nil
}

func copyArtifacts(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no build outputs found in %s", srcDir)
		}
		return err
	}
	var copied int
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
		fmt.Printf("wrote %s\n", filepath.Join(dstDir, entry.Name()))
		copied++
	}
	if copied == 0 {
		return fmt.Errorf("no artifacts found in %s", srcDir)
	}
	return nil
}

func resolveSources(root string, cfg config.Config) []string {
	sources := cfg.Sources()
	out := make([]string, 0, len(sources))
	for _, src := range sources {
		out = append(out, filepath.Join(root, src))
	}
	return out
}

func loadConfig(root string) (config.Config, error) {
	path := filepath.Join(root, "go2apk.yaml")
	if err := require(path); err != nil {
		return config.Config{}, err
	}
	cfg, err := config.Load(path)
	if err != nil {
		return config.Config{}, err
	}
	cfg.ResolveSecrets()
	return cfg, nil
}
