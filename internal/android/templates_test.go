package android

import (
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/internal/config"
)

func TestRenderBuildGradleObfuscation(t *testing.T) {
	cfg := config.Default()
	cfg.Obfuscate = true
	gradle := RenderBuildGradle(cfg)
	for _, want := range []string{"compileSdk 36", "targetSdk 36", "sourceCompatibility JavaVersion.VERSION_17", "minifyEnabled true", "shrinkResources true", "proguard-rules.pro"} {
		if !strings.Contains(gradle, want) {
			t.Fatalf("expected build.gradle to contain %q:\n%s", want, gradle)
		}
	}
}

func TestRenderBuildGradleDependenciesAndSourceDirs(t *testing.T) {
	cfg := config.Default()
	cfg.Source = "./app"
	cfg.SourceDirs = []string{"./app/ui", "./app/logic"}
	cfg.GradleDependencies = []string{"androidx.appcompat:appcompat:1.7.0", "com.google.android.material:material:1.12.0"}
	gradle := RenderBuildGradle(cfg)
	for _, want := range []string{
		`implementation "androidx.appcompat:appcompat:1.7.0"`,
		`implementation "com.google.android.material:material:1.12.0"`,
		`inputs.dir(rootProject.file('.././app'))`,
		`inputs.dir(rootProject.file('.././app/ui'))`,
		`inputs.dir(rootProject.file('.././app/logic'))`,
	} {
		if !strings.Contains(gradle, want) {
			t.Fatalf("expected build.gradle to contain %q:\n%s", want, gradle)
		}
	}
}
