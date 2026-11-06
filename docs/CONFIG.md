# Configuration

`go2apk init` writes a commented `go2apk.yaml` in the project root. Every
command reads it, so it is the single place to change how your app is built.

Unknown keys and invalid values are **rejected**, not ignored. A typo fails
immediately with the offending line number rather than quietly producing a
project that is slightly wrong.

```yaml
name: "My Go App"
package: "com.example.goapp"
version: "0.1.0"
min_sdk: 23
target_sdk: 36
orientation: "unspecified"
theme: "@style/AppTheme"
source: "."
source_dirs: ""
gradle_dependencies: ""
build_tools: "36.0.0"
ndk: "27.2.12479018"
go_build_tags: ""
keystore: ""
keystore_password: ""
key_alias: ""
key_password: ""
obfuscate: false
```

---

## App identity

| Key | Meaning |
| --- | --- |
| `name` | Display name shown under the icon. |
| `package` | Android application ID. Must be lower-case dot-separated segments, and no segment may be a Java reserved word. This is what identifies your app on a device and in a store, so changing it after release creates a new app. |
| `version` | Version string shown to users. |

## Android levels

| Key | Meaning |
| --- | --- |
| `min_sdk` | Lowest API level that may install the app. Minimum is 21. |
| `target_sdk` | API level the app is built against. Must be at least `min_sdk`. This determines which platform the generated code is compiled for. |
| `build_tools` | Android build-tools revision installed by `scripts/install-sdk.*` and used in CI. |
| `ndk` | NDK revision used to cross-compile the Go code. |

`build_tools` and `ndk` live here so that a local SDK install, `go2apk doctor`
and the generated CI workflows all agree. `go2apk doctor` tells you if the
configured revision is not the one actually installed.

## Activity

| Key | Meaning |
| --- | --- |
| `orientation` | One of `unspecified`, `behind`, `fullSensor`, `fullUser`, `landscape`, `locked`, `nosensor`, `portrait`, `reverseLandscape`, `reversePortrait`, `sensor`, `sensorLandscape`, `sensorPortrait`, `user`, `userLandscape`, `userPortrait`. |
| `theme` | Android theme resource, for example `@style/AppTheme`. |

## Source layout

| Key | Meaning |
| --- | --- |
| `source` | Directory holding the package that calls `ui.Run` or `ui.RunApp`. |
| `source_dirs` | Comma-separated extra first-party directories to parse and to track as Gradle inputs, for example `./app/ui, ./app/logic`. |
| `go_build_tags` | Comma-separated Go build tags passed by your own build scripts. |

`source` is parsed for the widget tree. `source_dirs` exists for projects that
split the UI across packages: every listed directory is parsed, and Gradle
watches them so editing one triggers a rebuild.

## Dependencies

| Key | Meaning |
| --- | --- |
| `gradle_dependencies` | Comma-separated Android Maven artifacts, for example `androidx.appcompat:appcompat:1.7.0`. |

Go2APK adds the AndroidX artifacts each widget needs on its own, so this is only
for libraries your app uses directly.

## Output

| Key | Meaning |
| --- | --- |
| `obfuscate` | Enable R8 minification and resource shrinking on release builds. `--obfuscate` on the command line overrides it for a single build. |

---

## Release signing

An APK must be signed to be installed, and signed with a key you control to be
published.

| Key | Meaning |
| --- | --- |
| `keystore` | Path to the keystore file, relative to the project root. |
| `keystore_password` | Password for the keystore. |
| `key_alias` | Alias of the key inside the keystore. |
| `key_password` | Password for the key. |

Create a keystore once and keep it — losing it means you can never publish an
update to the app:

```bash
keytool -genkeypair -v -keystore release.jks -keyalg RSA -keysize 2048 \
    -validity 10000 -alias go2apk
```

Then reference it without putting the passwords in the file:

```yaml
keystore: "release.jks"
key_alias: "go2apk"
keystore_password: ""
key_password: ""
```

```bash
export GO2APK_KEYSTORE_PASSWORD=...
export GO2APK_KEY_PASSWORD=...
go2apk release
```

The environment takes precedence over `go2apk.yaml`, which is what lets you
commit the config and keep the secrets out of it. If a keystore and alias are
set, both passwords are required — otherwise the build fails rather than
producing an APK you cannot publish.

With no keystore configured, `go2apk release` signs with the debug key and prints
a warning to stderr. The result installs on a device but cannot be uploaded to a
store.

`release.jks` and `*.jks` should be in your `.gitignore`.

---

## Environment variables

| Variable | Effect |
| --- | --- |
| `GO2APK_KEYSTORE` | Overrides `keystore` |
| `GO2APK_KEY_ALIAS` | Overrides `key_alias` |
| `GO2APK_KEYSTORE_PASSWORD` | Overrides `keystore_password` |
| `GO2APK_KEY_PASSWORD` | Overrides `key_password` |
| `GO2APK_GRADLE` | Path to a Gradle binary to use instead of provisioning one |
| `GO2APK_GRADLE_SHA256` | Expected SHA-256 of the Gradle distribution Go2APK downloads when a project has no wrapper |
| `ANDROID_HOME`, `ANDROID_SDK_ROOT` | Where the Android SDK lives. `ANDROID_HOME` is checked first, then `ANDROID_SDK_ROOT`, then `.go2apk/android-sdk` in the project. |
