Go2APK - Master Development Plan

Vision

Go2APK is an open-source tool that converts Go applications into installable Android APKs with minimal developer effort.

The goal is to allow developers to write applications primarily in Go while automatically generating a complete Android project, building an APK, signing it, and producing release artifacts.

The project should support both CLI and GitHub Actions automation.

---

Primary Objectives

- Convert Go source code into Android applications.
- Support GUI and non-GUI applications.
- Generate complete Android projects automatically.
- Produce Debug and Release APKs.
- Produce Android App Bundles (AAB).
- Support ARM64, ARMv7, x86 and x86_64.
- One-command builds.
- Cross-platform development.
- CI/CD ready.
- Beginner friendly.

---

Repository Layout

Go2APK/

cmd/
    go2apk/

internal/
    builder/
    android/
    gradle/
    sdk/
    signing/
    templates/
    gomobile/
    assets/
    workflow/
    util/

templates/
    android/
        app/
        gradle/
        manifests/

examples/

scripts/

.github/
    workflows/

docs/

testdata/

PLAN.md
README.md
LICENSE

---

Core Features

Project Initialization

go2apk init

Creates

- Android project
- Gradle wrapper
- Manifest
- Assets
- Icons
- Default package

---

Build

go2apk build

Automatically

- validates Go
- downloads dependencies
- prepares Android SDK
- generates bindings
- compiles native libraries
- creates APK
- signs APK
- outputs APK

---

Release Build

go2apk release

Outputs

- Release APK
- AAB
- Mapping files
- Checksums

---

Clean

go2apk clean

Removes

- build cache
- temporary files
- generated Android project

---

Doctor

go2apk doctor

Checks

- Go version
- Java
- Android SDK
- NDK
- Gradle
- Environment variables
- Missing packages

---

SDK Installer

go2apk sdk install

Downloads

- Android SDK
- Build Tools
- Platform Tools
- NDK
- Commandline Tools

Automatically configures environment.

---

Android Support

Minimum SDK configurable

Target SDK latest stable

Support

- ARM64
- ARMv7
- x86
- x86_64

---

Configuration File

Example

go2apk.yaml

name:
package:
version:
icon:
permissions:
min_sdk:
target_sdk:
orientation:
theme:

---

Permissions

Allow automatic generation of

- Internet
- Storage
- Notifications
- Camera
- Bluetooth
- NFC
- GPS
- Foreground Service
- Microphone
- Contacts

---

Assets

Support

- Images
- Fonts
- Videos
- Audio
- HTML
- CSS
- JavaScript
- SQLite
- JSON

---

UI

Future support

- Fyne
- Gio
- WebView
- Native Android XML
- Jetpack Compose wrapper

---

Example Projects

Create example applications.

Examples include

- Hello World
- Calculator
- File Manager
- HTTP Client
- REST API Client
- QR Scanner
- Camera Demo
- SQLite Demo
- GPS Demo
- Bluetooth Demo
- Notification Demo
- Background Service
- WebView Browser
- Markdown Viewer
- Settings App
- Todo App

Each example should compile successfully in CI.

---

Scripts

Create helper scripts.

scripts/

install-sdk.sh
install-sdk.ps1

build.sh
build.ps1

release.sh

clean.sh

doctor.sh

lint.sh

format.sh

test.sh

package.sh

generate-icons.sh

generate-keystore.sh

download-dependencies.sh

---

Documentation

Create

README.md

Installation Guide

Quick Start

FAQ

Troubleshooting

Architecture

Contributing

Developer Guide

Plugin Guide

Release Guide

Examples Guide

---

GitHub Workflows

CI

Every push

- Go formatting
- Lint
- Unit tests
- Build
- Generate APK
- Upload artifacts

---

Release

On tag

- Build APK
- Build AAB
- Generate SHA256
- Create GitHub Release
- Upload artifacts

---

Nightly

Daily

- Build all examples
- Test on latest Go
- Test on latest Android SDK
- Dependency updates

---

Security

Run

- govulncheck
- CodeQL
- Dependabot
- Secret scanning

---

Testing

Unit Tests

Integration Tests

Android Emulator Tests

Build Tests

CLI Tests

Workflow Tests

Example Tests

Regression Tests

---

Code Quality

Use

- gofmt
- gofumpt
- golangci-lint
- govulncheck
- staticcheck

Maintain high coverage.

---

Milestone 1

- CLI
- Project generation
- Android template
- APK generation
- Debug builds

---

Milestone 2

- Signing
- Release APK
- AAB
- Configuration file
- Assets

---

Milestone 3

- GUI support
- WebView support
- Examples
- GitHub Actions
- Documentation

---

Milestone 4

- Plugin system
- Automatic SDK installer
- Automatic icon generation
- Automatic keystore generation
- Build optimization

---

Stretch Goals

- Hot Reload
- Live Logcat
- APK Analyzer
- Bundle Analyzer
- Play Store upload
- Firebase integration
- Crash reporting
- Plugin Marketplace
- Template Marketplace
- Visual project generator
- Remote cloud builds
- Docker build image
- Reproducible builds
- Incremental compilation
- Multi-module support
- WASM support
- Desktop-to-Android migration tools

---

Definition of Done

A new user should be able to clone the repository, install Go2APK, execute:

go2apk init
go2apk build

and receive a working Android APK without manually configuring Gradle, the Android SDK, or project templates whenever automatic setup is possible. The repository must include documentation, examples, helper scripts, automated GitHub workflows, unit tests, integration tests, and release automation.
