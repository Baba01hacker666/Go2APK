# Go2APK

Go2APK is an early-stage CLI for turning Go applications into Android project scaffolding and, over time, installable APK/AAB release artifacts.

## Current commands

```bash
go run ./cmd/go2apk init
go run ./cmd/go2apk build
go run ./cmd/go2apk release
go run ./cmd/go2apk clean
go run ./cmd/go2apk doctor
```

`init` writes a starter `go2apk.yaml` and Android Gradle project under `android/`. `build` and `release` currently validate that the project has been initialized and create artifact directories that will be replaced by full Android SDK/gomobile integration in future milestones.

## Roadmap

See [plan.md](plan.md) for the full development plan, including Android SDK setup, Gradle wrapper generation, signing, ABI support, and GitHub Actions automation.
