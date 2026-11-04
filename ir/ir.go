package ir

// Program represents the parsed and analyzed Go program,
// containing all information needed to generate the Android application.
type Program struct {
	// UI represents the parsed widget tree from ui.Run (legacy single-page)
	UI Widget
	// Pages represents the parsed widget trees from ui.RunApp (multi-page)
	Pages []Page
	// Events maps generated Android event names to their Go handlers,
	// e.g. "button_1_onclick" -> HandleClick.
	Events map[string]EventHandler

	// Permissions lists Android permissions required by this app.
	// e.g., ["android.permission.CAMERA", "android.permission.INTERNET"]
	Permissions []string

	// Receivers lists broadcast receiver registrations declared via ui.BroadcastReceiver.
	Receivers []BroadcastReceiverDecl

	// HasVPN indicates if the program uses android.StartVPN, requiring VpnService setup.
	HasVPN bool

	AndroidImport string
}

// EventHandler is a Go callback bound to an Android event name.
type EventHandler struct {
	// Expr is Go source for the handler: the name of a function, or a func
	// literal copied out of the app's source.
	Expr string
	// Imports are the import paths Expr needs. An inline handler is emitted
	// into a generated file, which does not inherit the imports of the file it
	// was written in.
	Imports []string
	// Kind is what the handler expects to receive, and therefore which JNI
	// entry point can deliver the event.
	Kind HandlerKind
}

// HandlerKind describes the value a Go event handler is called with. It is
// decided by the widget's field type: ui.Button uses func(), ui.Switch
// func(bool) and ui.TextField func(string).
type HandlerKind int

const (
	// HandlerPlain is a func() handler, which receives nothing.
	HandlerPlain HandlerKind = iota
	// HandlerChecked is a func(bool) handler, which receives the new state of
	// a ui.Switch.
	HandlerChecked
	// HandlerText is a func(string) handler, which receives the new contents of
	// a ui.TextField.
	HandlerText
)

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

// Widget represents a node in the UI layout tree
type Widget interface {
	WidgetType() string
}

// ChildrenOf returns the child widgets of w, or nil when w is a leaf.
//
// It is intentionally defined next to the widget definitions: adding a widget
// means updating this switch rather than hunting for every tree walk scattered
// across the generators and linters.
func ChildrenOf(w Widget) []Widget {
	switch v := w.(type) {
	case ColumnWidget:
		return v.Children
	case RowWidget:
		return v.Children
	case ScrollViewWidget:
		return v.Children
	case CardViewWidget:
		return v.Children
	}
	return nil
}

// IDOf returns the identifier of w, or "" when the widget has none.
func IDOf(w Widget) string {
	switch v := w.(type) {
	case ColumnWidget:
		return v.ID
	case RowWidget:
		return v.ID
	case TextViewWidget:
		return v.ID
	case ButtonWidget:
		return v.ID
	case TextFieldWidget:
		return v.ID
	case ImageWidget:
		return v.ID
	case AudioWidget:
		return v.ID
	case XMLView:
		return v.ID
	case VideoWidget:
		return v.ID
	case WebViewWidget:
		return v.ID
	case ScrollViewWidget:
		return v.ID
	case CardViewWidget:
		return v.ID
	case ProgressBarWidget:
		return v.ID
	case SwitchWidget:
		return v.ID
	}
	return ""
}

// StyleOf returns the styling of w. Widgets without styling return the zero
// value, which every consumer treats as "use the generator's default".
func StyleOf(w Widget) Style {
	switch v := w.(type) {
	case ColumnWidget:
		return v.Style
	case RowWidget:
		return v.Style
	case TextViewWidget:
		return v.Style
	case ButtonWidget:
		return v.Style
	case TextFieldWidget:
		return v.Style
	case ImageWidget:
		return v.Style
	case XMLView:
		return v.Style
	case VideoWidget:
		return v.Style
	case WebViewWidget:
		return v.Style
	case ScrollViewWidget:
		return v.Style
	case CardViewWidget:
		return v.Style
	case ProgressBarWidget:
		return v.Style
	case SwitchWidget:
		return v.Style
	}
	return Style{}
}

// CSSOf returns the inline CSS of w, or "" when the widget carries none.
func CSSOf(w Widget) string {
	switch v := w.(type) {
	case ColumnWidget:
		return v.CSS
	case RowWidget:
		return v.CSS
	case TextViewWidget:
		return v.CSS
	case ButtonWidget:
		return v.CSS
	case TextFieldWidget:
		return v.CSS
	case ImageWidget:
		return v.CSS
	case XMLView:
		return v.CSS
	case VideoWidget:
		return v.CSS
	case WebViewWidget:
		return v.CSS
	case ScrollViewWidget:
		return v.CSS
	case CardViewWidget:
		return v.CSS
	case ProgressBarWidget:
		return v.CSS
	case SwitchWidget:
		return v.CSS
	}
	return ""
}

// Walk visits w and, depth first, each of its descendants.
func Walk(w Widget, fn func(Widget)) {
	if w == nil {
		return
	}
	fn(w)
	for _, child := range ChildrenOf(w) {
		Walk(child, fn)
	}
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
