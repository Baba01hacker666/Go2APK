package sdk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const (
	defaultCommandLineTools = "15859902"
	defaultBuildTools       = "36.0.0"
	defaultPlatform         = "android-36"
	defaultNDK              = "27.2.12479018"
)

// Install writes reproducible SDK setup scripts and directory placeholders.
func Install(root string, w io.Writer) error {
	sdkRoot := filepath.Join(root, ".go2apk", "android-sdk")
	for _, dir := range []string{sdkRoot, filepath.Join(root, "scripts")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "install-sdk.sh"), []byte(InstallScript()), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "install-sdk.ps1"), []byte(InstallPowerShell()), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(w, "prepared Android SDK installer for %s at %s\n", runtime.GOOS, sdkRoot)
	return nil
}

// InstallScript returns a POSIX installer for Android command-line tools.
func InstallScript() string {
	return fmt.Sprintf(`#!/usr/bin/env bash
set -euo pipefail
SDK_ROOT="${ANDROID_HOME:-${ANDROID_SDK_ROOT:-$(pwd)/.go2apk/android-sdk}}"
TOOLS_VERSION="${GO2APK_CMDLINE_TOOLS_VERSION:-%s}"
BUILD_TOOLS="${GO2APK_BUILD_TOOLS:-%s}"
PLATFORM="${GO2APK_PLATFORM:-%s}"
NDK_VERSION="${GO2APK_NDK_VERSION:-%s}"
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
`, defaultCommandLineTools, defaultBuildTools, defaultPlatform, defaultNDK)
}

// InstallPowerShell returns a Windows installer for Android command-line tools.
func InstallPowerShell() string {
	return fmt.Sprintf(`$ErrorActionPreference = "Stop"
$SdkRoot = if ($env:ANDROID_HOME) { $env:ANDROID_HOME } elseif ($env:ANDROID_SDK_ROOT) { $env:ANDROID_SDK_ROOT } else { Join-Path (Get-Location) ".go2apk/android-sdk" }
$ToolsVersion = if ($env:GO2APK_CMDLINE_TOOLS_VERSION) { $env:GO2APK_CMDLINE_TOOLS_VERSION } else { "%s" }
$BuildTools = if ($env:GO2APK_BUILD_TOOLS) { $env:GO2APK_BUILD_TOOLS } else { "%s" }
$Platform = if ($env:GO2APK_PLATFORM) { $env:GO2APK_PLATFORM } else { "%s" }
$NdkVersion = if ($env:GO2APK_NDK_VERSION) { $env:GO2APK_NDK_VERSION } else { "%s" }
New-Item -ItemType Directory -Force -Path (Join-Path $SdkRoot "cmdline-tools") | Out-Null
$Archive = "commandlinetools-win-$($ToolsVersion)_latest.zip"
$Temp = New-Item -ItemType Directory -Force -Path (Join-Path ([IO.Path]::GetTempPath()) (New-Guid))
Invoke-WebRequest "https://dl.google.com/android/repository/$Archive" -OutFile (Join-Path $Temp "tools.zip")
Expand-Archive (Join-Path $Temp "tools.zip") -DestinationPath $Temp -Force
$Latest = Join-Path $SdkRoot "cmdline-tools/latest"
Remove-Item $Latest -Recurse -Force -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $Latest | Out-Null
Copy-Item (Join-Path $Temp "cmdline-tools/*") $Latest -Recurse -Force
& (Join-Path $Latest "bin/sdkmanager.bat") --sdk_root=$SdkRoot --licenses
& (Join-Path $Latest "bin/sdkmanager.bat") --sdk_root=$SdkRoot "platform-tools" "platforms;$Platform" "build-tools;$BuildTools" "ndk;$NdkVersion"
Write-Host "ANDROID_HOME=$SdkRoot"
`, defaultCommandLineTools, defaultBuildTools, defaultPlatform, defaultNDK)
}
