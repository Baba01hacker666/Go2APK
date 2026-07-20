package android

import (
	"fmt"

	"github.com/go2apk/go2apk/internal/config"
)

// RenderManifest creates a minimal Android manifest for generated projects.
func RenderManifest(cfg config.Config) string {
	return fmt.Sprintf(`<manifest xmlns:android="http://schemas.android.com/apk/res/android">
    <uses-permission android:name="android.permission.INTERNET" />
    <application android:extractNativeLibs="true" android:theme="%s" android:label="%s" android:allowBackup="true" android:supportsRtl="true">
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

    buildTypes {
        release {
            minifyEnabled %t
            shrinkResources %t
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }
}

tasks.register('buildGoCalculator') {
    description = 'Builds the Go calculator JNI library without gomobile.'
    group = 'build'
    def script = rootProject.file('../scripts/build-go-calculator.sh')
    inputs.dir(rootProject.file('../native/calculator'))
    outputs.dir(project.file('src/main/jniLibs'))
    doLast {
        exec {
            workingDir rootProject.file('..')
            commandLine 'bash', script.absolutePath
        }
    }
}

preBuild.dependsOn('buildGoCalculator')
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

// RenderNativeCalculator creates the Java JNI bridge for the Go calculator library.
func RenderNativeCalculator(cfg config.Config) string {
	return fmt.Sprintf(`package %s;

final class NativeCalculator {
    private static final String LIBRARY_NAME = "go2apkcalc";

    static {
        System.loadLibrary(LIBRARY_NAME);
    }

    private NativeCalculator() {
    }

    static String calculate(String left, String operator, String right) {
        return calculateWithGo(left, operator, right);
    }

    private static native String calculateWithGo(String left, String operator, String right);
}
`, cfg.Package)
}

// RenderMainActivity creates a Java calculator Activity for Gradle fallback builds.
func RenderMainActivity(cfg config.Config) string {
	return fmt.Sprintf(`package %s;

import android.app.Activity;
import android.graphics.Color;
import android.os.Bundle;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.GridLayout;
import android.widget.LinearLayout;
import android.widget.TextView;

public class MainActivity extends Activity {
    private final CalculatorEngine engine = new CalculatorEngine();
    private TextView expressionView;
    private TextView resultView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("%s");
        setContentView(createCalculatorView());
        refreshDisplay();
    }

    private View createCalculatorView() {
        LinearLayout root = new LinearLayout(this);
        root.setOrientation(LinearLayout.VERTICAL);
        root.setGravity(Gravity.CENTER_HORIZONTAL);
        root.setPadding(24, 32, 24, 24);
        root.setBackgroundColor(Color.rgb(18, 18, 18));

        expressionView = new TextView(this);
        expressionView.setGravity(Gravity.END);
        expressionView.setTextColor(Color.rgb(189, 189, 189));
        expressionView.setTextSize(22);
        expressionView.setSingleLine(false);
        root.addView(expressionView, new LinearLayout.LayoutParams(
                LinearLayout.LayoutParams.MATCH_PARENT,
                LinearLayout.LayoutParams.WRAP_CONTENT));

        resultView = new TextView(this);
        resultView.setGravity(Gravity.END);
        resultView.setTextColor(Color.WHITE);
        resultView.setTextSize(42);
        resultView.setSingleLine(false);
        root.addView(resultView, new LinearLayout.LayoutParams(
                LinearLayout.LayoutParams.MATCH_PARENT,
                0,
                1));

        GridLayout keypad = new GridLayout(this);
        keypad.setColumnCount(4);
        keypad.setRowCount(5);
        String[] keys = {
                "C", "⌫", "%%", "÷",
                "7", "8", "9", "×",
                "4", "5", "6", "−",
                "1", "2", "3", "+",
                "0", ".", "±", "="
        };
        for (String key : keys) {
            Button button = new Button(this);
            button.setText(key);
            button.setTextSize(24);
            button.setAllCaps(false);
            button.setTextColor(Color.WHITE);
            button.setBackgroundColor(isOperator(key) ? Color.rgb(255, 149, 0) : Color.rgb(51, 51, 51));
            button.setOnClickListener(v -> {
                engine.press(key);
                refreshDisplay();
            });
            GridLayout.LayoutParams params = new GridLayout.LayoutParams();
            params.width = 0;
            params.height = 144;
            params.columnSpec = GridLayout.spec(GridLayout.UNDEFINED, 1f);
            params.setMargins(6, 6, 6, 6);
            keypad.addView(button, params);
        }
        root.addView(keypad, new LinearLayout.LayoutParams(
                LinearLayout.LayoutParams.MATCH_PARENT,
                LinearLayout.LayoutParams.WRAP_CONTENT));
        return root;
    }

    private void refreshDisplay() {
        expressionView.setText(engine.expression());
        resultView.setText(engine.display());
    }

    private boolean isOperator(String key) {
        return "+−×÷=".contains(key);
    }

    static final class CalculatorEngine {
        private String accumulator;
        private String pendingOperator;
        private String currentInput = "0";
        private boolean resetInput;
        private String expression = "";

        void press(String key) {
            if ("C".equals(key)) {
                clear();
            } else if ("⌫".equals(key)) {
                backspace();
            } else if ("±".equals(key)) {
                toggleSign();
            } else if (".".equals(key)) {
                appendDecimal();
            } else if ("=".equals(key)) {
                calculate();
            } else if ("+−×÷%%".contains(key)) {
                chooseOperator(key);
            } else {
                appendDigit(key);
            }
        }

        String display() {
            return currentInput;
        }

        String expression() {
            return expression.isEmpty() ? "Calculator" : expression;
        }

        private void clear() {
            accumulator = null;
            pendingOperator = null;
            currentInput = "0";
            expression = "";
            resetInput = false;
        }

        private void appendDigit(String digit) {
            if (resetInput || "0".equals(currentInput) || isError(currentInput)) {
                currentInput = digit;
                resetInput = false;
            } else {
                currentInput += digit;
            }
        }

        private void appendDecimal() {
            if (resetInput || isError(currentInput)) {
                currentInput = "0";
                resetInput = false;
            }
            if (!currentInput.contains(".")) {
                currentInput += ".";
            }
        }

        private void backspace() {
            if (resetInput || currentInput.length() <= 1 || (currentInput.length() == 2 && currentInput.startsWith("-"))) {
                currentInput = "0";
                resetInput = false;
                return;
            }
            currentInput = currentInput.substring(0, currentInput.length() - 1);
        }

        private void toggleSign() {
            if ("0".equals(currentInput) || isError(currentInput)) {
                return;
            }
            currentInput = currentInput.startsWith("-") ? currentInput.substring(1) : "-" + currentInput;
        }

        private void chooseOperator(String operator) {
            if (pendingOperator != null && !resetInput) {
                calculate();
            } else {
                accumulator = currentInput;
            }
            if (isError(currentInput)) {
                return;
            }
            pendingOperator = operator;
            expression = accumulator + " " + operator;
            resetInput = true;
        }

        private void calculate() {
            if (pendingOperator == null || accumulator == null) {
                return;
            }
            String right = currentInput;
            String result = NativeCalculator.calculate(accumulator, pendingOperator, right);
            expression = accumulator + " " + pendingOperator + " " + right + " =";
            currentInput = result;
            accumulator = isError(result) ? null : result;
            pendingOperator = null;
            resetInput = true;
        }

        private boolean isError(String value) {
            return value.startsWith("Cannot") || value.startsWith("invalid") || value.startsWith("unsupported") || value.startsWith("Remainder");
        }
    }
}

`, cfg.Package, cfg.Name)
}
