# Go2APK documentation

Start with the [README](../README.md) for a two-minute overview, then come back
here.

## Using Go2APK

| Document | What it covers |
| --- | --- |
| [Building a UI](UI.md) | Widgets, inline CSS, IDs, event handlers, layout rules |
| [Configuration](CONFIG.md) | Every `go2apk.yaml` key, environment variables, release signing |
| [Commands](CLI.md) | What each `go2apk` command does and when to run it |
| [Android platform APIs](ANDROID.md) | Permissions, intents, broadcast receivers, VPN |
| [Media and animations](MEDIA_AND_ANIMATIONS.md) | Image, Video, Audio, HTML text, `android.Animate` |
| [Generated files](GENERATED_FILES.md) | What Go2APK writes, what to commit, what to ignore |
| [Troubleshooting](TROUBLESHOOTING.md) | `go2apk doctor` and the errors it is meant to prevent |

## Working on Go2APK

| Document | What it covers |
| --- | --- |
| [Architecture](ARCHITECTURE.md) | The five-stage pipeline, package map, invariants to preserve |

## Conventions used in these documents

- Commands are written as `go2apk build`. Run them from the directory holding
  `go2apk.yaml`.
- Paths are relative to that same directory.
- Anything shown as generated — `android/`, `scripts/`, `*_gen.go`,
  `preview.html` — is rewritten on every build and must not be hand-edited.
