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
