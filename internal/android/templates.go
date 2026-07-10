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
}
`, cfg.Package, cfg.TargetSDK, cfg.Package, cfg.MinSDK, cfg.TargetSDK, cfg.Version)
}

// RenderStyles creates the default Android theme resource.
func RenderStyles() string {
	return `<resources>
    <style name="AppTheme" parent="android:style/Theme.Material.Light.NoActionBar" />
</resources>
`
}

// RenderMainActivity creates a tiny native Android entry point for generated apps.
func RenderMainActivity(cfg config.Config) string {
	return fmt.Sprintf(`package %s;

import android.app.Activity;
import android.os.Bundle;
import android.widget.TextView;

public class MainActivity extends Activity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        TextView view = new TextView(this);
        view.setText("Hello from %s");
        view.setTextSize(24);
        view.setPadding(32, 32, 32, 32);
        setContentView(view);
    }
}
`, cfg.Package, cfg.Name)
}
