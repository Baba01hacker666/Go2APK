package builder

import (
	"fmt"

	"github.com/Baba01hacker666/Go2APK/ir"
)

// LintProgram performs a static analysis of the parsed UI IR and warns the user
// about potential layout issues (e.g. invisible widgets due to width/height 0).
func LintProgram(prog *ir.Program) {
	if prog == nil {
		return
	}
	if prog.UI != nil {
		lintWidget(prog.UI, "", "Root", prog.Name)
	}
	for _, page := range prog.Pages {
		if page.Root != nil {
			lintWidget(page.Root, "", page.Name, prog.Name)
		}
	}
}

func lintWidget(w ir.Widget, parentLayout string, path string, appName string) {
	if w == nil {
		return
	}

	var style ir.Style
	var children []ir.Widget
	nodeName := w.WidgetType()

	switch v := w.(type) {
	case ir.ColumnWidget:
		style = v.Style
		children = v.Children
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.RowWidget:
		style = v.Style
		children = v.Children
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.ButtonWidget:
		style = v.Style
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.TextViewWidget:
		style = v.Style
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.TextFieldWidget:
		style = v.Style
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.ImageWidget:
		style = v.Style
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.VideoWidget:
		style = v.Style
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.SwitchWidget:
		style = v.Style
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.ProgressBarWidget:
		style = v.Style
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	case ir.ScrollViewWidget:
		style = v.Style
		if v.ID != "" {
			nodeName += "(" + v.ID + ")"
		}
	}


	currentPath := path + " -> " + nodeName

	// Warnings logic
	// In Android, our Button, TextField, etc. code generates specific defaults if style.Width/Height is 0.
	// E.g., Button defaultWidth = 0.
	// If a Button is in a Column (parentLayout == "Column"), and its Width == 0, and default is 0, it becomes invisible.
	
	if parentLayout == "Column" {
		// In a Column (vertical), if Width == 0, it collapses to 0 pixels horizontally.
		// Some widgets (like Button) might default to 0 in our generator if not specified.
		// Let's warn if Width is explicitly 0, or if it's a Button and Width isn't set.
		if style.Width == 0 {
			if w.WidgetType() == "Button" || w.WidgetType() == "Row" {
				fmt.Printf("\n[WARNING] %s:\n  Widget %s has Width=0 inside a Column. It will be completely invisible on screen!\n  Fix: Set Width to ui.MatchParent or ui.WrapContent.\n\n", appName, currentPath)
			}
		}
	} else if parentLayout == "Row" {
		// In a Row (horizontal), if Height == 0 and Weight == 0, it collapses to 0 pixels vertically.
		// Our generator sometimes sets defaultHeight = -2 (WrapContent). But if user sets Height=0 manually without Weight, it's invisible.
		if style.Height == 0 && style.Weight == 0 {
			fmt.Printf("\n[WARNING] %s:\n  Widget %s has Height=0 and Weight=0 inside a Row. It will be completely invisible on screen!\n  Fix: Set Height or add a Weight value.\n\n", appName, currentPath)
		}
	}

	for _, child := range children {
		if child != nil {
			lintWidget(child, w.WidgetType(), currentPath, appName)
		}
	}
}
