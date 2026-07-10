# Go2APK

Go2APK is a CLI for turning Go mobile applications into Android APKs, while still generating Android project scaffolding for environments that use Gradle directly.

## Install and run

```bash
go install ./cmd/go2apk
go2apk init
go2apk doctor
go2apk sdk install
go2apk build
go2apk release
```

## Commands

```bash
go run ./cmd/go2apk init           # write config, Android templates, scripts, and workflows
go run ./cmd/go2apk build          # build a debug APK with gomobile, or fall back to Gradle
go run ./cmd/go2apk release        # build a release APK with gomobile, or fall back to Gradle
go run ./cmd/go2apk clean          # remove dist artifacts
go run ./cmd/go2apk doctor         # check Go, Java, Gradle, SDK, sdkmanager, and adb
go run ./cmd/go2apk sdk install    # create Android SDK installer scripts and local SDK directory
go run ./cmd/go2apk workflow init  # regenerate GitHub Actions workflows
```

`init` writes a starter `go2apk.yaml`, Android Gradle project under `android/`, a minimal `MainActivity`, default theme resources, SDK installer scripts, and CI/release workflows. `build` and `release` first invoke `gomobile build -target=android` for the Go package configured by `source` and write APKs to `dist/debug` or `dist/release`. If gomobile is not installed, the commands fall back to the generated Gradle project and still prepare diagnostics when the local environment is incomplete.

## Building real Go apps

This repository includes a demo app at `examples/demo`, and the checked-in `go2apk.yaml` points `source` at that package. Set `source` to a different Go package when building your own app, then install gomobile once:

```bash
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
go2apk build
```

Debug APKs are copied to `dist/debug`; release APKs are copied to `dist/release`. The generated GitHub Actions workflow installs Android SDK packages plus gomobile, builds the demo APK, verifies that an APK exists, and uploads it as a workflow artifact named `go2apk-demo-debug-apk`.

## Android SDK setup

The generated SDK scripts install Android command-line tools, platform tools, API 35, build-tools 35.0.0, and NDK 27.2.12479018.

```bash
./scripts/install-sdk.sh
```

On Windows:

```powershell
./scripts/install-sdk.ps1
```

## GitHub Actions

The repository includes:

- `.github/workflows/ci.yml` for formatting, tests, CLI build, Android SDK setup, and debug artifacts.
- `.github/workflows/release.yml` for tag-based release builds, checksums, and GitHub Release uploads.

## Roadmap

See [plan.md](plan.md) for the full development plan, including deeper gomobile integration, signing, ABI support, examples, and release automation.
