package config

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestLoadAndroidDependencies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go2apk.yaml")
	contents := "android_dependencies: com.google.android.material:material:1.12.0, androidx.activity:activity:1.9.3\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"com.google.android.material:material:1.12.0", "androidx.activity:activity:1.9.3"}
	if !reflect.DeepEqual(cfg.AndroidDependencies, want) {
		t.Fatalf("dependencies = %#v, want %#v", cfg.AndroidDependencies, want)
	}
}
