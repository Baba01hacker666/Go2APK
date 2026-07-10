package workflow

import (
	"fmt"
	"os"
	"path/filepath"
)

// Init writes GitHub Actions workflows for Go validation, Android builds, and releases.
func Init(root string) error {
	files := map[string]string{
		".github/workflows/ci.yml":      CIYAML,
		".github/workflows/release.yml": ReleaseYAML,
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

const CIYAML = `name: CI

on:
  push:
  pull_request:

jobs:
  go:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - name: Format check
        run: test -z "$(gofmt -l .)"
      - name: Test
        run: go test ./...
      - name: Build CLI
        run: go build ./cmd/go2apk

  demo-apk:
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
        run: sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" "ndk;27.2.12479018"
      - name: Install gomobile
        run: |
          go install golang.org/x/mobile/cmd/gomobile@latest
          gomobile init
      - name: Resolve demo mobile module
        run: go mod tidy
        working-directory: examples/demo
      - name: Build demo debug APK
        run: go run ./cmd/go2apk build
      - name: Verify demo APK exists
        run: test -n "$(find dist/debug -maxdepth 1 -name '*.apk' -print -quit)"
      - uses: actions/upload-artifact@v4
        with:
          name: go2apk-demo-debug-apk
          path: dist/debug/*.apk
          if-no-files-found: error
`

const ReleaseYAML = `name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
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
        run: sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" "ndk;27.2.12479018"
      - name: Install gomobile
        run: |
          go install golang.org/x/mobile/cmd/gomobile@latest
          gomobile init
      - name: Resolve demo mobile module
        run: go mod tidy
        working-directory: examples/demo
      - name: Build release artifacts
        run: go run ./cmd/go2apk release
      - name: Verify release APK exists
        run: test -n "$(find dist/release -maxdepth 1 -name '*.apk' -print -quit)"
      - name: Checksums
        run: find dist/release -maxdepth 1 -type f -print0 | xargs -0 sha256sum > dist/release/SHA256SUMS.txt
      - uses: actions/upload-artifact@v4
        with:
          name: go2apk-demo-release-apk
          path: dist/release/*
          if-no-files-found: error
      - uses: softprops/action-gh-release@v2
        with:
          files: dist/release/*
`
