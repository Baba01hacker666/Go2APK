package ui

import "os"

// LoadFile reads a file from the filesystem. It is intended to load CSS and HTML
// assets into widgets. Go2APK's build system also intercepts this to embed the
// file contents during compilation.
func LoadFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

// Widget is the base interface for all UI components.
type Widget interface {
	isWidget()
}

// Column lays out its children vertically.
type Column struct {
	ID       string
	Style    Style
	CSS      string
	Children []Widget
}

func (Column) isWidget() {}

// Row lays out its children horizontally.
type Row struct {
	ID       string
	Style    Style
	CSS      string
	Children []Widget
}

func (Row) isWidget() {}

// TextView displays read-only text.
type TextView struct {
	ID    string
	Text  string
	Style Style
	CSS   string
}

func (TextView) isWidget() {}

// Button is a clickable widget.
type Button struct {
	ID      string
	Text    string
	Style   Style
	CSS     string
	OnClick func()
}

func (Button) isWidget() {}

// TextField allows the user to enter text.
type TextField struct {
	ID          string
	Placeholder string
	Style       Style
	CSS         string
	OnChanged   func(text string)
}

func (TextField) isWidget() {}

// Image displays an image from a URL or local path.
type Image struct {
	ID    string
	Src   string
	Style Style
	CSS   string
}

func (Image) isWidget() {}

// Audio plays an audio file from a URL.
type Audio struct {
	ID       string
	Src      string
	AutoPlay bool
}

func (Audio) isWidget() {}

// Video plays a video file from a URL.
type Video struct {
	ID    string
	Src   string
	Style Style
	CSS   string
}

func (Video) isWidget() {}

// WebView renders web content from a URL or raw HTML string.
type WebView struct {
	ID    string
	Src   string // URL to load
	HTML  string // Raw HTML string to load (if Src is empty)
	Style Style
	CSS   string
}

func (WebView) isWidget() {}

// ScrollView allows scrolling its child content.
type ScrollView struct {
	ID       string
	Style    Style
	CSS      string
	Children []Widget
}

func (ScrollView) isWidget() {}

// CardView provides a Material card with elevation and rounded corners.
type CardView struct {
	ID       string
	Style    Style
	CSS      string
	Children []Widget
}

func (CardView) isWidget() {}

// ProgressBar shows a loading spinner or progress bar.
type ProgressBar struct {
	ID    string
	Style Style
	CSS   string
}

func (ProgressBar) isWidget() {}

// Switch provides a toggle switch.
type Switch struct {
	ID        string
	Checked   bool
	Style     Style
	CSS       string
	OnChanged func(checked bool)
}

func (Switch) isWidget() {}

// Run starts the application with the given root widget.
// Note: In Go2APK, this function is statically analyzed at build time to generate
// the Android UI layout and JNI bindings, so the actual execution on Android
// happens via generated Java code.
func Run(w Widget) {
	// Stub implementation for compilation.
}
