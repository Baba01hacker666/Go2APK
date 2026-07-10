package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go2apk/go2apk/internal/android"
	"github.com/go2apk/go2apk/internal/config"
)

// Init creates the initial Go2APK configuration and Android project skeleton.
func Init(root string) error {
	cfg := config.Default()
	files := map[string]string{
		"go2apk.yaml": renderConfig(cfg),
		"android/settings.gradle": `pluginManagement { repositories { google(); mavenCentral(); gradlePluginPortal() } }
dependencyResolutionManagement { repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS); repositories { google(); mavenCentral() } }
rootProject.name = "Go2APKApp"
include ':app'
`,
		"android/build.gradle": `plugins {
    id 'com.android.application' version '8.5.2' apply false
}
`,
		"android/app/build.gradle":                 android.RenderBuildGradle(cfg),
		"android/app/src/main/AndroidManifest.xml": android.RenderManifest(cfg),
		"android/app/src/main/assets/.keep":        "",
	}

	for name, contents := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if _, err := os.Stat(path); err == nil {
			continue
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			return err
		}
	}
	fmt.Println("initialized Go2APK project")
	return nil
}

func renderConfig(cfg config.Config) string {
	return fmt.Sprintf("name: %s\npackage: %s\nversion: %s\nmin_sdk: %d\ntarget_sdk: %d\norientation: %s\ntheme: %s\n", cfg.Name, cfg.Package, cfg.Version, cfg.MinSDK, cfg.TargetSDK, cfg.Orientation, cfg.Theme)
}
