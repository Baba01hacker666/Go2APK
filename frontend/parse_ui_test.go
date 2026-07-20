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

func TestExtractUI(t *testing.T) {
	// Setup a temporary directory for parsing
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}

	// For a real test without a full Go module, we can use parser.Parse on a temp directory,
	// but to make it simple we'll just parse the AST directly using go/parser.
	// Since ExtractUI takes []*packages.Package, we use packages.Load overlay.

	cfg.Overlay = map[string][]byte{
		"/tmp/testapp/main.go":            []byte(testCode),
		"/tmp/testapp/ui/ui.go":           []byte(uiCode),
		"/tmp/testapp/android/android.go": []byte(androidCode),
	}

	// Load the package
	pkgs, err := packages.Load(cfg, "file=/tmp/testapp/main.go")
	if err != nil {
		t.Fatalf("packages.Load failed: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Fatalf("packages had errors")
	}

	uiTree, events, permissions, receivers, _ := ExtractUI(pkgs)

	if uiTree == nil {
		t.Fatalf("Expected root widget, got nil")
	}
	col, ok := uiTree.(ir.ColumnWidget)
	if !ok {
		t.Fatalf("Expected root to be ColumnWidget")
	}
	if col.ID != "root" {
		t.Errorf("Expected root ID to be 'root', got %q", col.ID)
	}
	if !strings.Contains(col.CSS, "background-color: #000000") {
		t.Errorf("Expected root CSS to contain background-color, got %q", col.CSS)
	}

	if len(col.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(col.Children))
	}

	tv, ok := col.Children[0].(ir.TextViewWidget)
	if !ok {
		t.Fatalf("Expected first child to be TextViewWidget")
	}
	if tv.ID != "tv1" || tv.Text != "Hello" || tv.CSS != "color: red; font-size: 24px;" {
		t.Errorf("TextView parsed incorrectly: %+v", tv)
	}

	btn, ok := col.Children[1].(ir.ButtonWidget)
	if !ok {
		t.Fatalf("Expected second child to be ButtonWidget")
	}
	if btn.ID != "btn1" || btn.Text != "Click Me" {
		t.Errorf("Button parsed incorrectly: %+v", btn)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if len(permissions) != 1 || permissions[0] != "android.permission.CAMERA" {
		t.Errorf("Permissions parsed incorrectly: %v", permissions)
	}

	if len(receivers) != 1 || receivers[0].Action != "android.intent.action.BATTERY_LOW" || receivers[0].Name != "on_battery" {
		t.Errorf("Receivers parsed incorrectly: %v", receivers)
	}
}
