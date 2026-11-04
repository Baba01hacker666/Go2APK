package android

import (
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/ir"
)

// everyWidget returns one instance of each widget kind, all carrying an ID.
func everyWidget() []ir.Widget {
	return []ir.Widget{
		ir.ColumnWidget{ID: "col"},
		ir.RowWidget{ID: "row"},
		ir.TextViewWidget{ID: "tv"},
		ir.ButtonWidget{ID: "btn"},
		ir.TextFieldWidget{ID: "tf"},
		ir.ImageWidget{ID: "img"},
		ir.AudioWidget{ID: "aud"},
		ir.VideoWidget{ID: "vid"},
		ir.WebViewWidget{ID: "web"},
		ir.XMLView{ID: "xml"},
		ir.ScrollViewWidget{ID: "scroll"},
		ir.CardViewWidget{ID: "card"},
		ir.ProgressBarWidget{ID: "bar"},
		ir.SwitchWidget{ID: "sw"},
	}
}

// TestEveryWidgetDeclaresAField guards a past bug: declareFields knew only some
// widget kinds, so buildView assigned to fields that were never declared and
// the generated Activity did not compile.
func TestEveryWidgetDeclaresAField(t *testing.T) {
	for _, w := range everyWidget() {
		if got := describe(w).FieldType; got == "" {
			t.Errorf("describe(%s).FieldType is empty; a widget with an ID must declare a Java field", w.WidgetType())
		}
	}
}

func TestWalkVisitsNestedChildren(t *testing.T) {
	tree := ir.ColumnWidget{ID: "root", Children: []ir.Widget{
		ir.CardViewWidget{ID: "card", Children: []ir.Widget{
			ir.TextViewWidget{ID: "nested"},
		}},
		ir.ButtonWidget{ID: "btn"},
	}}

	var seen []string
	walk(tree, func(w ir.Widget) {
		seen = append(seen, w.WidgetType())
	})
	want := []string{"Column", "CardView", "TextView", "Button"}
	if strings.Join(seen, ",") != strings.Join(want, ",") {
		t.Errorf("walk visited %v, want %v", seen, want)
	}
}

func TestRequiredDependencies(t *testing.T) {
	// CardView needs the AndroidX artifact; nothing else does.
	prog := &ir.Program{UI: ir.CardViewWidget{ID: "card"}}
	got := RequiredDependencies(prog)
	if len(got) != 1 || got[0] != cardViewDependency {
		t.Fatalf("RequiredDependencies = %v, want [%s]", got, cardViewDependency)
	}

	prog = &ir.Program{UI: ir.ColumnWidget{Children: []ir.Widget{
		ir.CardViewWidget{ID: "a"}, ir.CardViewWidget{ID: "b"},
	}}}
	if got := RequiredDependencies(prog); len(got) != 1 {
		t.Errorf("RequiredDependencies = %v, want the dependency listed once", got)
	}

	prog = &ir.Program{UI: ir.ColumnWidget{Children: []ir.Widget{ir.TextViewWidget{ID: "t"}}}}
	if got := RequiredDependencies(prog); len(got) != 0 {
		t.Errorf("RequiredDependencies = %v, want none", got)
	}
	if got := RequiredDependencies(nil); got != nil {
		t.Errorf("RequiredDependencies(nil) = %v, want nil", got)
	}
}

func TestLayoutDefaultsFor(t *testing.T) {
	// A Button collapses unless the app gives it a width; that is why the
	// linter has to consult these defaults rather than the raw style.
	if got := LayoutDefaultsFor(ir.ButtonWidget{}); got.Width != 0 {
		t.Errorf("Button default width = %d, want 0 so an unset width is reported", got.Width)
	}
	if got := LayoutDefaultsFor(ir.TextViewWidget{}); got.Width != -1 {
		t.Errorf("TextView default width = %d, want -1 (match parent)", got.Width)
	}
	if got := LayoutDefaultsFor(ir.ColumnWidget{}); got.Height != -1 {
		t.Errorf("Column default height = %d, want -1 (match parent)", got.Height)
	}
}
