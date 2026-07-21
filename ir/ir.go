package ir

// Program represents the parsed and analyzed Go program,
// containing all information needed to generate the Android application.
type Program struct {
	Name       string
	Packages   map[string]*Package
	Entrypoint *Function

	// UI represents the parsed widget tree from ui.Run (legacy single-page)
	UI Widget
	// Pages represents the parsed widget trees from ui.RunApp (multi-page)
	Pages []Page
	// Events maps UI event handlers to Go function names
	Events map[string]string // e.g., "button_1_onclick" -> "HandleClick"

	// Permissions lists Android permissions required by this app.
	// e.g., ["android.permission.CAMERA", "android.permission.INTERNET"]
	Permissions []string

	// Receivers lists broadcast receiver registrations declared via ui.BroadcastReceiver.
	Receivers []BroadcastReceiverDecl

	// HasVPN indicates if the program uses android.StartVPN, requiring VpnService setup.
	HasVPN bool

	AndroidImport string
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

// Page represents a parsed screen in a multi-page app.
type Page struct {
	Name string
	Root Widget
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

type Animation struct {
	Type     string
	Duration int
	Delay    int
	Loop     bool
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
	Animation       Animation
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

// XMLView represents a native Android XML layout
type XMLView struct {
	ID     string
	Layout string
	Style  Style
	CSS    string
}

func (x XMLView) WidgetType() string { return "XMLView" }

// VideoWidget represents an video player.
type VideoWidget struct {
	ID    string
	Src   string
	Style Style
	CSS   string
}

func (VideoWidget) WidgetType() string { return "Video" }

// WebViewWidget represents a web view.
type WebViewWidget struct {
	ID    string
	Src   string
	HTML  string
	Style Style
	CSS   string
}

func (WebViewWidget) WidgetType() string { return "WebView" }

// ScrollViewWidget represents a scrolling container.
type ScrollViewWidget struct {
	ID       string
	Style    Style
	CSS      string
	Children []Widget
}

func (ScrollViewWidget) WidgetType() string { return "ScrollView" }

// CardViewWidget represents a material design card.
type CardViewWidget struct {
	ID       string
	Style    Style
	CSS      string
	Children []Widget
}

func (CardViewWidget) WidgetType() string { return "CardView" }

// ProgressBarWidget represents a loading indicator.
type ProgressBarWidget struct {
	ID    string
	Style Style
	CSS   string
}

func (ProgressBarWidget) WidgetType() string { return "ProgressBar" }

// SwitchWidget represents a toggle switch.
type SwitchWidget struct {
	ID            string
	Checked       bool
	Style         Style
	CSS           string
	OnChangedFunc string
}

func (SwitchWidget) WidgetType() string { return "Switch" }
