# Commands

Every command is run from the directory that holds `go2apk.yaml`. `go2apk` with
no arguments, or with an unrecognised flag, prints the usage text and exits
non-zero — it never silently ignores a flag it does not understand.

## The usual flow

```bash
go2apk init          # once, per project
go2apk doctor        # once, per machine
go2apk preview       # while iterating on layout
go2apk check         # while iterating on the widget tree
go2apk build         # produce an APK
go2apk release       # produce a signed APK and AAB
```

`build` and `release` both generate before they invoke Gradle, so running
`generate` by hand is only useful when you want to inspect the output without
waiting for a build.

---

## `go2apk init [--force]`

Writes `go2apk.yaml`, the Android project under `android/`, the helper scripts
under `scripts/`, and GitHub Actions workflows under `.github/workflows/`.

Without `--force`, a file that already exists is left alone. That makes `init`
safe to re-run, and it means the config you have edited is never overwritten.

The generated workflows are derived from `go2apk.yaml`, so they build against the
same SDK platform, build-tools and NDK versions that `scripts/install-sdk.sh`
installs. Regenerate them with `workflow init` after changing those versions.

## `go2apk check`

Parses the app, lowers it to the IR and reports layout warnings. It writes
nothing. This is the fastest way to see whether the tool understands your source.

It also compiles the package, so a genuine Go error in the app is reported here
rather than surfacing later as a Gradle failure.

## `go2apk generate`

Regenerates the Android project and the two Go files written next to your app
source, without running Gradle. Use it to inspect the generated Java, or to
refresh `android/` after editing `go2apk.yaml`.

Generation is deterministic: running it twice on an unchanged project produces
byte-identical output. CI enforces this, because a project that regenerates
differently every time cannot be diffed or reviewed.

Before parsing, Go2APK deletes the Go files it previously generated. They live
inside the package being parsed, so a stale one would otherwise be read as though
it were your code.

## `go2apk preview`

Renders `preview.html`, an approximation of the Android layout rendered as HTML.
It is a layout aid, not a simulator: it shows hierarchy, spacing and colour, and
it does not run any of your Go code.

## `go2apk build [--obfuscate]`

Generates and then builds a debug APK into `dist/debug/`. `--obfuscate` enables
R8 minification and resource shrinking, and overrides `obfuscate` in
`go2apk.yaml` for this build.

The command fails if no APK was produced. It does not report success on a build
that merely ran.

## `go2apk release [--obfuscate]`

Generates and then builds a release APK and AAB into `dist/release/`.

If no keystore is configured, the output is signed with the debug key. That is
installable but **not publishable**, and the command says so on stderr. See
[Configuration](CONFIG.md#release-signing) for how to configure signing without
committing secrets.

## `go2apk clean`

Removes `dist/`. The Android project under `android/` is left in place, since it
is regenerated on the next build.

## `go2apk doctor`

Reports which build prerequisites are present and which are missing, including
the specific SDK platform, build-tools and NDK revisions that `go2apk.yaml`
asks for. Run this first when a build fails for reasons that are not in your
code. See [Troubleshooting](TROUBLESHOOTING.md).

## `go2apk sdk install`

Writes `scripts/install-sdk.sh` and `scripts/install-sdk.ps1` for the toolchain
versions in `go2apk.yaml`. It does not download anything itself; run the script
it writes, which uses Google's official command-line tools.

## `go2apk workflow init`

Regenerates the GitHub Actions workflows from `go2apk.yaml`, without touching
anything else.

---

## Environment variables

| Variable | Effect |
| --- | --- |
| `GO2APK_GRADLE` | Path to a Gradle binary to use instead of provisioning one |
| `GO2APK_GRADLE_SHA256` | Expected checksum of the downloaded Gradle distribution |
| `GO2APK_KEYSTORE` | Keystore path, overrides `keystore` |
| `GO2APK_KEY_ALIAS` | Key alias, overrides `key_alias` |
| `GO2APK_KEYSTORE_PASSWORD` | Keystore password, overrides `keystore_password` |
| `GO2APK_KEY_PASSWORD` | Key password, overrides `key_password` |

`ANDROID_HOME` and `ANDROID_SDK_ROOT` are both honoured, in that order, before
falling back to `.go2apk/android-sdk` in the project.
