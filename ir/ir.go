package ir

// Program represents the parsed and analyzed Go program,
// containing all information needed to generate the Android application.
type Program struct {
	Name       string
	Packages   map[string]*Package
	Entrypoint *Function

	// Future: UI hierarchy, Assets, Resources, JNI bindings
}

// Package represents a parsed Go package.
type Package struct {
	Name  string
	Path  string
	Types map[string]*Type
	Funcs map[string]*Function
}

// Type represents a resolved Go type.
type Type struct {
	Name string
}

// Function represents a resolved Go function.
type Function struct {
	Name string
}
