// Package workflow generates GitHub Actions workflows for a Go2APK app.
package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Baba01hacker666/Go2APK/internal/config"
)

// CLIPath is the command users install to run Go2APK in CI.
const CLIPath = "github.com/Baba01hacker666/Go2APK/cmd/go2apk@latest"

// Init writes CI and release workflows for the project.
//
// The workflows are rendered from the project's configuration so that the SDK
// packages they install match the platform the generated Android project
// compiles against; a workflow that installed a different platform than
// target_sdk would fail in Gradle.
func Init(root string, cfg config.Config) error {
	files := map[string]string{
		".github/workflows/ci.yml":      CI(cfg),
		".github/workflows/release.yml": Release(cfg),
	}
	for name, body := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return err
		}
	}
	fmt.Println("generated GitHub Actions workflows")
	return nil
}

// CI returns a workflow that builds a debug APK on every push and pull request.
func CI(cfg config.Config) string {
	return fmt.Sprintf(`name: CI

on:
  push:
  pull_request:

jobs:
  build:
    name: Build debug APK
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true

      - uses: actions/setup-java@v4
        with:
          distribution: temurin
          java-version: '17'

      - uses: android-actions/setup-android@v3

      - name: Install Android SDK packages
        run: sdkmanager "platform-tools" "platforms;%[1]s" "build-tools;%[2]s" "ndk;%[3]s"

      - name: Install go2apk
        run: go install %[4]s

      - name: Generate Android project
        run: go2apk init

      - name: Build debug APK
        run: go2apk build

      - name: Verify APK exists
        run: test -n "$(find dist/debug -maxdepth 1 -name '*.apk' -print -quit)"

      - uses: actions/upload-artifact@v4
        with:
          name: debug-apk
          path: dist/debug/*.apk
          if-no-files-found: error
`, cfg.Platform(), cfg.BuildTools, cfg.NDK, CLIPath)
}

// Release returns a workflow that builds and publishes signed release artifacts
// when a version tag is pushed.
func Release(cfg config.Config) string {
	return fmt.Sprintf(`name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    name: Build release artifacts
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: actions/setup-java@v4
        with:
          distribution: temurin
          java-version: '17'

      - uses: android-actions/setup-android@v3

      - name: Install Android SDK packages
        run: sdkmanager "platform-tools" "platforms;%[1]s" "build-tools;%[2]s" "ndk;%[3]s"

      - name: Install go2apk
        run: go install %[4]s

      - name: Generate Android project
        run: go2apk init

      # Supply GO2APK_KEYSTORE, GO2APK_KEY_ALIAS, GO2APK_KEYSTORE_PASSWORD and
      # GO2APK_KEY_PASSWORD as repository secrets to publish a signed artifact.
      # Without them the APK is signed with the debug key and is not publishable.
      - name: Build release artifacts
        env:
          GO2APK_KEYSTORE: ${{ secrets.GO2APK_KEYSTORE }}
          GO2APK_KEY_ALIAS: ${{ secrets.GO2APK_KEY_ALIAS }}
          GO2APK_KEYSTORE_PASSWORD: ${{ secrets.GO2APK_KEYSTORE_PASSWORD }}
          GO2APK_KEY_PASSWORD: ${{ secrets.GO2APK_KEY_PASSWORD }}
        run: go2apk release

      - name: Verify release APK exists
        run: test -n "$(find dist/release -maxdepth 1 -name '*.apk' -print -quit)"

      - name: Checksums
        run: find dist/release -maxdepth 1 -type f -print0 | xargs -0 sha256sum > dist/release/SHA256SUMS.txt

      - uses: actions/upload-artifact@v4
        with:
          name: release-artifacts
          path: dist/release/*
          if-no-files-found: error

      - uses: softprops/action-gh-release@v2
        with:
          files: dist/release/*
`, cfg.Platform(), cfg.BuildTools, cfg.NDK, CLIPath)
}
