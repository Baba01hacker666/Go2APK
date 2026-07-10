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
	tools := []string{"go", "java", "gradle", "sdkmanager", "adb"}
	for _, tool := range tools {
		path, err := exec.LookPath(tool)
		if err != nil {
			fmt.Fprintf(w, "missing: %s\n", tool)
			continue
		}
		fmt.Fprintf(w, "found: %s (%s)\n", tool, path)
	}

	if sdkRoot := androidSDKRoot(); sdkRoot != "" {
		fmt.Fprintf(w, "found: Android SDK (%s)\n", sdkRoot)
	} else {
		fmt.Fprintln(w, "missing: Android SDK (set ANDROID_HOME or ANDROID_SDK_ROOT, or run go2apk sdk install)")
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
