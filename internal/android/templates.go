package android

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Baba01hacker666/Go2APK/internal/config"
	"github.com/Baba01hacker666/Go2APK/ir"
)

// RenderManifest creates an Android manifest for the project, including
// any permissions and broadcast receivers declared via the ui package.
func RenderManifest(cfg config.Config, prog *ir.Program) string {
	var sb strings.Builder

	sb.WriteString(`<manifest xmlns:android="http://schemas.android.com/apk/res/android">` + "\n")

	// Always include INTERNET; add declared permissions
	internetIncluded := false
	if prog != nil {
		for _, perm := range prog.Permissions {
			if perm == "android.permission.INTERNET" {
				internetIncluded = true
			}
			sb.WriteString(fmt.Sprintf("    <uses-permission android:name=\"%s\" />\n", xmlAttr(perm)))
		}
	}
	if !internetIncluded {
		sb.WriteString("    <uses-permission android:name=\"android.permission.INTERNET\" />\n")
	}

	sb.WriteString(fmt.Sprintf(
		"    <application android:extractNativeLibs=\"true\" android:theme=\"%s\" android:label=\"%s\" android:allowBackup=\"true\" android:supportsRtl=\"true\" android:usesCleartextTraffic=\"true\">\n",
		xmlAttr(cfg.Theme), xmlAttr(cfg.Name),
	))

	// Main activity (or first page)
	sb.WriteString(fmt.Sprintf(`        <activity android:name=".MainActivity" android:exported="true" android:screenOrientation="%s">`+"\n", xmlAttr(cfg.Orientation)))
	sb.WriteString("            <intent-filter>\n")
	sb.WriteString("                <action android:name=\"android.intent.action.MAIN\" />\n")
	sb.WriteString("                <category android:name=\"android.intent.category.LAUNCHER\" />\n")
	sb.WriteString("            </intent-filter>\n")
	sb.WriteString("        </activity>\n")

	// Additional activities for multi-page apps
	if prog != nil && len(prog.Pages) > 1 {
		for i := 1; i < len(prog.Pages); i++ {
			page := prog.Pages[i]
			sb.WriteString(fmt.Sprintf(`        <activity android:name=".%s" android:exported="false" android:screenOrientation="%s"></activity>`+"\n", xmlAttr(page.Name), xmlAttr(cfg.Orientation)))
		}
	}

	// Broadcast receivers
	if prog != nil {
		for _, recv := range prog.Receivers {
			exported := "false"
			if recv.Exported {
				exported = "true"
			}
			sb.WriteString(fmt.Sprintf("        <receiver android:name=\".Go2APKBroadcastReceiver\" android:exported=\"%s\">\n", exported))
			sb.WriteString("            <intent-filter>\n")
			sb.WriteString(fmt.Sprintf("                <action android:name=\"%s\" />\n", xmlAttr(recv.Action)))
			sb.WriteString("            </intent-filter>\n")
			sb.WriteString("        </receiver>\n")
		}

		if prog.HasVPN {
			sb.WriteString("        <service android:name=\".Go2ApkVpnService\" android:permission=\"android.permission.BIND_VPN_SERVICE\" android:exported=\"false\">\n")
			sb.WriteString("            <intent-filter>\n")
			sb.WriteString("                <action android:name=\"android.net.VpnService\" />\n")
			sb.WriteString("            </intent-filter>\n")
			sb.WriteString("        </service>\n")
		}
	}

	sb.WriteString("    </application>\n")
	sb.WriteString("</manifest>\n")
	return sb.String()
}

// RenderBuildGradle creates a starter Android application Gradle file.
//
// extraDeps are Android artifacts the generated Java needs but the user did not
// declare, for example the card view library required by ui.CardView. Without
// them the generated project would not compile.
func RenderBuildGradle(cfg config.Config, extraDeps ...string) string {
	deps := mergeDependencies(cfg.GradleDependencies, extraDeps)
	return fmt.Sprintf(`plugins {
    id 'com.android.application'
}

android {
    namespace '%s'
    compileSdk %d

    defaultConfig {
        applicationId '%s'
        minSdk %d
        targetSdk %d
        versionName '%s'
        versionCode 1
    }

    packaging {
        jniLibs {
            useLegacyPackaging true
        }
    }

    sourceSets {
        main {
            jniLibs.srcDirs = ['src/main/jniLibs']
        }
    }

    compileOptions {
        sourceCompatibility JavaVersion.VERSION_17
        targetCompatibility JavaVersion.VERSION_17
    }

    signingConfigs {
%s    }

    buildTypes {
        release {
            minifyEnabled %t
            shrinkResources %t
            signingConfig signingConfigs.%s
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }
%s}
%s
tasks.register('buildGoApp', Exec) {
    description = 'Builds the Go JNI library.'
    group = 'build'
    def script = rootProject.file('../scripts/build-go-app.sh')
%s    outputs.dir(project.file('src/main/jniLibs'))
    workingDir rootProject.file('..')
    commandLine 'bash', script.absolutePath, '%s'
}

preBuild.dependsOn('buildGoApp')
`, cfg.Package, cfg.TargetSDK, cfg.Package, cfg.MinSDK, cfg.TargetSDK, cfg.Version, releaseSigningConfig(cfg), cfg.Obfuscate, cfg.Obfuscate, signingConfigName(cfg), renderDependencyMetadata(deps), renderGradleDependencies(deps), renderGradleInputs(cfg), cfg.Source)
}

