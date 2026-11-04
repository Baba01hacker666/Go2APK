// Package doctor reports whether the toolchain an Android build needs is
// present, so a failed build can be diagnosed without reading Gradle output.
package doctor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Baba01hacker666/Go2APK/internal/config"
)

// Check reports the state of each prerequisite.
//
// It returns a nil error so that the report is always printed; a missing tool is
// reported in the output rather than as a failure of the check itself.
func Check(w io.Writer, cfg config.Config) error {
	tools := []string{"go", "java", "adb"}
	for _, tool := range tools {
		path, err := exec.LookPath(tool)
		if err != nil {
			fmt.Fprintf(w, "missing: %s\n", tool)
			continue
		}
		fmt.Fprintf(w, "found: %s (%s)\n", tool, path)
	}

	gradlewPath := filepath.Join("android", "gradlew")
	if _, err := os.Stat(gradlewPath); err == nil {
		fmt.Fprintf(w, "found: gradle (wrapper at %s)\n", gradlewPath)
	} else if path, err := exec.LookPath("gradle"); err == nil {
		fmt.Fprintf(w, "found: gradle (%s)\n", path)
	} else {
		fmt.Fprintf(w, "missing: gradle (no wrapper in android/ and no gradle on PATH; "+
			"go2apk will download one, or set GO2APK_GRADLE)\n")
	}

	sdkRoot := androidSDKRoot()
	if sdkRoot == "" {
		fmt.Fprintln(w, "missing: Android SDK (set ANDROID_HOME or ANDROID_SDK_ROOT, or run go2apk sdk install)")
		return nil
	}
	fmt.Fprintf(w, "found: Android SDK (%s)\n", sdkRoot)

	if _, err := os.Stat(filepath.Join(sdkRoot, "cmdline-tools", "latest", "bin", "sdkmanager")); err == nil {
		fmt.Fprintf(w, "found: sdkmanager\n")
	} else if _, err := exec.LookPath("sdkmanager"); err == nil {
		fmt.Fprintf(w, "found: sdkmanager\n")
	} else {
		fmt.Fprintln(w, "missing: sdkmanager")
	}

	// A build fails inside Gradle if the exact platform or build-tools revision
	// the generated project asks for is not installed, which is worth seeing
	// before starting a multi-minute build.
	reportDir(w, "platform", filepath.Join(sdkRoot, "platforms"), cfg.Platform(), true)
	reportDir(w, "build-tools", filepath.Join(sdkRoot, "build-tools"), cfg.BuildTools, true)
	reportDir(w, "ndk", filepath.Join(sdkRoot, "ndk"), cfg.NDK, false)

	return nil
}

// reportDir lists the versions installed under dir and flags whether want is
// among them. When want is empty any version will do.
func reportDir(w io.Writer, label, dir, want string, exact bool) {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		if want != "" {
			fmt.Fprintf(w, "missing: %s %s (none installed in %s)\n", label, want, dir)
		} else {
			fmt.Fprintf(w, "missing: %s (none installed in %s)\n", label, dir)
		}
		return
	}

	installed := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			installed = append(installed, entry.Name())
		}
	}
	sort.Strings(installed)

	if want == "" {
		fmt.Fprintf(w, "found: %s (%s)\n", label, strings.Join(installed, ", "))
		return
	}
	for _, have := range installed {
		if have == want {
			fmt.Fprintf(w, "found: %s %s\n", label, want)
			return
		}
	}
	if exact {
		fmt.Fprintf(w, "missing: %s %s (installed: %s)\n", label, want, strings.Join(installed, ", "))
		return
	}
	fmt.Fprintf(w, "warning: %s %s not installed, falling back to %s (installed: %s)\n",
		label, want, installed[len(installed)-1], strings.Join(installed, ", "))
}

// androidSDKRoot finds the Android SDK, preferring the environment and falling
// back to the local directory go2apk sdk install prepares.
func androidSDKRoot() string {
	for _, key := range []string{"ANDROID_HOME", "ANDROID_SDK_ROOT"} {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	local := filepath.Join(".go2apk", "android-sdk")
	if _, err := os.Stat(local); err == nil {
		return local
	}
	return ""
}
