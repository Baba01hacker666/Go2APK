#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ANDROID_APP_DIR="$ROOT_DIR/android/app"
NATIVE_DIR="$ROOT_DIR/native/calculator"
OUT_DIR="$ANDROID_APP_DIR/src/main/jniLibs"
API_LEVEL="${GO2APK_ANDROID_API:-23}"
NDK_HOME="${ANDROID_NDK_HOME:-${ANDROID_NDK_ROOT:-}}"

if [[ -z "$NDK_HOME" ]]; then
  SDK_ROOT="${ANDROID_HOME:-${ANDROID_SDK_ROOT:-$ROOT_DIR/android-sdk}}"
  if [[ -d "$SDK_ROOT/ndk" ]]; then
    NDK_HOME="$(find "$SDK_ROOT/ndk" -mindepth 1 -maxdepth 1 -type d | sort -V | tail -n 1)"
  fi
fi

if [[ -z "$NDK_HOME" || ! -d "$NDK_HOME" ]]; then
  echo "Android NDK not found. Set ANDROID_NDK_HOME or install the SDK/NDK with scripts/install-sdk.sh." >&2
  exit 1
fi

HOST_TAG="linux-x86_64"
case "$(uname -s)" in
  Darwin) HOST_TAG="darwin-x86_64" ;;
  Linux) HOST_TAG="linux-x86_64" ;;
esac
TOOLCHAIN="$NDK_HOME/toolchains/llvm/prebuilt/$HOST_TAG/bin"

build_abi() {
  local abi="$1" goarch="$2" cc_prefix="$3"
  mkdir -p "$OUT_DIR/$abi"
  (cd "$NATIVE_DIR" && \
    CGO_ENABLED=1 GOOS=android GOARCH="$goarch" CC="$TOOLCHAIN/${cc_prefix}${API_LEVEL}-clang" \
      go build -buildmode=c-shared -trimpath -o "$OUT_DIR/$abi/libgo2apkcalc.so" .)
}

build_abi arm64-v8a arm64 aarch64-linux-android
build_abi armeabi-v7a arm armv7a-linux-androideabi
build_abi x86 386 i686-linux-android
build_abi x86_64 amd64 x86_64-linux-android
