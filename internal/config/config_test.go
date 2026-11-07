package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadObfuscate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go2apk.yaml")
	if err := os.WriteFile(path, []byte("obfuscate: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Obfuscate {
		t.Fatal("expected obfuscate to load as true")
	}
}

func TestLoadAdvancedLists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go2apk.yaml")
	contents := "source: ./app\nsource_dirs: ./app/ui, ./app/logic, ./app/ui\ngradle_dependencies: androidx.appcompat:appcompat:1.7.0, com.google.android.material:material:1.12.0\ngo_build_tags: pro, android\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := cfg.Sources(), []string{"./app", "./app/ui", "./app/logic"}; len(got) != len(want) || got[2] != want[2] {
		t.Fatalf("Sources() = %#v, want %#v", got, want)
	}
	if len(cfg.GradleDependencies) != 2 {
		t.Fatalf("expected two gradle dependencies, got %#v", cfg.GradleDependencies)
	}
	if len(cfg.GoBuildTags) != 2 {
		t.Fatalf("expected two Go build tags, got %#v", cfg.GoBuildTags)
	}
}

func TestLoadRejectsUnknownKeys(t *testing.T) {
	// A misspelled key used to be ignored silently, so a project could appear
	// to configure something that never took effect.
	_, err := Load(writeConfig(t, "nmae: Typo\n"))
	if err == nil {
		t.Fatal("expected an error for an unknown key")
	}
	if !strings.Contains(err.Error(), "unknown key") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		message string
	}{
		{"bad package case", "package: Com.Example.App", "valid Android package"},
		{"package with dash", "package: com.example.my-app", "valid Android package"},
		{"reserved package segment", "package: com.example.class", "reserved word"},
		{"min sdk too low", "min_sdk: 16", "below the minimum"},
		{"target below min", "min_sdk: 30\ntarget_sdk: 24", "greater than or equal"},
		{"unknown orientation", "orientation: sidewards", "not one of"},
		{"empty name", "name: \"\"", "must not be empty"},
		{"empty source", "source: \"\"", "must name the directory"},
		{"non numeric sdk", "min_sdk: twenty", "invalid min_sdk"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load(writeConfig(t, tc.config))
			if err == nil {
				t.Fatalf("expected an error for %q", tc.config)
			}
			if !strings.Contains(err.Error(), tc.message) {
				t.Errorf("error %q does not mention %q", err, tc.message)
			}
		})
	}
}

func TestDefaultConfigIsValid(t *testing.T) {
	if err := Default().Validate(); err != nil {
		t.Errorf("the shipped default configuration is not valid: %v", err)
	}
}

func TestLoadReportsLineNumbers(t *testing.T) {
	_, err := Load(writeConfig(t, "name: fine\npackage: com.example\nnotakey: 1\n"))
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), ":3:") {
		t.Errorf("error %q does not point at line 3", err)
	}
}

func TestLoadIgnoresCommentsAndBlankLines(t *testing.T) {
	cfg, err := Load(writeConfig(t, "# a comment\n\nname: Commented App\n"))
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Name != "Commented App" {
		t.Errorf("Name = %q, want Commented App", cfg.Name)
	}
	if cfg.Package != Default().Package {
		t.Errorf("Package = %q, want the default %q", cfg.Package, Default().Package)
	}
}

func writeConfig(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "go2apk.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
