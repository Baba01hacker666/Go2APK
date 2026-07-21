package parser

import (
	"fmt"
	"golang.org/x/tools/go/packages"
)

// ParseConfig holds configuration for the parsing step.
type ParseConfig struct {
	Dir  string
	Dirs []string
}

// Parse loads and parses the Go source code starting from the given directory.
func Parse(cfg ParseConfig) ([]*packages.Package, error) {
	mode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedSyntax

	dirs := cfg.Dirs
	if len(dirs) == 0 {
		dirs = []string{cfg.Dir}
	}
	if len(dirs) == 0 || dirs[0] == "" {
		dirs = []string{"."}
	}

	var pkgs []*packages.Package
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		config := &packages.Config{
			Mode: mode,
			Dir:  dir,
		}

		loaded, err := packages.Load(config, ".")
		if err != nil {
			return nil, fmt.Errorf("failed to load packages from %s: %w", dir, err)
		}
		pkgs = append(pkgs, loaded...)
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}

	return pkgs, nil
}
