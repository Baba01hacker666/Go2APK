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

func TestRenderBuildGradleReleaseSigning(t *testing.T) {
	// Without a keystore the release build must still reference a signing
	// config, otherwise Gradle produces an "-unsigned" APK that cannot be
	// installed.
	cfg := config.Default()
	gradle := RenderBuildGradle(cfg)
	if !strings.Contains(gradle, "signingConfig signingConfigs.debug") {
		t.Errorf("unsigned release should fall back to the debug signing config:\n%s", gradle)
	}
	if strings.Contains(gradle, "signingConfigs.release") {
		t.Errorf("no keystore configured, but a release signing config was emitted:\n%s", gradle)
	}

	cfg.Keystore = "./release.keystore"
	cfg.KeystorePassword = "store pass"
	cfg.KeyAlias = "mykey"
	cfg.KeyPassword = "key pass"
	gradle = RenderBuildGradle(cfg)
	if !cfg.Signed() {
		t.Fatal("expected the config to report itself signed")
	}
	for _, want := range []string{
		"signingConfig signingConfigs.release",
		"storeFile rootProject.file('../release.keystore')",
		`storePassword "store pass"`,
		`keyAlias "mykey"`,
		`keyPassword "key pass"`,
	} {
		if !strings.Contains(gradle, want) {
			t.Errorf("signed release is missing %q:\n%s", want, gradle)
		}
	}
}

func TestRenderBuildGradleEscapesSigningPasswords(t *testing.T) {
	cfg := config.Default()
	cfg.Keystore = "./ks.jks"
	cfg.KeystorePassword = `pa"ss\word`
	cfg.KeyAlias = "alias"
	cfg.KeyPassword = "x"

	gradle := RenderBuildGradle(cfg)
	if strings.Contains(gradle, `pa"ss\word`) {
		t.Error("signing password was not escaped for the Gradle script")
	}
	if !strings.Contains(gradle, `storePassword "pa\"ss\\word"`) {
		t.Errorf("expected an escaped password, got:\n%s", gradle)
	}
}
