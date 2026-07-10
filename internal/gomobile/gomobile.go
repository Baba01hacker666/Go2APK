package gomobile

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go2apk/go2apk/internal/config"
)

// Options describes a gomobile Android build.
type Options struct {
	Root    string
	Config  config.Config
	Release bool
	Stdout  io.Writer
	Stderr  io.Writer
}

// BuildAPK uses gomobile to compile the configured Go package into an Android APK.
func BuildAPK(opts Options) (string, error) {
	if opts.Root == "" {
		return "", errors.New("project root is required")
	}
	if opts.Config.Source == "" {
		opts.Config.Source = "."
	}
	gomobile, err := exec.LookPath("gomobile")
	if err != nil {
		return "", fmt.Errorf("gomobile executable was not found; install it with: go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init")
	}
	outDir := filepath.Join(opts.Root, "dist", buildKind(opts.Release))
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}
	apkName := safeArtifactName(opts.Config.Name, buildKind(opts.Release)+".apk")
	apkPath := filepath.Join(outDir, apkName)

	buildDir := opts.Root
	buildSource := opts.Config.Source
	if sourceDir, ok := sourceModuleDir(opts.Root, opts.Config.Source); ok {
		buildDir = sourceDir
		buildSource = "."
	}

	args := []string{"build", "-target=android", "-androidapi", fmt.Sprint(opts.Config.MinSDK), "-o", apkPath}
	if opts.Release {
		args = append(args, "-ldflags=-s -w")
	}
	args = append(args, buildSource)

	cmd := exec.Command(gomobile, args...)
	cmd.Dir = buildDir
	cmd.Stdout = opts.Stdout
	cmd.Stderr = opts.Stderr
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("gomobile build failed: %w", err)
	}
	return apkPath, nil
}

func buildKind(release bool) string {
	if release {
		return "release"
	}
	return "debug"
}

func safeArtifactName(name, suffix string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		name = "app"
	}
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-', r == '_', r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	cleaned := strings.Trim(b.String(), "-_.")
	if cleaned == "" {
		cleaned = "app"
	}
	return cleaned + "-" + suffix
}

func sourceModuleDir(root, source string) (string, bool) {
	if source == "" || source == "." {
		return "", false
	}
	if strings.HasPrefix(source, "./") || strings.HasPrefix(source, "../") {
		dir := filepath.Clean(filepath.Join(root, source))
		if info, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !info.IsDir() {
			return dir, true
		}
	}
	return "", false
}
