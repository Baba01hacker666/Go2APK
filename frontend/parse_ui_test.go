package frontend

import (
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/ir"
	"golang.org/x/tools/go/packages"
)

const testCode = `package main

import (
	"github.com/Baba01hacker666/Go2APK/ui"
	"github.com/Baba01hacker666/Go2APK/android"
)

func myHandler() {}

func main() {
	android.Permission("android.permission.CAMERA")
	android.BroadcastReceiver("android.intent.action.BATTERY_LOW", "on_battery", func() {})

	ui.Run(ui.Column{
		ID: "root",
		CSS: "background-color: #000000; padding: 16px;",
		Children: []ui.Widget{
			ui.TextView{
				ID: "tv1",
				Text: "Hello",
				CSS: "color: red; font-size: 24px;",
			},
			ui.Button{
				ID: "btn1",
				Text: "Click Me",
				OnClick: myHandler,
			},
		},
	})
}
`

// A minimal mock ui package to satisfy the parser.
const uiCode = `package ui
type Widget interface{}
type Column struct { ID, CSS string; Children []Widget }
type TextView struct { ID, Text, CSS string }
type Button struct { ID, Text, CSS string; OnClick func() }
func Run(w Widget) {}
`

const androidCode = `package android
func Permission(name string) {}
func BroadcastReceiver(action, name string, cb func()) {}
`

// loadTestProgram parses testCode with the mock ui and android packages, as if
// it were the entry point of a real app.
func loadTestProgram(t *testing.T, code string) []*packages.Package {
	t.Helper()
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedSyntax | packages.NeedTypes,
		Overlay: map[string][]byte{
			"/go2apk-test/main.go":            []byte(code),
			"/go2apk-test/ui/ui.go":           []byte(uiCode),
			"/go2apk-test/android/android.go": []byte(androidCode),
		},
	}
	pkgs, err := packages.Load(cfg, "file=/go2apk-test/main.go")
	if err != nil {
		t.Fatalf("packages.Load failed: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Fatalf("packages had errors")
	}
	return pkgs
}

func TestExtractUI(t *testing.T) {
	ui, err := ExtractUI(loadTestProgram(t, testCode))
	if err != nil {
		t.Fatalf("ExtractUI failed: %v", err)
	}

	if ui.Root == nil {
		t.Fatal("expected a root widget, got nil")
	}
	col, ok := ui.Root.(ir.ColumnWidget)
	if !ok {
		t.Fatalf("expected root to be ColumnWidget, got %T", ui.Root)
	}
	if col.ID != "root" {
		t.Errorf("root ID = %q, want %q", col.ID, "root")
	}
	if !strings.Contains(col.CSS, "background-color: #000000") {
		t.Errorf("root CSS = %q, want it to contain the background colour", col.CSS)
	}
	if len(col.Children) != 2 {
		t.Fatalf("got %d children, want 2", len(col.Children))
	}

	tv, ok := col.Children[0].(ir.TextViewWidget)
	if !ok {
		t.Fatalf("first child = %T, want TextViewWidget", col.Children[0])
	}
	if tv.ID != "tv1" || tv.Text != "Hello" || tv.CSS != "color: red; font-size: 24px;" {
		t.Errorf("TextView parsed incorrectly: %+v", tv)
	}

	btn, ok := col.Children[1].(ir.ButtonWidget)
	if !ok {
		t.Fatalf("second child = %T, want ButtonWidget", col.Children[1])
	}
	if btn.ID != "btn1" || btn.Text != "Click Me" {
		t.Errorf("Button parsed incorrectly: %+v", btn)
	}

	if got, want := ui.Events["btn1_onclick"].Expr, "myHandler"; got != want {
		t.Errorf("event btn1_onclick = %q, want %q", got, want)
	}
	if len(ui.Events) != 1 {
		t.Errorf("got %d events, want 1: %v", len(ui.Events), ui.Events)
	}

	if len(ui.Permissions) != 1 || ui.Permissions[0] != "android.permission.CAMERA" {
		t.Errorf("permissions = %v, want [android.permission.CAMERA]", ui.Permissions)
	}
	if len(ui.Receivers) != 1 {
		t.Fatalf("got %d receivers, want 1", len(ui.Receivers))
	}
	if ui.Receivers[0].Action != "android.intent.action.BATTERY_LOW" || ui.Receivers[0].Name != "on_battery" {
		t.Errorf("receiver parsed incorrectly: %+v", ui.Receivers[0])
	}
	if ui.AndroidImport != DefaultAndroidImport {
		t.Errorf("AndroidImport = %q, want %q", ui.AndroidImport, DefaultAndroidImport)
	}
}

