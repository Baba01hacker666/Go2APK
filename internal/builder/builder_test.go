package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/internal/android"
	"github.com/Baba01hacker666/Go2APK/internal/config"
)

func TestRemoveGeneratedGoOnlyRemovesTheToolsOwnFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Default()
	cfg.Source = "."

	// A stale generated file is deleted before parsing, so it cannot be read
	// back as if it were the author's code.
	stale := filepath.Join(dir, "events_gen.go")
	if err := os.WriteFile(stale, []byte(android.GeneratedFileHeader+"\npackage main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ours := filepath.Join(dir, "jni_exports_gen.go")
	if err := os.WriteFile(ours, []byte(android.GeneratedFileHeader+"\npackage main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A file of the same name the author wrote must survive.
	theirs := filepath.Join(dir, "handwritten.go")
	if err := os.WriteFile(theirs, []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Point the config at the temp directory.
	cfg.Source = "."
	root := filepath.Dir(dir)
	cfg.Source = filepath.Base(dir)
	if err := removeGeneratedGo(root, cfg); err != nil {
		t.Fatalf("removeGeneratedGo: %v", err)
	}
	for _, path := range []string{stale, ours} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("expected %s to be removed", path)
		}
	}
	if _, err := os.Stat(theirs); err != nil {
		t.Errorf("the author's file was removed: %v", err)
	}
}

func TestWriteGeneratedGoWritesWhenNeeded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events_gen.go")

	if err := writeGeneratedGo(path, android.GeneratedFileHeader+"\npackage main\n"); err != nil {
		t.Fatalf("writeGeneratedGo: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), android.GeneratedFileHeader) {
		t.Errorf("generated file is missing its header:\n%s", b)
	}
}

func TestReportWarningsSilentWhenClean(t *testing.T) {
	var out strings.Builder
	ReportWarnings(&out, nil)
	if out.String() != "" {
		t.Errorf("expected no output for a clean program, got %q", out.String())
	}
}
