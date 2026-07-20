package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Baba01hacker666/Go2APK/internal/android"
	"github.com/Baba01hacker666/Go2APK/internal/config"
	"github.com/Baba01hacker666/Go2APK/internal/sdk"
	"github.com/Baba01hacker666/Go2APK/internal/workflow"
)

// Init creates the initial Go2APK configuration and Android project skeleton.
func Init(root string) error {
	cfg := config.Default()
	javaDir := filepath.ToSlash(filepath.Join("android/app/src/main/java", strings.ReplaceAll(cfg.Package, ".", "/")))
	activityPath := filepath.ToSlash(filepath.Join(javaDir, "MainActivity.java"))
	nativeBridgePath := filepath.ToSlash(filepath.Join(javaDir, "NativeBridge.java"))
	files := map[string]string{
		"go2apk.yaml": renderConfig(cfg),
		"android/settings.gradle": `pluginManagement { repositories { google(); mavenCentral(); gradlePluginPortal() } }
dependencyResolutionManagement { repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS); repositories { google(); mavenCentral() } }
rootProject.name = "Go2APKCalculator"
include ':app'
`,
		"android/build.gradle": `plugins {
    id 'com.android.application' version '8.5.2' apply false
}
`,
		"android/app/build.gradle":                 android.RenderBuildGradle(cfg),
		"android/app/proguard-rules.pro":           android.RenderProguardRules(),
		"android/app/src/main/AndroidManifest.xml": android.RenderManifest(cfg),
		activityPath:                                 android.RenderDynamicMainActivity(cfg, nil),
		nativeBridgePath:                             android.RenderNativeBridge(cfg),
		"android/app/src/main/assets/.keep":          "",
		"android/app/src/main/res/values/styles.xml": android.RenderStyles(),
		"scripts/install-sdk.sh":                     sdk.InstallScript(),
		"scripts/build-go-app.sh": `#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(pwd)"
OUT_DIR="$ROOT_DIR/android/app/src/main/jniLibs"
SRC_DIR="native/app"
LIB_NAME="libgo2apkapp.so"

HOST_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
HOST_ARCH=$(uname -m)

if [ "$HOST_ARCH" = "aarch64" ] && [ -n "${PREFIX:-}" ]; then
    echo "Using Termux native clang for arm64..."
    mkdir -p "$OUT_DIR/arm64-v8a"
    (cd "$SRC_DIR" && CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC=clang \
        go build -buildmode=c-shared -o "$OUT_DIR/arm64-v8a/$LIB_NAME" .)
else
    SDK_ROOT="${ANDROID_HOME:-${ANDROID_SDK_ROOT:-$(pwd)/.go2apk/android-sdk}}"
    NDK_HOME=""

    if [ -d "$SDK_ROOT/ndk" ]; then
        NDK_HOME=$(find "$SDK_ROOT/ndk" -mindepth 1 -maxdepth 1 -type d | sort -r | head -n 1)
    fi

    if [ -z "$NDK_HOME" ] || [ ! -d "$NDK_HOME" ]; then
        echo "NDK not found in $SDK_ROOT/ndk. Please run go2apk sdk install or set ANDROID_HOME."
        exit 1
    fi

    if [ "$HOST_OS" = "linux" ]; then
        HOST_TAG="linux-x86_64"
    elif [ "$HOST_OS" = "darwin" ]; then
        if [ "$HOST_ARCH" = "arm64" ]; then
            HOST_TAG="darwin-aarch64"
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
`,
		"scripts/install-sdk.ps1":       sdk.InstallPowerShell(),
		".github/workflows/ci.yml":      workflow.CIYAML,
		".github/workflows/release.yml": workflow.ReleaseYAML,
	}

	for name, contents := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if _, err := os.Stat(path); err == nil {
			continue
		}
		perm := os.FileMode(0o644)
		if name == "scripts/install-sdk.sh" || name == "scripts/build-go-app.sh" {
			perm = 0o755
		}
		if err := os.WriteFile(path, []byte(contents), perm); err != nil {
			return err
		}
	}
	fmt.Println("initialized Go2APK project")
	return nil
}

func renderConfig(cfg config.Config) string {
	return fmt.Sprintf("name: %s\npackage: %s\nversion: %s\nmin_sdk: %d\ntarget_sdk: %d\norientation: %s\ntheme: %s\nsource: %s\nobfuscate: %t\n", cfg.Name, cfg.Package, cfg.Version, cfg.MinSDK, cfg.TargetSDK, cfg.Orientation, cfg.Theme, cfg.Source, cfg.Obfuscate)
}
