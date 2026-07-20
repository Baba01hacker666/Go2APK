package ir

// Program represents the parsed and analyzed Go program,
// containing all information needed to generate the Android application.
type Program struct {
	Name       string
	Packages   map[string]*Package
	Entrypoint *Function

	// UI represents the parsed widget tree from ui.Run
	UI Widget
	// Events maps UI event handlers to Go function names
	Events map[string]string // e.g., "button_1_onclick" -> "HandleClick"
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

// Widget represents a node in the UI layout tree
type Widget interface {
	WidgetType() string
}

type ColumnWidget struct {
	ID       string
	Children []Widget
}

func (ColumnWidget) WidgetType() string { return "Column" }

type RowWidget struct {
	ID       string
	Children []Widget
}

func (RowWidget) WidgetType() string { return "Row" }

type TextViewWidget struct {
	ID   string
	Text string
}

func (TextViewWidget) WidgetType() string { return "TextView" }

type ButtonWidget struct {
	ID          string
	Text        string
	OnClickFunc string
}

func (ButtonWidget) WidgetType() string { return "Button" }

type TextFieldWidget struct {
	ID            string
	Placeholder   string
	OnChangedFunc string
}

func (TextFieldWidget) WidgetType() string { return "TextField" }