func TestExtractUIGeneratesIDForHandlerWithoutOne(t *testing.T) {
	const code = `package main

import "github.com/Baba01hacker666/Go2APK/ui"

func handler() {}

func main() {
	ui.Run(ui.Column{Children: []ui.Widget{
		ui.Button{Text: "no id", OnClick: handler},
	}})
}
`
	ui, err := ExtractUI(loadTestProgram(t, code))
	if err != nil {
		t.Fatalf("ExtractUI failed: %v", err)
	}
	col, ok := ui.Root.(ir.ColumnWidget)
	if !ok {
		t.Fatalf("root = %T, want ColumnWidget", ui.Root)
	}
	btn, ok := col.Children[0].(ir.ButtonWidget)
	if !ok {
		t.Fatalf("child = %T, want ButtonWidget", col.Children[0])
	}
	if btn.ID == "" {
		t.Fatal("expected a generated ID for a Button with a handler")
	}
	if got, want := ui.Events[btn.ID+"_onclick"].Expr, "handler"; got != want {
		t.Errorf("event %q = %q, want %q; events = %v", btn.ID+"_onclick", got, want, ui.Events)
	}
}

func TestExtractUIRejectsInvalidWidgetID(t *testing.T) {
	const code = `package main

import "github.com/Baba01hacker666/Go2APK/ui"

func main() {
	ui.Run(ui.Column{Children: []ui.Widget{
		ui.TextView{ID: "not-a-java-name", Text: "hi"},
	}})
}
`
	_, err := ExtractUI(loadTestProgram(t, code))
	if err == nil {
		t.Fatal("expected an error for an ID that is not a Java identifier")
	}
	if !strings.Contains(err.Error(), "not a valid Java identifier") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExtractUIMultiPage(t *testing.T) {
	const code = `package main

import "github.com/Baba01hacker666/Go2APK/ui"

func main() {
	ui.RunApp(
		ui.Page{Name: "MainPage", Root: ui.TextView{ID: "a", Text: "one"}},
		ui.Page{Name: "SettingsPage", Root: ui.TextView{ID: "b", Text: "two"}},
	)
}
`
	uiCode := uiCode + `
type Page struct { Name string; Root Widget }
func RunApp(pages ...Page) {}
`
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedImports,
		Overlay: map[string][]byte{
			"/go2apk-test/main.go":  []byte(code),
			"/go2apk-test/ui/ui.go": []byte(uiCode),
		},
	}
	pkgs, err := packages.Load(cfg, "file=/go2apk-test/main.go")
	if err != nil {
		t.Fatalf("packages.Load failed: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Fatal("packages had errors")
	}
	prog, err := ExtractUI(pkgs)
	if err != nil {
		t.Fatalf("ExtractUI failed: %v", err)
	}
	if len(prog.Pages) != 2 {
		t.Fatalf("got %d pages, want 2", len(prog.Pages))
	}
	if prog.Pages[0].Name != "MainPage" || prog.Pages[1].Name != "SettingsPage" {
		t.Errorf("pages parsed incorrectly: %+v", prog.Pages)
	}
}

func TestExtractUIAcceptsInlineFuncHandlers(t *testing.T) {
	const code = `package main

import (
	"github.com/Baba01hacker666/Go2APK/android"
	"github.com/Baba01hacker666/Go2APK/ui"
)

func main() {
	ui.Run(ui.Column{Children: []ui.Widget{
		ui.Button{ID: "btn", Text: "go", OnClick: func() {
			android.UpdateText("btn", "done")
		}},
	}})
}
`
	prog, err := ExtractUI(loadTestProgram(t, code))
	if err != nil {
		t.Fatalf("ExtractUI failed: %v", err)
	}
	got := prog.Events["btn_onclick"].Expr
	// The callback is kept as a Go expression so that it can be emitted into
	// events_gen.go as written.
	if !strings.Contains(got, "func()") || !strings.Contains(got, `android.UpdateText("btn", "done")`) {
		t.Errorf("event handler = %q, want the func literal source", got)
	}
}

func TestExtractUIRejectsInvalidHandler(t *testing.T) {
	const code = `package main

import "github.com/Baba01hacker666/Go2APK/ui"

func makeHandler() func() { return func() {} }

func main() {
	ui.Run(ui.Column{Children: []ui.Widget{
		ui.Button{ID: "btn", Text: "go", OnClick: makeHandler()},
	}})
}
`
	_, err := ExtractUI(loadTestProgram(t, code))
	if err == nil {
		t.Fatal("expected an error when OnClick is not a function")
	}
	if !strings.Contains(err.Error(), "OnClick") {
		t.Errorf("unexpected error: %v", err)
	}
}
