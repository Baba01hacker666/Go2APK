// Package frontend turns Go source code into the intermediate representation
// that the Android code generators consume.
package frontend

import (
	"fmt"
	"os"

	"github.com/Baba01hacker666/Go2APK/ir"
	"golang.org/x/tools/go/packages"
)

// Frontend orchestrates parsing and semantic analysis to produce the IR.
type Frontend struct{}

// New returns a ready-to-use Frontend.
func New() *Frontend {
	return &Frontend{}
}

// BuildIR parses the code at dir and produces an intermediate representation.
func (f *Frontend) BuildIR(dir string) (*ir.Program, error) {
	return f.BuildIRDirs([]string{dir})
}

// BuildIRDirs parses one or more app source roots. The first directory is
// treated as the main application package; extra directories make split
// UI/logic projects explicit for checks and Gradle inputs while normal Go
// imports keep working.
func (f *Frontend) BuildIRDirs(dirs []string) (*ir.Program, error) {
	pkgs, err := LoadPackages(dirs)
	if err != nil {
		return nil, err
	}
	ui, err := ExtractUI(pkgs)
	if err != nil {
		return nil, err
	}
	return &ir.Program{
		UI:            ui.Root,
		Pages:         ui.Pages,
		Events:        ui.Events,
		Permissions:   ui.Permissions,
		Receivers:     ui.Receivers,
		HasVPN:        ui.HasVPN,
		AndroidImport: ui.AndroidImport,
	}, nil
}

// LoadPackages loads and type-checks the package in each source root.
//
// Sources are analysed with the full syntax and type information the extractor
// needs; packages that fail to load are reported with the compiler's own
// message rather than being skipped, so a broken app fails loudly instead of
// generating an empty APK.
func LoadPackages(dirs []string) ([]*packages.Package, error) {
	const mode = packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedSyntax |
		packages.NeedTypes

	cleaned := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if dir != "" {
			cleaned = append(cleaned, dir)
		}
	}
	if len(cleaned) == 0 {
		cleaned = []string{"."}
	}

	var pkgs []*packages.Package
	for _, dir := range cleaned {
		if _, err := os.Stat(dir); err != nil {
			return nil, fmt.Errorf("source directory %s: %w", dir, err)
		}
		loaded, err := packages.Load(&packages.Config{Mode: mode, Dir: dir}, ".")
		if err != nil {
			return nil, fmt.Errorf("failed to load packages from %s: %w", dir, err)
		}
		pkgs = append(pkgs, loaded...)
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("one or more source packages failed to compile")
	}
	return pkgs, nil
}
