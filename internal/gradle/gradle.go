// Package gradle locates or provisions the Gradle build tool for a project.
package gradle

import (
	"fmt"
	"os"
	"path/filepath"
)

// Command returns the Gradle executable to use for an Android project, preferring
// the project's own wrapper so that builds are reproducible.
//
// When a project has no wrapper, Go2APK provisions a known-good Gradle release
// instead of depending on whatever version happens to be installed.
func Command(root string) (string, []string, error) {
	if name := os.Getenv("GO2APK_GRADLE"); name != "" {
		return name, nil, nil
	}

	unixWrapper := filepath.Join(root, "android", "gradlew")
	windowsWrapper := filepath.Join(root, "android", "gradlew.bat")
	if _, err := os.Stat(unixWrapper); err == nil {
		return unixWrapper, nil, nil
	}
	if _, err := os.Stat(windowsWrapper); err == nil {
		return windowsWrapper, nil, nil
	}

	bin, err := EnsureGradle(root)
	if err != nil {
		return "", nil, fmt.Errorf("no Gradle wrapper in %s and Go2APK could not provision one: %w\n"+
			"Install Gradle 8.9 or newer, or set GO2APK_GRADLE to an existing gradle binary",
			filepath.Join(root, "android"), err)
	}
	return bin, nil, nil
}
