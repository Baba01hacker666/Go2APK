package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/internal/android"
)

func TestWriteGeneratedGoRemovesStaleFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events_gen.go")

	// A file the tool generated is removed once it is no longer needed; this is
	// what stops a stale events_gen.go from breaking the next parse.
	stale := android.GeneratedFileHeader + "\npackage main\n\nfunc init() { android.RegisterEvent(\"old\", Old) }\n"
	if err := os.WriteFile(path, []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeGeneratedGo(path, false, ""); err != nil {
		t.Fatalf("writeGeneratedGo: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected the stale generated file to be removed")
	}
}

func TestWriteGeneratedGoKeepsUserFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events_gen.go")

	// A file the user wrote must never be deleted, even though it has the same
	// name the tool uses.
	if err := os.WriteFile(path, []byte("package main\n\nfunc helper() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeGeneratedGo(path, false, ""); err != nil {
		t.Fatalf("writeGeneratedGo: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("the user's file was removed: %v", err)
	}
}

func TestWriteGeneratedGoWritesWhenNeeded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events_gen.go")

	if err := writeGeneratedGo(path, true, android.GeneratedFileHeader+"\npackage main\n"); err != nil {
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
