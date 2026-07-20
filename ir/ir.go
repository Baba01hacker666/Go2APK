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

	// Permissions lists Android permissions required by this app.
	// e.g., ["android.permission.CAMERA", "android.permission.INTERNET"]
	Permissions []string

	// Receivers lists broadcast receiver registrations declared via ui.BroadcastReceiver.
	Receivers []BroadcastReceiverDecl
}

// BroadcastReceiverDecl holds the manifest + runtime info for a broadcast receiver.
type BroadcastReceiverDecl struct {
	// Name is the Go-side event name fired when this broadcast is received
	// e.g. "on_battery_low" → registered via ui.RegisterEvent
	Name string
	// Action is the Android intent action string
	// e.g. "android.intent.action.BATTERY_LOW"
	Action string
	// Exported controls android:exported in the manifest
	Exported bool
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

type Style struct {
	BackgroundColor string
	TextColor       string
	TextSize        int
	Padding         int
	Margin          int
	Width           int
	Height          int
	Weight          float32
}

type ColumnWidget struct {
	ID       string
	Style    Style
	CSS      string
	Children []Widget
}

func (ColumnWidget) WidgetType() string { return "Column" }

type RowWidget struct {
	ID       string
	Style    Style
	CSS      string
	Children []Widget
}

func (RowWidget) WidgetType() string { return "Row" }

type TextViewWidget struct {
	ID    string
	Text  string
	Style Style
	CSS   string
}

func (TextViewWidget) WidgetType() string { return "TextView" }

type ButtonWidget struct {
	ID          string
	Text        string
	Style       Style
	CSS         string
	OnClickFunc string
}

func (ButtonWidget) WidgetType() string { return "Button" }

type TextFieldWidget struct {
	ID            string
	Placeholder   string
	Style         Style
	CSS           string
	OnChangedFunc string // Function name to call when text changes
}

func (TextFieldWidget) WidgetType() string { return "TextField" }

type ImageWidget struct {
	ID    string
	Src   string
	Style Style
	CSS   string
}

func (ImageWidget) WidgetType() string { return "Image" }

type AudioWidget struct {
	ID       string
	Src      string
	AutoPlay bool
}

func (AudioWidget) WidgetType() string { return "Audio" }

type VideoWidget struct {
	ID    string
	Src   string
	Style Style
	CSS   string
}

func (VideoWidget) WidgetType() string { return "Video" }
