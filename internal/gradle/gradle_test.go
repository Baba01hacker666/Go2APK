package gradle

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

// fakeDistribution builds a zip shaped like a real Gradle distribution: every
// path sits under a single gradle-<version> directory.
func fakeDistribution(t *testing.T, version string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "gradle-"+version+"-bin.zip")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	entries := map[string]string{
		"gradle-" + version + "/bin/gradle":                "#!/bin/sh\\necho gradle\\n",
		"gradle-" + version + "/bin/gradle.bat":            "@echo off\\r\\n",
		"gradle-" + version + "/lib/gradle-launcher-8.jar": "not really a jar",
	}
	for name, body := range entries {
		out, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := out.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestInstallFromArchiveIsIdempotent covers a bug where the distribution was
// renamed onto the install root, so the binary ended up at
// .go2apk/gradle/bin/gradle while every later lookup used
// .go2apk/gradle/gradle-<version>/bin/gradle. The first call therefore
// "succeeded" and the next one re-downloaded, forever.
func TestInstallFromArchiveIsIdempotent(t *testing.T) {
	root := t.TempDir()
	installRoot := filepath.Join(root, ".go2apk", "gradle")
	want := filepath.Join(installRoot, "gradle-"+Version, "bin", "gradle")

	got, err := installFromArchive(fakeDistribution(t, Version), installRoot)
	if err != nil {
		t.Fatalf("installFromArchive failed: %v", err)
	}
	if got != want {
		t.Fatalf("installFromArchive returned %q, want %q", got, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("the binary is not where the returned path says: %v", err)
	}
	if info, err := os.Stat(want); err == nil && info.Mode().Perm()&0o111 == 0 {
		t.Errorf("gradle binary is not executable: %v", info.Mode())
	}

	// The whole point: a second install must land in exactly the same place,
	// so the caller's existence check can short-circuit it.
	got2, err := installFromArchive(fakeDistribution(t, Version), installRoot)
	if err != nil {
		t.Fatalf("second installFromArchive failed: %v", err)
	}
	if got2 != want {
		t.Errorf("second install returned %q, want %q", got2, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Errorf("binary disappeared after reinstall: %v", err)
	}
}

func TestInstallFromArchiveRejectsAnUnexpectedLayout(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "wrong.zip")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	out, err := w.Create("bin/gradle")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := out.Write([]byte("#!/bin/sh\\n")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()

	if _, err := installFromArchive(path, filepath.Join(root, ".go2apk", "gradle")); err == nil {
		t.Fatal("expected an archive without a gradle-<version> directory to be rejected")
	}
}

func TestInstallFromArchiveRejectsEntriesThatEscape(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "evil.zip")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	out, err := w.Create("../escaped/gradle")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := out.Write([]byte("x")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()

	dest := filepath.Join(root, ".go2apk", "gradle")
	if _, err := installFromArchive(path, dest); err == nil {
		t.Fatal("expected a path traversal entry to be rejected")
	}
	if _, err := os.Stat(filepath.Join(root, ".go2apk", "escaped")); !os.IsNotExist(err) {
		t.Error("an entry outside the destination was written anyway")
	}
}

func TestVerifyChecksum(t *testing.T) {
	known := "d725d707bfabd4dfdc958c624003b3c80accc03f7037b5122c4b1d0ef15cecab"

	t.Setenv(checksumEnv, "")
	if err := verifyChecksum(known); err != nil {
		t.Errorf("an unset checksum should warn, not fail: %v", err)
	}

	t.Setenv(checksumEnv, known)
	if err := verifyChecksum(known); err != nil {
		t.Errorf("matching checksum rejected: %v", err)
	}

	t.Setenv(checksumEnv, "sha256:"+known)
	if err := verifyChecksum(known); err != nil {
		t.Errorf("sha256:-prefixed checksum rejected: %v", err)
	}

	t.Setenv(checksumEnv, "0000000000000000000000000000000000000000000000000000000000000000")
	if err := verifyChecksum(known); err == nil {
		t.Error("mismatched checksum accepted")
	}
}
