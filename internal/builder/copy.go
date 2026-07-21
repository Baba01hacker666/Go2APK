package builder

import (
	"io"
	"os"
	"path/filepath"
)

// copyDir recursively copies a directory tree, attempting to preserve permissions.
func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		srcF, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcF.Close()
		dstF, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer dstF.Close()
		_, err = io.Copy(dstF, srcF)
		return err
	})
}
