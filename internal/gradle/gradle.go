package gradle

import (
	"os"
	"path/filepath"
)

// Command returns the preferred Gradle executable for an Android project.
func Command(root string) (string, []string) {
	unixWrapper := filepath.Join(root, "android", "gradlew")
	windowsWrapper := filepath.Join(root, "android", "gradlew.bat")
	if _, err := os.Stat(unixWrapper); err == nil {
		return unixWrapper, nil
	}
	if _, err := os.Stat(windowsWrapper); err == nil {
		return windowsWrapper, nil
	}
	
	// Fallback to downloaded Gradle
	if bin, err := EnsureGradle(root); err == nil {
		return bin, nil
	}

	return "gradle", nil
}
