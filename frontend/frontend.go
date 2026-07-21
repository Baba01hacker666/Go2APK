package frontend

import (
	"fmt"

	"github.com/Baba01hacker666/Go2APK/ir"
	"github.com/Baba01hacker666/Go2APK/parser"
)

// Frontend orchestrates parsing and semantic analysis to produce the IR.
type Frontend struct{}

func New() *Frontend {
	return &Frontend{}
}

// BuildIR parses the code at dir and produces an intermediate representation.
func (f *Frontend) BuildIR(dir string) (*ir.Program, error) {
	pkgs, err := parser.Parse(parser.ParseConfig{Dir: dir})
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	fmt.Printf("Parsed %d packages in dir %s\n", len(pkgs), dir)

	prog := &ir.Program{
		Packages: make(map[string]*ir.Package),
	}

	for _, pkg := range pkgs {
		irPkg := &ir.Package{
			Name:  pkg.Name,
			Path:  pkg.PkgPath,
			Types: make(map[string]*ir.Type),
			Funcs: make(map[string]*ir.Function),
		}

		fmt.Printf("Package %s has %d syntax trees\n", pkg.Name, len(pkg.Syntax))

		// TODO: Traverse AST/Types to populate IR structs and detect entrypoint

		prog.Packages[pkg.PkgPath] = irPkg
	}

	androidImport := "github.com/Baba01hacker666/Go2APK/android"
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, imp := range file.Imports {
				importPath := imp.Path.Value[1 : len(imp.Path.Value)-1]
				if len(importPath) > 8 && importPath[len(importPath)-8:] == "/android" && (len(importPath) > 16 && importPath[len(importPath)-15:] == "Go2APK/android") {
					androidImport = importPath
				}
			}
		}
	}

	uiTree, pages, events, permissions, receivers, hasVPN := ExtractUI(pkgs)
	prog.UI = uiTree
	prog.Pages = pages
	prog.Events = events
	prog.Permissions = permissions
	prog.Receivers = receivers
	prog.HasVPN = hasVPN
	prog.AndroidImport = androidImport

	return prog, nil
}
