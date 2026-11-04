# Go2APK

Go2APK turns a declarative Go UI application into a native Android APK. It
parses your Go source, builds an intermediate representation of the widget tree,
generates the Java Android project from it, and builds the result with Gradle
and the Android NDK.

It does **not** use gomobile, Fyne, or any other mobile framework. The Java
sources, the manifest and the Gradle scripts are all generated from your Go
code.

**Status:** this is a young project. Expect rough edges, and please report
anything that generates an APK which does not run.

## Install

```bash
go install github.com/Baba01hacker666/Go2APK/cmd/go2apk@latest
```

Building an APK also needs a JDK 17+, and the Android SDK with platform tools,
a platform matching your `target_sdk`, build tools and the NDK. `go2apk doctor`
tells you what is missing.

## Quick start

```bash
go2apk init          # write go2apk.yaml, the Android project, scripts and CI
go2apk doctor        # check the toolchain
go2apk sdk install   # write scripts/install-sdk.sh and install-sdk.ps1
go2apk preview       # render preview.html to check the layout in a browser
go2apk build         # build a debug APK into dist/debug
go2apk release       # build a release APK and AAB into dist/release
```

Point `source` in `go2apk.yaml` at the Go package that calls `ui.Run` or
`ui.RunApp`. The demo calculator lives in `examples/demo` and is the default.

## Writing an app

```go
package main

import (
	"strconv"

	"github.com/Baba01hacker666/Go2APK/android"
	"github.com/Baba01hacker666/Go2APK/ui"
)

var count int

func reset() { count = 0 }

func main() {
	ui.Run(ui.Column{
		Style: ui.Style{Width: ui.MatchParent, Padding: 24},
		Children: []ui.Widget{
			ui.TextView{ID: "label", Text: "Counter"},
			ui.Button{
				ID:      "add",
				Text:    "Add one",
				Style:   ui.Style{Width: ui.WrapContent},
				OnClick: func() {
					count++
					android.UpdateText("label", strconv.Itoa(count))
				},
			},
			ui.Button{ID: "reset", Text: "Reset",
				Style: ui.Style{Width: ui.WrapContent}, OnClick: reset},
		},
	})
}
```

Handlers may be a named function or an inline `func() {}`. A widget with a
handler but no `ID` is given a generated one, since the event is dispatched by
name over JNI.

Further reading: [UI styling and CSS](docs/UI.md),
[Android platform APIs](docs/ANDROID.md),
[media and animations](docs/MEDIA_AND_ANIMATIONS.md), and
[how the tool works](WORKINGS.md).

## Commands

| Command | What it does |
| --- | --- |
| `go2apk init [--force]` | Write `go2apk.yaml`, the Android project, helper scripts and CI workflows. Existing files are kept unless `--force` is given. |
| `go2apk check` | Parse the app, validate the widget tree and report layout warnings. |
| `go2apk generate` | Regenerate the Android project without running Gradle. |
| `go2apk preview` | Render `preview.html`, an approximation of the Android layout. |
| `go2apk build [--obfuscate]` | Build a debug APK into `dist/debug`. |
| `go2apk release [--obfuscate]` | Build a release APK and AAB into `dist/release`. |
| `go2apk clean` | Remove `dist/`. |
| `go2apk doctor` | Report which build prerequisites are present. |
| `go2apk sdk install` | Write the SDK installer scripts for the configured toolchain. |
| `go2apk workflow init` | Regenerate the GitHub Actions workflows. |

## Configuration

`go2apk init` writes a commented `go2apk.yaml`. Unknown keys and invalid values
are rejected rather than ignored, so a typo fails immediately.

| Key | Meaning |
| --- | --- |
| `name`, `package`, `version` | App identity. `package` must be a valid Android package name. |
| `min_sdk`, `target_sdk` | API levels. `target_sdk` must be at least `min_sdk`. |
| `build_tools`, `ndk` | Android toolchain versions used by the SDK installers and CI. |
| `orientation`, `theme` | Activity attributes. |
| `source` | Directory holding the package that calls `ui.Run` / `ui.RunApp`. |
| `source_dirs` | Extra first-party directories to parse and to track as Gradle inputs. |
| `gradle_dependencies` | Android Maven artifacts, comma separated. |
| `go_build_tags` | Go build tags for your own build scripts. |
| `keystore`, `key_alias`, `keystore_password`, `key_password` | Release signing. |
| `obfuscate` | Enable R8 minification and resource shrinking for release builds. |

Signing passwords are better supplied through `GO2APK_KEYSTORE`,
`GO2APK_KEY_ALIAS`, `GO2APK_KEYSTORE_PASSWORD` and `GO2APK_KEY_PASSWORD` so they
are never committed; the environment takes precedence over the file. With no
keystore configured, release builds fall back to the debug key, which is
installable but not publishable — `go2apk release` says so.

Other environment variables: `GO2APK_GRADLE` selects a Gradle binary, and
`GO2APK_GRADLE_SHA256` pins the checksum of the Gradle distribution Go2APK
downloads when a project has no wrapper.

## What is generated

`android/`, `scripts/`, `preview.html` and the `events_gen.go` /
`jni_exports_gen.go` files written next to your app source are all derived from
`go2apk.yaml` and your Go code. They are regenerated on every build and are not
meant to be edited or committed; the `.gitignore` that `go2apk init` expects you
to use excludes them. Customise the build through `go2apk.yaml`, not by editing
generated files.

## Development

```bash
make verify   # gofmt check, go vet, go test -race
make test
make cover
make build
make help     # list all targets
```

`gofmt` and `go vet` are enforced in CI, along with the tests and an end-to-end
APK build of the demo app.

## Licence

No licence has been chosen yet. Add a `LICENSE` file before publishing, since
without one nobody can legally use this code.
