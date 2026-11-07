package builder

import (
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/ir"
)

func TestLintReportsButtonWithoutWidthInColumn(t *testing.T) {
	// A Button's generated default width is 0, so inside a Column it collapses
	// to nothing. This is the case the linter exists to catch.
	prog := &ir.Program{UI: ir.ColumnWidget{Children: []ir.Widget{
		ir.ButtonWidget{ID: "go", Text: "Go"},
	}}}

	warnings := LintProgram(prog)
	if len(warnings) != 1 {
		t.Fatalf("got %d warnings, want 1: %+v", len(warnings), warnings)
	}
	if !strings.Contains(warnings[0].Message, "width is 0") {
		t.Errorf("unexpected message: %s", warnings[0].Message)
	}
	if !strings.Contains(warnings[0].Path, "Button(go)") {
		t.Errorf("path = %q, want it to identify the button", warnings[0].Path)
	}
}

func TestLintIgnoresWidgetsWithUsableDefaults(t *testing.T) {
	// A TextView's default width is match-parent, so an unset Width is fine and
	// must not be reported.
	prog := &ir.Program{UI: ir.ColumnWidget{Children: []ir.Widget{
		ir.TextViewWidget{ID: "label", Text: "hi"},
		ir.RowWidget{ID: "row"},
	}}}

	if warnings := LintProgram(prog); len(warnings) != 0 {
		t.Errorf("got %d warnings, want none: %+v", len(warnings), warnings)
	}
}

func TestLintHonoursAnExplicitWidth(t *testing.T) {
	prog := &ir.Program{UI: ir.ColumnWidget{Children: []ir.Widget{
		ir.ButtonWidget{ID: "go", Style: ir.Style{Width: -1}},
	}}}

	if warnings := LintProgram(prog); len(warnings) != 0 {
		t.Errorf("got %d warnings, want none: %+v", len(warnings), warnings)
	}
}

func TestLintRowDoesNotWarnOnTheZeroWidthPattern(t *testing.T) {
	// The 0dp + weight pattern is how buttons are normally laid out in a Row:
	// the weight stretches them, so they are not invisible.
	prog := &ir.Program{UI: ir.RowWidget{Children: []ir.Widget{
		ir.ButtonWidget{ID: "a", Style: ir.Style{Width: 0, Weight: 1}},
		ir.ButtonWidget{ID: "b"},
	}}}

	if warnings := LintProgram(prog); len(warnings) != 0 {
		t.Errorf("got %+v, want none: a Row stretches its children by weight", warnings)
	}
}

func TestLintColumnIgnoresWeight(t *testing.T) {
	// In a Column, weight distributes height, so it cannot rescue a zero width.
	prog := &ir.Program{UI: ir.ColumnWidget{Children: []ir.Widget{
		ir.ButtonWidget{ID: "a", Style: ir.Style{Width: 0, Weight: 1}},
	}}}

	if warnings := LintProgram(prog); len(warnings) != 1 {
		t.Errorf("got %+v, want one warning: weight does not affect width in a Column", warnings)
	}
}

func TestLintCoversEveryPage(t *testing.T) {
	prog := &ir.Program{Pages: []ir.Page{
		{Name: "MainPage", Root: ir.ColumnWidget{Children: []ir.Widget{ir.ButtonWidget{ID: "a"}}}},
		{Name: "SettingsPage", Root: ir.ColumnWidget{Children: []ir.Widget{ir.TextViewWidget{ID: "b"}}}},
	}}

	warnings := LintProgram(prog)
	if len(warnings) != 1 || warnings[0].Screen != "MainPage" {
		t.Fatalf("got %+v, want a single warning reported against MainPage", warnings)
	}
}

func TestLintNilProgram(t *testing.T) {
	if warnings := LintProgram(nil); warnings != nil {
		t.Errorf("LintProgram(nil) = %+v, want nil", warnings)
	}
}

func TestReportWarningsWritesToTheGivenWriter(t *testing.T) {
	var out strings.Builder
	ReportWarnings(&out, []Warning{{Screen: "MainPage", Path: "Column -> Button(go)", Message: "width is 0"}})

	got := out.String()
	if !strings.Contains(got, "warning: MainPage: Column -> Button(go)") {
		t.Errorf("unexpected output: %s", got)
	}
	if !strings.Contains(got, "1 layout warning(s)") {
		t.Errorf("expected a count, got: %s", got)
	}
}
