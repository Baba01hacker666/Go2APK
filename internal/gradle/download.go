package gradle

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func EnsureGradle(root string) (string, error) {
	gradleHome := filepath.Join(root, ".go2apk", "gradle-8.9")
	gradleBin := filepath.Join(gradleHome, "gradle-8.9", "bin", "gradle")
	if _, err := os.Stat(gradleBin); err == nil {
		return gradleBin, nil
	}

	fmt.Println("Downloading Gradle 8.9...")
	if err := os.MkdirAll(gradleHome, 0755); err != nil {
		return "", err
	}

	url := "https://services.gradle.org/distributions/gradle-8.9-bin.zip"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	zipFile := filepath.Join(gradleHome, "gradle.zip")
	out, err := os.Create(zipFile)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return "", err
	}
	defer os.Remove(zipFile)

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(gradleHome, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(gradleHome)+string(os.PathSeparator)) {
			continue
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return "", err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return "", err
		}
	}
	return gradleBin, nil
}
