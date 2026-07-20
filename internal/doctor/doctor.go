package doctor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Check reports whether common build prerequisites are available.
func Check(w io.Writer) error {
	tools := []string{"go", "java", "adb"}
	for _, tool := range tools {
		path, err := exec.LookPath(tool)
		if err != nil {
			fmt.Fprintf(w, "missing: %s\n", tool)
			continue
		}
		fmt.Fprintf(w, "found: %s (%s)\n", tool, path)
	}

	cwd, _ := os.Getwd()

	gradlewPath := filepath.Join(cwd, "android", "gradlew")
	if _, err := os.Stat(gradlewPath); err == nil {
		fmt.Fprintf(w, "found: gradle (local wrapper at %s)\n", gradlewPath)
	} else if path, err := exec.LookPath("gradle"); err == nil {
		fmt.Fprintf(w, "found: gradle (%s)\n", path)
	} else {
		fmt.Fprintf(w, "missing: gradle (neither local gradlew nor global gradle found)\n")
	}

	if sdkRoot := androidSDKRoot(); sdkRoot != "" {
		fmt.Fprintf(w, "found: Android SDK (%s)\n", sdkRoot)

		sdkmanagerPath := filepath.Join(sdkRoot, "cmdline-tools", "latest", "bin", "sdkmanager")
		if _, err := os.Stat(sdkmanagerPath); err == nil {
			fmt.Fprintf(w, "found: sdkmanager (local at %s)\n", sdkmanagerPath)
		} else if path, err := exec.LookPath("sdkmanager"); err == nil {
			fmt.Fprintf(w, "found: sdkmanager (%s)\n", path)
		} else {
			fmt.Fprintf(w, "missing: sdkmanager\n")
		}
	} else {
		fmt.Fprintln(w, "missing: Android SDK (set ANDROID_HOME or ANDROID_SDK_ROOT, or run go2apk sdk install)")
		fmt.Fprintf(w, "missing: sdkmanager\n")
	}
	return nil
}

func androidSDKRoot() string {
	for _, key := range []string{"ANDROID_HOME", "ANDROID_SDK_ROOT"} {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		local := filepath.Join(cwd, ".go2apk", "android-sdk")
		if _, err := os.Stat(local); err == nil {
			return local
		}
	}
	return ""
}
