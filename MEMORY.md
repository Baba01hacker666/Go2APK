# Project Memory

This file serves as a memory store of important decisions, architectural features, and project state for the Go2APK builder.

## Features Implemented
1. **Multi-Page Architecture**: Replaced `ui.Run()` with `ui.RunApp()` accepting multiple `ui.Page` definitions. The Go parser extracts these pages and the code generator emits a distinct Android `Activity` for each page (e.g., `MainActivity`, `AdvancedCalculatorActivity`).
2. **Activity Navigation**: Added `android.Navigate(activityName string)` to allow seamless traversal between pages by automatically invoking `startActivity` behind a dynamic intent.
3. **HTTP Operations**: Added `android.HTTPGet` to perform standard GET requests natively in Go while injecting required internet permissions (`android.permission.INTERNET`) and cleartext traffic bypass into the Android Manifest.
4. **Native Go Parsing Evaluator**: Dropped external evaluator modules (like `Knetic/govaluate`) in favor of a zero-dependency abstract syntax tree evaluator utilizing Go's native `go/ast` and `go/parser` packages to process math strings.
5. **Code Generation Mode**: Added the `go2apk generate` command to run the front-end parser and write Java native bindings without triggering the expensive Gradle build phase.

## Current State
- The demo calculator application successfully uses all of the above features to render a basic and advanced interface that evaluate complex strings completely within Go.
- `builder.Generate()` bypasses gradle overhead for rapid previewing.
