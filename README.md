# Go2APK

Go2APK is an early-stage CLI for turning Go applications into Android project scaffolding and, over time, installable APK/AAB release artifacts.

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
go run ./cmd/go2apk build          # validate and run assembleDebug when Gradle/SDK are available
go run ./cmd/go2apk release        # validate and run assembleRelease + bundleRelease when available
go run ./cmd/go2apk clean          # remove dist artifacts
go run ./cmd/go2apk doctor         # check Go, Java, Gradle, SDK, sdkmanager, and adb
go run ./cmd/go2apk sdk install    # create Android SDK installer scripts and local SDK directory
go run ./cmd/go2apk workflow init  # regenerate GitHub Actions workflows
```

`init` writes a starter `go2apk.yaml`, Android Gradle project under `android/`, a minimal `MainActivity`, default theme resources, SDK installer scripts, and CI/release workflows. `build` and `release` run Gradle tasks when Android SDK tooling is present and always prepare `dist/debug` or `dist/release` so CI can upload useful diagnostics when a local environment is incomplete.

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
