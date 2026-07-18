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
