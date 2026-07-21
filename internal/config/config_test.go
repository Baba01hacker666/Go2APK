package config

import (
	"os"
	"path/filepath"
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
