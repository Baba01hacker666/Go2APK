package gradle

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// Version is the Gradle release Go2APK provisions when a project has no
	// wrapper. It must satisfy the Android Gradle Plugin version the generated
	// build.gradle pins.
	Version = "8.9"

	// checksumEnv pins the SHA-256 of the distribution. Set it to make the
	// download verifiable; when it is empty the download proceeds unverified and
	// Go2APK says so.
	checksumEnv = "GO2APK_GRADLE_SHA256"

	downloadTimeout = 10 * time.Minute
)

// DistributionURL is the location Gradle is fetched from.
func DistributionURL() string {
	return "https://services.gradle.org/distributions/gradle-" + Version + "-bin.zip"
}

// EnsureGradle returns the path to a usable Gradle binary inside the project,
// downloading and unpacking it on first use.
//
// The distribution is unpacked into a staging directory and moved into place
// only once it is complete, so an interrupted download can never leave a
// half-installed Gradle that later fails in a confusing way.
func EnsureGradle(root string) (string, error) {
	installRoot := filepath.Join(root, ".go2apk", "gradle")
	gradleBin := filepath.Join(installRoot, "gradle-"+Version, "bin", "gradle")
	if _, err := os.Stat(gradleBin); err == nil {
		return gradleBin, nil
	}

	if err := os.MkdirAll(filepath.Dir(installRoot), 0o755); err != nil {
		return "", err
	}

	fmt.Printf("Downloading Gradle %s...\n", Version)
	archive, err := download(DistributionURL())
	if err != nil {
		return "", err
	}
	defer os.Remove(archive)

	staging, err := os.MkdirTemp(filepath.Dir(installRoot), ".gradle-staging-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(staging)

	if err := unzip(archive, staging); err != nil {
		return "", fmt.Errorf("unpacking Gradle: %w", err)
	}
	if _, err := os.Stat(filepath.Join(staging, "gradle-"+Version, "bin", "gradle")); err != nil {
		return "", fmt.Errorf("Gradle archive did not contain gradle-%s/bin/gradle: %w", Version, err)
	}

	// Replace any previous attempt, then move the verified tree into place.
	if err := os.RemoveAll(installRoot); err != nil {
		return "", err
	}
	if err := os.Rename(filepath.Join(staging, "gradle-"+Version), installRoot); err != nil {
		return "", err
	}
	if err := os.Chmod(gradleBin, 0o755); err != nil {
		return "", err
	}
	return gradleBin, nil
}

// download fetches url to a temporary file and verifies it when a checksum was
// configured through GO2APK_GRADLE_SHA256.
func download(url string) (string, error) {
	client := &http.Client{Timeout: downloadTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("downloading %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("downloading %s: unexpected status %s", url, resp.Status)
	}

	file, err := os.CreateTemp("", "go2apk-gradle-*.zip")
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	written, err := io.Copy(io.MultiWriter(file, hash), resp.Body)
	if err != nil {
		os.Remove(file.Name())
		return "", fmt.Errorf("downloading %s: %w", url, err)
	}
	if resp.ContentLength > 0 && written != resp.ContentLength {
		os.Remove(file.Name())
		return "", fmt.Errorf("downloading %s: got %d bytes, want %d", url, written, resp.ContentLength)
	}

	if err := verifyChecksum(hex.EncodeToString(hash.Sum(nil))); err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

// verifyChecksum compares the downloaded archive against GO2APK_GRADLE_SHA256.
// With no checksum configured it warns rather than silently trusting the bytes.
func verifyChecksum(got string) error {
	want := strings.TrimSpace(os.Getenv(checksumEnv))
	if want == "" {
		fmt.Printf("warning: %s is not set; the Gradle download cannot be verified\n", checksumEnv)
		fmt.Printf("         sha256 of the downloaded archive is %s\n", got)
		return nil
	}
	want = strings.ToLower(strings.TrimPrefix(want, "sha256:"))
	if want != got {
		return fmt.Errorf("checksum mismatch: %s is %s but the archive downloaded is %s", checksumEnv, want, got)
	}
	return nil
}

// unzip extracts archive into dest, refusing entries that would escape it.
func unzip(archive, dest string) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	root := filepath.Clean(dest) + string(os.PathSeparator)
	for _, f := range r.File {
		target := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(target, root) {
			return fmt.Errorf("archive entry %q escapes the destination directory", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := extractFile(f, target); err != nil {
			return err
		}
	}
	return nil
}

func extractFile(f *zip.File, target string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode().Perm()|0o100)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, rc); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// ErrNoGradle is returned when no usable Gradle can be found or provisioned.
var ErrNoGradle = errors.New("gradle not available")
