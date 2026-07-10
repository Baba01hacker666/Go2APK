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

  android:
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
      - name: Install SDK packages
        run: sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" "ndk;27.2.12479018"
      - name: Prepare project
        run: go run ./cmd/go2apk init
      - name: Build debug artifact
        run: go run ./cmd/go2apk build
      - uses: actions/upload-artifact@v4
        with:
          name: go2apk-debug
          path: dist/debug
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
      - name: Install SDK packages
        run: sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" "ndk;27.2.12479018"
      - name: Build release artifacts
        run: go run ./cmd/go2apk release
      - name: Checksums
        run: find dist/release -type f -maxdepth 1 -print0 | xargs -0 sha256sum > dist/release/SHA256SUMS.txt
      - uses: softprops/action-gh-release@v2
        with:
          files: dist/release/*
`
