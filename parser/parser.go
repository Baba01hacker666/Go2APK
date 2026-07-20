package parser

import (
	"fmt"
	"golang.org/x/tools/go/packages"
)

// ParseConfig holds configuration for the parsing step.
type ParseConfig struct {
	Dir string
}

// Parse loads and parses the Go source code starting from the given directory.
func Parse(cfg ParseConfig) ([]*packages.Package, error) {
	mode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo

	config := &packages.Config{
		Mode: mode,
		Dir:  cfg.Dir,
	}

	pkgs, err := packages.Load(config, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}
	
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}

	return pkgs, nil
}
