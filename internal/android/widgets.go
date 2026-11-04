package android

import "github.com/Baba01hacker666/Go2APK/ir"

// javaWidget is the Java-oriented view of an IR widget.
//
// Every template that walks the widget tree works against this descriptor
// instead of matching on the concrete IR type. That keeps the properties that
// are genuinely per-widget (Java class name, how text is written, which
// AndroidX artifacts are required) in a single place, so a new widget kind
// cannot silently be missed by one of the generators.
type javaWidget struct {
	ID       string
	Style    ir.Style
	CSS      string
	Children []ir.Widget

	// FieldType is the Java type of the generated activity field. It is empty
	// for headless widgets that never join the view hierarchy, in which case no
	// field is declared even when the widget carries an ID.
	FieldType string

	// SetTextExpr is the Java statement (operating on the widget variable)
	// used to push new text from Go. It is empty for widgets with no text.
	SetTextExpr string

	// Readable reports whether the widget's text can be read back from Go.
	Readable bool

	// Dependencies lists AndroidX artifacts the widget cannot compile without.
	Dependencies []string
}

const cardViewDependency = "androidx.cardview:cardview:1.0.0"

// describe reports how the Java generator should treat one IR widget.
//
// The per-widget data (ID, style, CSS, children) comes from the IR helpers so
// that only the genuinely Java-specific properties are listed here.
func describe(w ir.Widget) javaWidget {
	d := javaWidget{
		ID:       ir.IDOf(w),
		Style:    ir.StyleOf(w),
		CSS:      ir.CSSOf(w),
		Children: ir.ChildrenOf(w),
	}
	switch w.(type) {
	case ir.ColumnWidget:
		d.FieldType = "LinearLayout"
	case ir.RowWidget:
		d.FieldType = "LinearLayout"
	case ir.TextViewWidget:
		d.FieldType = "TextView"
		d.SetTextExpr = "setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT))"
		d.Readable = true
	case ir.ButtonWidget:
		d.FieldType = "Button"
		d.SetTextExpr = "setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT))"
		d.Readable = true
	case ir.TextFieldWidget:
		d.FieldType = "EditText"
		d.SetTextExpr = "setText(text)"
		d.Readable = true
	case ir.ImageWidget:
		d.FieldType = "ImageView"
	case ir.AudioWidget:
		// Headless: a MediaPlayer is not a View and is never attached, which is
		// also why the IR carries no style information for it.
		d.FieldType = "android.media.MediaPlayer"
	case ir.VideoWidget:
		d.FieldType = "VideoView"
	case ir.WebViewWidget:
		d.FieldType = "android.webkit.WebView"
	case ir.XMLView:
		d.FieldType = "View"
	case ir.ScrollViewWidget:
		d.FieldType = "ScrollView"
	case ir.CardViewWidget:
		d.FieldType = "androidx.cardview.widget.CardView"
		d.Dependencies = []string{cardViewDependency}
	case ir.ProgressBarWidget:
		d.FieldType = "ProgressBar"
	case ir.SwitchWidget:
		d.FieldType = "Switch"
	}
	return d
}

// walk visits w and, depth first, every descendant.
func walk(w ir.Widget, fn func(ir.Widget)) {
	ir.Walk(w, fn)
}

// RequiredDependencies returns the AndroidX artifacts the program's widget tree
// needs in order to compile, in addition to any the user declared.
func RequiredDependencies(prog *ir.Program) []string {
	if prog == nil {
		return nil
	}
	seen := make(map[string]bool)
	var deps []string
	add := func(w ir.Widget) {
		for _, dep := range describe(w).Dependencies {
			if !seen[dep] {
				seen[dep] = true
				deps = append(deps, dep)
			}
		}
	}
	walk(prog.UI, add)
	for _, page := range prog.Pages {
		walk(page.Root, add)
	}
	return deps
}

// LayoutDefaults are the width, height and weight the generator applies when
// the app leaves the corresponding style field unset. Sizes use the same
// convention as the ui package: -1 is match-parent, -2 is wrap-content.
type LayoutDefaults struct {
	Width  int
	Height int
	Weight float32
}

// LayoutDefaultsFor returns the sizing defaults for a widget kind.
//
// The generator and the layout linter both read these values, so a warning
// always describes what the generated Java will really do instead of
// second-guessing it. A default width of 0 means the widget genuinely collapses
// unless the app sets a width.
func LayoutDefaultsFor(w ir.Widget) LayoutDefaults {
	switch w.(type) {
	case ir.ColumnWidget:
		return LayoutDefaults{Width: -1, Height: -1}
	case ir.RowWidget:
		return LayoutDefaults{Width: -1, Height: -2}
	case ir.TextViewWidget:
		return LayoutDefaults{Width: -1, Height: -2}
	case ir.ButtonWidget:
		// A Button relies on the app to pick a width; it stretches by weight
		// when it is placed in a Row and otherwise collapses.
		return LayoutDefaults{Width: 0, Height: -2, Weight: 1}
	case ir.TextFieldWidget:
		return LayoutDefaults{Width: -1, Height: -2}
	case ir.ImageWidget:
		return LayoutDefaults{Width: -1, Height: -2}
	case ir.VideoWidget:
		return LayoutDefaults{Width: -1, Height: -2}
	case ir.WebViewWidget:
		return LayoutDefaults{Width: -1, Height: -1}
	case ir.XMLView:
		return LayoutDefaults{Width: -1, Height: -2}
	case ir.ScrollViewWidget:
		return LayoutDefaults{Width: -1, Height: -1}
	case ir.CardViewWidget:
		return LayoutDefaults{Width: -1, Height: -2}
	case ir.ProgressBarWidget:
		return LayoutDefaults{Width: -2, Height: -2}
	case ir.SwitchWidget:
		return LayoutDefaults{Width: -2, Height: -2}
	}
	return LayoutDefaults{Width: -2, Height: -2}
}
