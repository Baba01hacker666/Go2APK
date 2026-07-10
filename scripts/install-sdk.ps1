$ErrorActionPreference = "Stop"
$SdkRoot = if ($env:ANDROID_HOME) { $env:ANDROID_HOME } elseif ($env:ANDROID_SDK_ROOT) { $env:ANDROID_SDK_ROOT } else { Join-Path (Get-Location) ".go2apk/android-sdk" }
$ToolsVersion = if ($env:GO2APK_CMDLINE_TOOLS_VERSION) { $env:GO2APK_CMDLINE_TOOLS_VERSION } else { "11076708" }
$BuildTools = if ($env:GO2APK_BUILD_TOOLS) { $env:GO2APK_BUILD_TOOLS } else { "35.0.0" }
$Platform = if ($env:GO2APK_PLATFORM) { $env:GO2APK_PLATFORM } else { "android-35" }
$NdkVersion = if ($env:GO2APK_NDK_VERSION) { $env:GO2APK_NDK_VERSION } else { "27.2.12479018" }
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
