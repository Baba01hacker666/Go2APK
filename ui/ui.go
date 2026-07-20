package ui

// Widget is the base interface for all UI components.
type Widget interface {
	isWidget()
}

// Column lays out its children vertically.
type Column struct {
	ID       string
	Style    Style
	Children []Widget
}

func (Column) isWidget() {}

// Row lays out its children horizontally.
type Row struct {
	ID       string
	Style    Style
	Children []Widget
}

func (Row) isWidget() {}

// TextView displays read-only text.
type TextView struct {
	ID    string
	Text  string
	Style Style
}

func (TextView) isWidget() {}

// Button is a clickable widget.
type Button struct {
	ID      string
	Text    string
	Style   Style
	OnClick func()
}

func (Button) isWidget() {}

// TextField allows the user to enter text.
type TextField struct {
	ID          string
	Placeholder string
	Style       Style
	OnChanged   func(text string)
}

func (TextField) isWidget() {}

// Run starts the application with the given root widget.
// Note: In Go2APK, this function is statically analyzed at build time to generate
// the Android UI layout and JNI bindings, so the actual execution on Android
// happens via generated Java code.
func Run(w Widget) {
	// Stub implementation for compilation.
}

// UpdateText dynamically updates the text of a TextView widget.
func UpdateText(id string, text string) {
	updateTextNative(id, text)
}

var eventRegistry = make(map[string]func())

func RegisterEvent(name string, handler func()) {
	eventRegistry[name] = handler
}

func handleEvent(name string) {
	UpdateText("display", "Go: "+name) // DEBUG
	if handler, ok := eventRegistry[name]; ok {
		handler()
	}
}
