# Project Memory

This file serves as a memory store of important decisions, architectural features, and project state for the Go2APK builder.

## Features Implemented
1. **Multi-Page Architecture**: Replaced `ui.Run()` with `ui.RunApp()` accepting multiple `ui.Page` definitions. The Go parser extracts these pages and the code generator emits a distinct Android `Activity` for each page (e.g., `MainActivity`, `AdvancedCalculatorActivity`).
2. **Activity Navigation**: Added `android.Navigate(activityName string)` to allow seamless traversal between pages by automatically invoking `startActivity` behind a dynamic intent.
3. **HTTP Operations**: Added `android.HTTPGet` to perform standard GET requests natively in Go while injecting required internet permissions (`android.permission.INTERNET`) and cleartext traffic bypass into the Android Manifest.
4. **Native Go Parsing Evaluator**: Dropped external evaluator modules (like `Knetic/govaluate`) in favor of a zero-dependency abstract syntax tree evaluator utilizing Go's native `go/ast` and `go/parser` packages to process math strings.
5. **Code Generation Mode**: Added the `go2apk generate` command to run the front-end parser and write Java native bindings without triggering the expensive Gradle build phase.

## Strict Guidelines for AI Agents (Lessons Learned)
To prevent repeating past mistakes, any AI agent interacting with this repository MUST adhere to the following rules:

1. **DO NOT BUILD GRADLE LOCALLY**: This machine is not meant for heavy Gradle compilations. Do not run `go run ./cmd/go2apk build ...` unless explicitly forced. Instead, always use `go run ./cmd/go2apk generate ...` to view the generated `.java` bindings and verify structural correctness without triggering Gradle.
2. **Handle Stale Generated Files**: If frontend parsing fails on a build/generate command due to syntax errors in `events_gen.go`, always delete `events_gen.go` first and retry. Stale events can cause the AST parser to fail.
3. **No External Math Libraries**: We wrote a custom expression evaluator using `go/ast` and `go/parser` for the calculator demo. Do NOT attempt to install `github.com/Knetic/govaluate` or `github.com/apache/casbin-govaluate`.
4. **Android Permissions**: If adding a new feature that requires permissions (like `HTTPGet`), you must explicitly call `android.Permission(...)` dynamically in the Go code so the builder injects it into `AndroidManifest.xml`.
5. **No Gomobile**: Under no circumstances should `gomobile` or `golang.org/x/mobile` be imported or used. This project is a strict from-scratch native implementation.

## Current State
- The demo calculator application successfully uses all of the above features to render a basic and advanced interface that evaluate complex strings completely within Go.
- `builder.Generate()` bypasses gradle overhead for rapid previewing.
