package android

import (
	"fmt"

	"github.com/go2apk/go2apk/internal/config"
)

// RenderManifest creates a minimal Android manifest for generated projects.
func RenderManifest(cfg config.Config) string {
	return fmt.Sprintf(`<manifest xmlns:android="http://schemas.android.com/apk/res/android">
    <uses-permission android:name="android.permission.INTERNET" />
    <application android:theme="%s" android:label="%s" android:allowBackup="true" android:supportsRtl="true">
        <activity android:name=".MainActivity" android:exported="true" android:screenOrientation="%s">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>
</manifest>
`, cfg.Theme, cfg.Name, cfg.Orientation)
}

// RenderBuildGradle creates a starter Android application Gradle file.
func RenderBuildGradle(cfg config.Config) string {
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

    buildTypes {
        release {
            minifyEnabled %t
            shrinkResources %t
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }
}
`, cfg.Package, cfg.TargetSDK, cfg.Package, cfg.MinSDK, cfg.TargetSDK, cfg.Version, cfg.Obfuscate, cfg.Obfuscate)
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
	return `# Add app-specific keep rules here when reflection or JNI entry points require them.
`
}
