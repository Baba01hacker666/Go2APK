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

	uiTree, events, permissions, receivers, hasVPN := ExtractUI(pkgs)
	prog.UI = uiTree
	prog.Events = events
	prog.Permissions = permissions
	prog.Receivers = receivers
	prog.HasVPN = hasVPN

	return prog, nil
}
