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

## Quick start

```bash
go2apk init       # write go2apk.yaml, the Android project, scripts and CI
go2apk doctor     # check the toolchain
go2apk preview    # render preview.html to check the layout in a browser
go2apk build      # build a debug APK into dist/debug
```

Building an APK also needs a JDK 17+, and the Android SDK with platform tools, a
platform matching your `target_sdk`, build tools and the NDK. `go2apk doctor`
tells you exactly which of those is missing.

## Writing an app

```go
package main

import (
	"strconv"

	"github.com/Baba01hacker666/Go2APK/android"
	"github.com/Baba01hacker666/Go2APK/ui"
)

var count int

func main() {
	ui.Run(ui.Column{
		Style: ui.Style{Width: ui.MatchParent, Padding: 24},
		Children: []ui.Widget{
			ui.TextView{ID: "label", Text: "Counter"},
			ui.Button{
				ID:    "add",
				Text:  "Add one",
				Style: ui.Style{Width: ui.WrapContent},
				OnClick: func() {
					count++
					android.UpdateText("label", strconv.Itoa(count))
				},
			},
		},
	})
}
```

The demo calculator in `examples/demo` is a fuller example and is the default
`source` in `go2apk.yaml`.

## Documentation

| | |
| --- | --- |
| [Building a UI](docs/UI.md) | Widgets, inline CSS, IDs, event handlers, layout rules |
| [Configuration](docs/CONFIG.md) | Every `go2apk.yaml` key, environment variables, release signing |
| [Commands](docs/CLI.md) | What each `go2apk` command does |
| [Android platform APIs](docs/ANDROID.md) | Permissions, intents, broadcast receivers, VPN |
| [Media and animations](docs/MEDIA_AND_ANIMATIONS.md) | Image, Video, Audio, HTML text, `android.Animate` |
| [Generated files](docs/GENERATED_FILES.md) | What Go2APK writes, and what to commit |
| [Troubleshooting](docs/TROUBLESHOOTING.md) | `go2apk doctor` and the errors it prevents |
| [Architecture](docs/ARCHITECTURE.md) | How the tool works internally |

## Contributing

```bash
make verify   # gofmt check, go vet, go test -race
make help     # list all targets
```

`gofmt` and `go vet` are enforced in CI, along with the tests, a check that
regenerating an unchanged project changes nothing, and an end-to-end APK build
of the demo app.

Two rules worth knowing before you change the generator, both explained in
[Architecture](docs/ARCHITECTURE.md#invariants-worth-preserving): everything
interpolated into Java or XML must be escaped, and generation must stay
deterministic.

## Licence

No licence has been chosen yet. Add a `LICENSE` file before publishing, since
without one nobody can legally use this code.
