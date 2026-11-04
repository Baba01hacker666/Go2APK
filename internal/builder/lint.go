package builder

import (
	"fmt"
	"io"

	"github.com/Baba01hacker666/Go2APK/internal/android"
	"github.com/Baba01hacker666/Go2APK/ir"
)

// Warning describes a layout problem the linter found in the widget tree.
type Warning struct {
	// Screen is the page the widget belongs to, or "Root" for single-page apps.
	Screen string
	// Path locates the widget, for example "Column -> Row -> Button(btn_1)".
	Path string
	// Message explains the problem and how to fix it.
	Message string
}

// LintProgram statically analyses the parsed UI and reports widgets that
// Android would lay out with a zero-sized edge and therefore never show.
//
// The rules mirror the defaults the Java generator applies: a widget in a
// vertical Column needs a width, and one in a horizontal Row needs a height or
// a layout weight. Catching this at check time is far cheaper than debugging a
// blank screen on a device.
func LintProgram(prog *ir.Program) []Warning {
	if prog == nil {
		return nil
	}
	var warnings []Warning
	if prog.UI != nil {
		warnings = append(warnings, lintTree("Root", prog.UI)...)
	}
	for _, page := range prog.Pages {
		if page.Root != nil {
			warnings = append(warnings, lintTree(page.Name, page.Root)...)
		}
	}
	return warnings
}

func lintTree(screen string, root ir.Widget) []Warning {
	var out []Warning
	var visit func(ir.Widget, ir.Widget, string)
	visit = func(w ir.Widget, parent ir.Widget, path string) {
		if w == nil {
			return
		}
		label := w.WidgetType()
		if id := ir.IDOf(w); id != "" {
			label += "(" + id + ")"
		}
		current := label
		if path != "" {
			current = path + " -> " + label
		}

		if msg := lintWidget(w, parent); msg != "" {
			out = append(out, Warning{Screen: screen, Path: current, Message: msg})
		}
		for _, child := range ir.ChildrenOf(w) {
			visit(child, w, current)
		}
	}
	visit(root, nil, "")
	return out
}

// lintWidget returns a message when w is guaranteed to be invisible inside
// parent, and "" otherwise.
//
// Sizes are resolved through the generator's own defaults: an unset style field
// is zero in the IR, which means "use the default" rather than "collapse", so
// only the resulting size can be judged.
func lintWidget(w, parent ir.Widget) string {
	if parent == nil {
		return ""
	}
	style := ir.StyleOf(w)
	defaults := android.LayoutDefaultsFor(w)

	switch parent.WidgetType() {
	case "Column":
		// A vertical LinearLayout gives its children no intrinsic width, so a
		// resolved width of 0 collapses the widget horizontally.
		if width := resolveSize(style.Width, defaults.Width); width == 0 {
			return "resolved width is 0 inside a Column, so the widget will be invisible. " +
				"Set Width to ui.MatchParent or ui.WrapContent."
		}
	case "Row":
		// A horizontal LinearLayout gives its children no intrinsic height
		// unless a weight stretches them.
		height := resolveSize(style.Height, defaults.Height)
		weight := style.Weight
		if weight == 0 {
			weight = defaults.Weight
		}
		if height == 0 && weight == 0 {
			return "resolved height is 0 with no Weight inside a Row, so the widget will be invisible. " +
				"Set Height or give the widget a Weight."
		}
	}
	return ""
}

// resolveSize applies the generator's fallback when a style field is unset.
func resolveSize(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

// ReportWarnings writes lint findings to w in a form that is easy to read and
// to grep, and reports how many were found.
func ReportWarnings(w io.Writer, warnings []Warning) {
	for _, warn := range warnings {
		fmt.Fprintf(w, "warning: %s: %s\n    %s\n", warn.Screen, warn.Path, warn.Message)
	}
	if n := len(warnings); n > 0 {
		fmt.Fprintf(w, "%d layout warning(s)\n", n)
	}
}
