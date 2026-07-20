package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go2apk/go2apk/internal/android"
	"github.com/go2apk/go2apk/internal/config"
	"github.com/go2apk/go2apk/internal/sdk"
	"github.com/go2apk/go2apk/internal/workflow"
)

// Init creates the initial Go2APK configuration and Android project skeleton.
func Init(root string) error {
	cfg := config.Default()
	javaDir := filepath.ToSlash(filepath.Join("android/app/src/main/java", strings.ReplaceAll(cfg.Package, ".", "/")))
	activityPath := filepath.ToSlash(filepath.Join(javaDir, "MainActivity.java"))
	nativeCalculatorPath := filepath.ToSlash(filepath.Join(javaDir, "NativeCalculator.java"))
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
		activityPath:                                 android.RenderMainActivity(cfg),
		nativeCalculatorPath:                         android.RenderNativeCalculator(cfg),
		"android/app/src/main/assets/.keep":          "",
		"android/app/src/main/res/values/styles.xml": android.RenderStyles(),
		"scripts/install-sdk.sh":                     sdk.InstallScript(),
		"scripts/build-go-calculator.sh": `#!/usr/bin/env bash
set -euo pipefail
echo "Copy scripts/build-go-calculator.sh from the Go2APK repository or run go2apk init from an up-to-date install." >&2
exit 1
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
		if name == "scripts/install-sdk.sh" || name == "scripts/build-go-calculator.sh" {
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
