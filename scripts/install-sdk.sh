#!/usr/bin/env bash
set -euo pipefail
SDK_ROOT="${ANDROID_HOME:-${ANDROID_SDK_ROOT:-$(pwd)/.go2apk/android-sdk}}"
TOOLS_VERSION="${GO2APK_CMDLINE_TOOLS_VERSION:-15859902}"
BUILD_TOOLS="${GO2APK_BUILD_TOOLS:-36.0.0}"
PLATFORM="${GO2APK_PLATFORM:-android-36}"
NDK_VERSION="${GO2APK_NDK_VERSION:-27.2.12479018}"
mkdir -p "$SDK_ROOT/cmdline-tools"
case "$(uname -s)" in
  Darwin) archive="commandlinetools-mac-${TOOLS_VERSION}_latest.zip" ;;
  Linux) archive="commandlinetools-linux-${TOOLS_VERSION}_latest.zip" ;;
  *) echo "Unsupported OS for automatic Android SDK installation" >&2; exit 1 ;;
esac
url="https://dl.google.com/android/repository/${archive}"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
curl -fsSL "$url" -o "$tmp/tools.zip"
unzip -q "$tmp/tools.zip" -d "$tmp"
rm -rf "$SDK_ROOT/cmdline-tools/latest"
mkdir -p "$SDK_ROOT/cmdline-tools/latest"
cp -R "$tmp/cmdline-tools/." "$SDK_ROOT/cmdline-tools/latest/"
yes | "$SDK_ROOT/cmdline-tools/latest/bin/sdkmanager" --sdk_root="$SDK_ROOT" --licenses >/dev/null
yes | "$SDK_ROOT/cmdline-tools/latest/bin/sdkmanager" --sdk_root="$SDK_ROOT" \
  "platform-tools" "platforms;$PLATFORM" "build-tools;$BUILD_TOOLS" "ndk;$NDK_VERSION"
echo "export ANDROID_HOME=$SDK_ROOT"
echo "export ANDROID_SDK_ROOT=$SDK_ROOT"
