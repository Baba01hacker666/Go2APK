package frontend

import (
	"fmt"

	"github.com/go2apk/go2apk/ir"
	"github.com/go2apk/go2apk/parser"
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
		
		// TODO: Traverse AST/Types to populate IR structs and detect entrypoint
		
		prog.Packages[pkg.PkgPath] = irPkg
	}

	return prog, nil
}
