package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/internal/config"
)

func TestInitWritesScaffold(t *testing.T) {
	root := t.TempDir()
	if err := Init(root, false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	pkgDir := filepath.Join("android", "app", "src", "main", "java", "com", "example", "goapp")
	for _, want := range []string{
		"go2apk.yaml",
		"android/settings.gradle",
		"android/build.gradle",
		"android/app/build.gradle",
		"android/app/proguard-rules.pro",
		"android/app/src/main/AndroidManifest.xml",
		filepath.Join(pkgDir, "MainActivity.java"),
		filepath.Join(pkgDir, "Go2ApkActivity.java"),
		filepath.Join(pkgDir, "NativeBridge.java"),
		"scripts/install-sdk.sh",
		"scripts/install-sdk.ps1",
		"scripts/build-go-app.sh",
	} {
		if _, err := os.Stat(filepath.Join(root, want)); err != nil {
			t.Errorf("expected %s to be created: %v", want, err)
		}
	}
}

func TestInitMakesShellScriptsExecutable(t *testing.T) {
	root := t.TempDir()
	if err := Init(root, false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	for _, name := range []string{"scripts/install-sdk.sh", "scripts/build-go-app.sh"} {
		info, err := os.Stat(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("%s was not created: %v", name, err)
		}
		if info.Mode().Perm()&0o111 == 0 {
			t.Errorf("%s is not executable (mode %v)", name, info.Mode().Perm())
		}
	}
}

func TestInitIsIdempotentWithoutForce(t *testing.T) {
	root := t.TempDir()
	if err := Init(root, false); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, "go2apk.yaml")
	if err := os.WriteFile(configPath, []byte("name: Edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Init(root, false); err != nil {
		t.Fatalf("second Init failed: %v", err)
	}
	got, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "name: Edited\n" {
		t.Errorf("Init overwrote an existing file without --force: %q", got)
	}

	if err := Init(root, true); err != nil {
		t.Fatalf("Init --force failed: %v", err)
	}
	got, err = os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(string(got), "name: Edited") {
		t.Errorf("Init --force did not overwrite the file: %q", got)
	}
}

// TestInitConfigRoundTrips makes sure the generated go2apk.yaml can be read
// back, including values that contain characters the old unquoted format used
// to mangle.
func TestInitConfigRoundTrips(t *testing.T) {
	root := t.TempDir()
	if err := Init(root, false); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(filepath.Join(root, "go2apk.yaml"))
	if err != nil {
		t.Fatalf("the generated config does not load: %v", err)
	}
	if cfg.Name != config.Default().Name {
		t.Errorf("Name = %q, want %q", cfg.Name, config.Default().Name)
	}
	if cfg.Package != config.Default().Package {
		t.Errorf("Package = %q, want %q", cfg.Package, config.Default().Package)
	}
	if cfg.Source != config.Default().Source {
		t.Errorf("Source = %q, want %q", cfg.Source, config.Default().Source)
	}
}

// TestInitSettingsMatchesProjectName guards a past inconsistency: settings.gradle
// hard-coded a project name that had nothing to do with the config.
func TestInitSettingsMatchesProjectName(t *testing.T) {
	root := t.TempDir()
	if err := Init(root, false); err != nil {
		t.Fatal(err)
	}
	settings, err := os.ReadFile(filepath.Join(root, "android", "settings.gradle"))
	if err != nil {
		t.Fatal(err)
	}
	want := `rootProject.name = "` + gradleProjectName(config.Default().Name) + `"`
	if !strings.Contains(string(settings), want) {
		t.Errorf("settings.gradle should contain %q, got:\n%s", want, settings)
	}
}

func TestGradleProjectName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"My Go App", "My Go App"},
		{"App: The Sequel", "App The Sequel"},
		{`Quote"Name`, "QuoteName"},
		{"  padded  ", "padded"},
		{"///", "Go2APKApp"},
		{"", "Go2APKApp"},
	}
	for _, tc := range tests {
		if got := gradleProjectName(tc.in); got != tc.want {
			t.Errorf("gradleProjectName(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestInitHonoursTheConfiguredPackage covers a bug where the Java sources were
// written under the *default* package directory even for a project whose
// go2apk.yaml named a different one, leaving two copies of every class in
// mutually invisible packages.
func TestInitHonoursTheConfiguredPackage(t *testing.T) {
	root := t.TempDir()
	custom := "com.acme.widgets"
	if err := os.WriteFile(filepath.Join(root, "go2apk.yaml"), []byte(
		"name: \"Widgets\"\npackage: \""+custom+"\"\nmin_sdk: 24\ntarget_sdk: 36\n"+
			"orientation: portrait\nsource: \".\"\nbuild_tools: \"36.0.0\"\nndk: \"27.2.12479018\"\n",
	), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Init(root, false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	pkgDir := filepath.Join(root, "android/app/src/main/java", filepath.FromSlash(strings.ReplaceAll(custom, ".", "/")))
	for _, name := range []string{"MainActivity.java", "Go2ApkActivity.java", "NativeBridge.java"} {
		path := filepath.Join(pkgDir, name)
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("expected %s under the configured package: %v", name, err)
		}
		// The directory is only half the contract; the file must declare the
		// same package or javac rejects it.
		if !strings.Contains(string(src), "package "+custom+";") {
			t.Errorf("%s does not declare package %s", name, custom)
		}
	}

	// Nothing should have been written under the default package.
	if _, err := os.Stat(filepath.Join(root, "android/app/src/main/java/com/example/goapp")); !os.IsNotExist(err) {
		t.Error("Init wrote into the default package directory as well as the configured one")
	}
}

// TestInitReportsAnInvalidConfig rather than quietly generating a project from
// defaults that do not match what the author wrote.
func TestInitReportsAnInvalidConfig(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go2apk.yaml"), []byte("package: \"Not A Package\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Init(root, false); err == nil {
		t.Fatal("expected Init to reject an invalid package name")
	}
}
