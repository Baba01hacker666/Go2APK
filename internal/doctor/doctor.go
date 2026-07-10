package doctor

import (
	"fmt"
	"io"
	"os/exec"
)

// Check reports whether common build prerequisites are available.
func Check(w io.Writer) error {
	tools := []string{"go", "java", "gradle"}
	for _, tool := range tools {
		path, err := exec.LookPath(tool)
		if err != nil {
			fmt.Fprintf(w, "missing: %s\n", tool)
			continue
		}
		fmt.Fprintf(w, "found: %s (%s)\n", tool, path)
	}
	return nil
}
