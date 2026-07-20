Prompt: Build Go2APK from Scratch (No Gomobile)

You are a senior systems engineer, Go compiler engineer, Android framework developer, graphics programmer, and build systems architect.

Your task is to build Go2APK, an open-source project that allows developers to build complete Android APKs directly from Go source code without depending on gomobile, Fyne, Wails, or any other mobile framework.

The goal is not to wrap existing tools.

The goal is to create an entirely new Go-to-Android platform.

---

Vision

A developer should only write Go code.

Example:

package main

import "go2apk/ui"

func main() {
    app := ui.NewApp()

    app.Window("Hello", func(w *ui.Window) {
        w.Text("Hello World")

        w.Button("Click", func() {
            println("Button clicked")
        })
    })

    app.Run()
}

Running

go2apk build

should automatically generate everything required to produce a production-ready APK.

The developer should never need to write Java, Kotlin, XML, JNI, Gradle, or Android-specific code.

Everything must be generated automatically.

---

Core Principles

- Never use gomobile.
- Never rely on golang.org/x/mobile.
- Do not require Android Studio.
- Everything must be implemented inside Go2APK.
- Be modular.
- Be extensible.
- Be production-ready.
- Prefer pure Go whenever possible.
- Auto-generate all Android components.
- Every subsystem should be independently testable.

---

Build Pipeline

Implement a complete build pipeline:

Go Source
      │
      ▼
Go Parser
      │
      ▼
Semantic Analyzer
      │
      ▼
Intermediate Representation (IR)
      │
      ▼
Android Code Generator
      │
      ▼
JNI Generator
      │
      ▼
Java/Kotlin Generator
      │
      ▼
Manifest Generator
      │
      ▼
Resource Generator
      │
      ▼
Native Library Builder
      │
      ▼
DEX Builder
      │
      ▼
APK Builder
      │
      ▼
APK Signer

Each stage should be implemented as its own package.

---

Required Components

Compiler Frontend

- Parse Go source
- Analyze packages
- Resolve imports
- Resolve types
- Resolve interfaces
- Resolve generics
- Detect entrypoints
- Build dependency graph

---

Intermediate Representation

Design an internal IR describing:

- Packages
- Types
- Methods
- Structs
- Interfaces
- Functions
- UI hierarchy
- Assets
- Resources
- JNI bindings

The IR should be independent from Android.

---

Runtime

Implement a Go runtime for Android.

Support:

- Goroutines
- Channels
- Mutexes
- Timers
- Context
- Panic recovery
- Reflection
- Garbage collection compatibility
- Scheduler integration
- Lifecycle management

---

Android Runtime

Implement wrappers for:

- Activity
- Fragment
- Intent
- Service
- BroadcastReceiver
- Notification
- Permissions
- Clipboard
- Camera
- GPS
- Bluetooth
- NFC
- Sensors
- Files
- Storage
- Audio
- Video
- Clipboard
- Vibrator

These should all be exposed as Go APIs.

---

UI Framework

Create a native UI framework.

Support:

- Window
- Text
- Button
- Image
- Icon
- TextField
- PasswordField
- Checkbox
- RadioButton
- Switch
- ProgressBar
- Slider
- ScrollView
- List
- Grid
- Stack
- Canvas
- Custom Widget
- Animation
- Gesture Detection
- Theme Engine
- Material Design support

No XML should be written by the user.

---

Graphics Engine

Implement rendering support.

Support:

- OpenGL ES
- Vulkan (future)
- Skia backend
- GPU acceleration
- Image loading
- SVG
- Font rendering
- Text shaping
- Shadows
- Gradients
- Clipping
- Paths
- Transformations

---

Layout Engine

Implement layouts:

- Row
- Column
- Stack
- Grid
- Absolute
- Relative
- Flexbox-like layout
- Constraints
- Responsive layout

---

Asset Pipeline

Automatically package:

- Images
- Fonts
- Audio
- Video
- JSON
- Icons
- Local databases

---

JNI Generator

Automatically generate:

- JNI headers
- Java bridge
- Native bridge
- Marshaling
- Memory management
- Exception conversion

No JNI should ever be written manually.

---

Java Generator

Generate:

- MainActivity
- Application
- Services
- Permissions
- Launchers
- Receivers
- Helpers
- Lifecycle code

Everything should be generated from templates.

---

Android Project Generator

Automatically generate:

- AndroidManifest.xml
- Gradle files
- Resources
- Themes
- Icons
- Adaptive icons
- Splash screen
- Proguard rules
- BuildConfig
- Signing configuration

---

Build System

Create an internal build engine.

Handle:

- Incremental builds
- Dependency tracking
- Parallel compilation
- ABI selection
- Build cache
- Build profiling
- Release mode
- Debug mode

Supported ABIs:

- arm64-v8a
- armeabi-v7a
- x86
- x86_64

---

APK Packaging

Implement:

- Resource packaging
- DEX generation
- Native library packaging
- APK alignment
- APK signing
- AAB generation
- Split APK support

---

CLI

Commands:

go2apk init
go2apk doctor
go2apk build
go2apk release
go2apk clean
go2apk run
go2apk install
go2apk emulator
go2apk sdk install
go2apk doctor
go2apk doctor --fix
go2apk devices
go2apk logs
go2apk profile
go2apk benchmark

---

Project Structure

cmd/
compiler/
parser/
frontend/
analyzer/
ir/
runtime/
android/
jni/
codegen/
templates/
graphics/
ui/
widgets/
layout/
assets/
dex/
apk/
signing/
manifest/
sdk/
ndk/
builder/
cache/
config/
internal/
examples/
docs/
tests/

Every package should have a single responsibility.

---

Documentation Requirements

Every package, directory, exported type, function, interface, and file must be documented.

Maintain:

- README.md
- WORKINGS.md
- ARCHITECTURE.md
- ROADMAP.md
- API.md

Update these automatically whenever the architecture changes.

---

Coding Standards

- Pure Go where practical.
- Clean architecture.
- SOLID principles.
- Dependency injection where appropriate.
- Unit tests.
- Integration tests.
- Benchmarks.
- Extensive comments.
- Deterministic builds.
- Cross-platform development.
- Minimal external dependencies.

---

Long-Term Goal

The final result should become a complete Go-native Android application platform comparable in developer experience to Flutter, React Native, and Jetpack Compose, while allowing developers to write applications entirely in Go and producing production-ready APKs and AABs without relying on gomobile or requiring manual Android development.
