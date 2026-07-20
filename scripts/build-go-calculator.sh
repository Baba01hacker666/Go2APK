#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(pwd)"
OUT_DIR="$ROOT_DIR/android/app/src/main/jniLibs"
SRC_DIR="native/calculator"
LIB_NAME="libgo2apkcalc.so"

HOST_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
HOST_ARCH=$(uname -m)

if [ "$HOST_ARCH" = "aarch64" ] && [ -n "${PREFIX:-}" ]; then
    # We are likely in Termux, where clang is natively targeting Android
    echo "Using Termux native clang for arm64..."
    mkdir -p "$OUT_DIR/arm64-v8a"
    (cd "$SRC_DIR" && CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC=clang \
        go build -buildmode=c-shared -o "$OUT_DIR/arm64-v8a/$LIB_NAME" .)
else
    # We are on a standard desktop/CI environment, use the NDK
    SDK_ROOT="${ANDROID_HOME:-${ANDROID_SDK_ROOT:-$(pwd)/.go2apk/android-sdk}}"
    NDK_HOME=""

    if [ -d "$SDK_ROOT/ndk" ]; then
        NDK_HOME=$(find "$SDK_ROOT/ndk" -mindepth 1 -maxdepth 1 -type d | sort -r | head -n 1)
    fi

    if [ -z "$NDK_HOME" ] || [ ! -d "$NDK_HOME" ]; then
        echo "NDK not found in $SDK_ROOT/ndk. Please run go2apk sdk install or set ANDROID_HOME."
        exit 1
    fi

    # Determine NDK host tag
    if [ "$HOST_OS" = "linux" ]; then
        HOST_TAG="linux-x86_64"
    elif [ "$HOST_OS" = "darwin" ]; then
        if [ "$HOST_ARCH" = "arm64" ]; then
            HOST_TAG="darwin-aarch64"
            # Fallback for older NDKs that only had x86_64 Mac binaries
            if [ ! -d "$NDK_HOME/toolchains/llvm/prebuilt/$HOST_TAG" ]; then
                HOST_TAG="darwin-x86_64"
            fi
        else
            HOST_TAG="darwin-x86_64"
        fi
    else
        HOST_TAG="windows-x86_64"
    fi

    TOOLCHAIN="$NDK_HOME/toolchains/llvm/prebuilt/$HOST_TAG/bin"

    if [ ! -d "$TOOLCHAIN" ]; then
        echo "Toolchain not found at $TOOLCHAIN"
        exit 1
    fi

    API=24

    build_for_abi() {
        local GOOS=$1
        local GOARCH=$2
        local ABI=$3
        local CC_PREFIX=$4
        
        local CC="${TOOLCHAIN}/${CC_PREFIX}${API}-clang"
        
        echo "Building $ABI ($GOARCH)..."
        mkdir -p "$OUT_DIR/$ABI"
        
        (cd "$SRC_DIR" && CGO_ENABLED=1 GOOS=$GOOS GOARCH=$GOARCH CC="$CC" \
            go build -buildmode=c-shared -o "$OUT_DIR/$ABI/$LIB_NAME" .)
    }

    build_for_abi "android" "arm64" "arm64-v8a" "aarch64-linux-android"
    build_for_abi "android" "arm" "armeabi-v7a" "armv7a-linux-androideabi"
    build_for_abi "android" "amd64" "x86_64" "x86_64-linux-android"
    build_for_abi "android" "386" "x86" "i686-linux-android"
fi

echo "Build complete."