// signingConfigName is the Gradle signing config the release build type uses.
// With no keystore configured it falls back to the debug key so that a release
// build still produces an installable APK instead of an "-unsigned" artifact.
func signingConfigName(cfg config.Config) string {
	if cfg.Signed() {
		return "release"
	}
	return "debug"
}

// releaseSigningConfig emits the release signing block, plus the debug fallback
// it needs so that signingConfigs.debug always resolves.
func releaseSigningConfig(cfg config.Config) string {
	var sb strings.Builder
	if cfg.Signed() {
		sb.WriteString(fmt.Sprintf("        release {\n"))
		sb.WriteString(fmt.Sprintf("            storeFile rootProject.file('../%s')\n", gradlePath(cfg.Keystore)))
		sb.WriteString(fmt.Sprintf("            storePassword %s\n", javaStringLiteral(cfg.KeystorePassword)))
		sb.WriteString(fmt.Sprintf("            keyAlias %s\n", javaStringLiteral(cfg.KeyAlias)))
		sb.WriteString(fmt.Sprintf("            keyPassword %s\n", javaStringLiteral(cfg.KeyPassword)))
		sb.WriteString("        }\n")
	}
	if !cfg.Signed() {
		// AGP creates the debug config from $HOME/.android/debug.keystore.
		sb.WriteString("        // No keystore configured; release builds fall back to the debug key.\n")
		sb.WriteString("        // Set keystore, key_alias, keystore_password and key_password in\n")
		sb.WriteString("        // go2apk.yaml (or the GO2APK_* environment variables) to ship a\n")
		sb.WriteString("        // properly signed artifact.\n")
	}
	return sb.String()
}

// gradlePath normalises a project-relative path for use with rootProject.file.
func gradlePath(path string) string {
	return filepath.ToSlash(strings.TrimPrefix(path, "./"))
}

// mergeDependencies combines the user's dependencies with those the generated
// code requires, dropping duplicates while preserving order.
func mergeDependencies(userDeps, extraDeps []string) []string {
	var out []string
	seen := make(map[string]bool, len(userDeps)+len(extraDeps))
	for _, dep := range append(append([]string(nil), userDeps...), extraDeps...) {
		if dep == "" || seen[dep] {
			continue
		}
		seen[dep] = true
		out = append(out, dep)
	}
	return out
}

func renderDependencyMetadata(deps []string) string {
	if len(deps) == 0 {
		return ""
	}
	return "\n    dependenciesInfo {\n        includeInApk true\n        includeInBundle true\n    }\n"
}

func renderGradleDependencies(deps []string) string {
	if len(deps) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\ndependencies {\n")
	for _, dep := range deps {
		sb.WriteString(fmt.Sprintf("    implementation %q\n", dep))
	}
	sb.WriteString("}\n\n")
	return sb.String()
}

func renderGradleInputs(cfg config.Config) string {
	var sb strings.Builder
	for _, src := range cfg.Sources() {
		sb.WriteString(fmt.Sprintf("    inputs.dir(rootProject.file('../%s'))\n", filepath.ToSlash(src)))
	}
	return sb.String()
}

// RenderStyles creates the default Android theme resource.
func RenderStyles() string {
	return `<resources>
    <style name="AppTheme" parent="android:style/Theme.Material.Light.NoActionBar" />
</resources>
`
}

// RenderProguardRules creates a conservative app-specific rules file for R8.
func RenderProguardRules() string {
	return `# Keep JNI bridge classes so Go can call Java methods
-keep class *.NativeBridge { *; }
-keep class *.MainActivity { *; }
-keep class *.Go2APKBroadcastReceiver { *; }
`
}
