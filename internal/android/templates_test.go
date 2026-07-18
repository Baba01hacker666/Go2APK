package android

import (
	"strings"
	"testing"

	"github.com/go2apk/go2apk/internal/config"
)

func TestRenderBuildGradleObfuscation(t *testing.T) {
	cfg := config.Default()
	cfg.Obfuscate = true
	gradle := RenderBuildGradle(cfg)
	for _, want := range []string{"minifyEnabled true", "shrinkResources true", "proguard-rules.pro"} {
		if !strings.Contains(gradle, want) {
			t.Fatalf("expected build.gradle to contain %q:\n%s", want, gradle)
		}
	}
}
