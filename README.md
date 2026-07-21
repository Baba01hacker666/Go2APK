# Go2APK

Go2APK is a toolchain that converts declarative Go UI applications into native Android APKs **from scratch**. It parses your Go code, constructs an Intermediate Representation (IR), generates native Android views via Java, and compiles the app using the Android NDK and Gradle. It explicitly **avoids using gomobile**.

## Documentation
* [UI Styling and CSS](docs/UI.md)
* [Android Platform APIs (Intents, Permissions, Broadcasts)](docs/ANDROID.md)
* [Media and Animations (Images, Video, Audio)](docs/MEDIA_AND_ANIMATIONS.md)
* [Internal Workings](WORKINGS.md)

## Install and run

```bash
go install ./cmd/go2apk
go2apk new my-app
cd my-app
go2apk init
go2apk doctor
go2apk sdk install
go2apk preview
go2apk build
go2apk release
```

## Commands

```bash
go run ./cmd/go2apk new my-app     # create a split UI/logic starter app
go run ./cmd/go2apk init           # write config, Android templates, scripts, and workflows
go run ./cmd/go2apk preview        # instantly generate a preview.html mocking your Android UI
go run ./cmd/go2apk build          # build a debug APK via Gradle and the NDK
go run ./cmd/go2apk release        # build a release APK via Gradle and the NDK
go run ./cmd/go2apk clean          # remove dist artifacts
go run ./cmd/go2apk doctor         # check Go, Java, Gradle, SDK, sdkmanager, and adb
go run ./cmd/go2apk sdk install    # create Android SDK installer scripts and local SDK directory
go run ./cmd/go2apk workflow init  # regenerate GitHub Actions workflows
```

## How It Works

1. `init` writes a starter `go2apk.yaml`, an Android Gradle project under `android/`, and helper scripts.
2. `preview` builds an HTML representation of your Go widget tree so you can test layout and CSS styling instantly in your browser without Android SDK dependencies.
3. `build` parses the Go AST from the package defined in `go2apk.yaml`, dynamically generates Java Android layout instructions, and runs a Gradle build that compiles your Go code into a native JNI library linked to a custom `NativeBridge.java`.

## Building real Go apps

This repository includes a demo app at `examples/demo`, and the checked-in `go2apk.yaml` points `source` at that package. Set `source` to a different Go package when building your own app.

```bash
go2apk build
```

Debug APKs are copied to `dist/debug`; release APKs are copied to `dist/release`.

## Obfuscation

Set `obfuscate: true` in `go2apk.yaml`, or pass `--obfuscate` to `go2apk build`/`go2apk release`, to enable Android release obfuscation in the generated Gradle project. This turns on R8 minification and resource shrinking for the release build type and uses `android/app/proguard-rules.pro` for app-specific keep rules.

## Android SDK setup

The generated SDK scripts install Android command-line tools, platform tools, API 36, build-tools 36.0.0, and NDK 27.2.12479018.

```bash
./scripts/install-sdk.sh
```

## GitHub Actions

The repository includes:

- `.github/workflows/ci.yml` for formatting, tests, CLI build, Android SDK setup, and debug artifacts.
- `.github/workflows/release.yml` for tag-based release builds, checksums, and GitHub Release uploads.


## App structure and dependencies

Use `go2apk new my-app` to create a starter app where UI composition and app logic live in separate Go files (`home_ui.go` and `logic.go`). UI files can return reusable `ui.Component` functions, while logic files hold event handlers and business code.

Go dependencies work like normal Go modules: run `go get example.com/package` inside your app module and import it from your logic files. Android Maven libraries can be added to `go2apk.yaml` with `android_dependencies`:

```yaml
android_dependencies: com.google.android.material:material:1.12.0, androidx.activity:activity:1.9.3
```

The generated Gradle app automatically turns those coordinates into `implementation` dependencies.
