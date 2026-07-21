package ui

// Component is a small reusable function that returns a widget tree. It is a
// convention for separating UI composition from business logic across files.
type Component func() Widget

// Build evaluates a component and returns its widget tree. It keeps app entry
// points concise when the UI is split into feature files.
func Build(component Component) Widget {
	if component == nil {
		return TextView{}
	}
	return component()
}

// V creates a vertical layout with children.
func V(children ...Widget) Column { return Column{Children: children} }

// H creates a horizontal layout with children.
func H(children ...Widget) Row { return Row{Children: children} }

// Screen creates a named page from a component, making multi-file routing easier.
func Screen(name string, component Component) Page {
	return Page{Name: name, Root: Build(component)}
}
