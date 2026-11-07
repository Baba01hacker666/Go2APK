package android

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/internal/config"
	"github.com/Baba01hacker666/Go2APK/ir"
)

var (
	fieldDeclRE     = regexp.MustCompile(`(?m)^\s*private\s+\S+\s+(\w+);$`)
	fieldAssignRE   = regexp.MustCompile(`this\.(\w+)\s*=`)
	javaStringLitRE = regexp.MustCompile(`"(?:[^"\\\n]|\\.)*"`)
)

// TestRenderDynamicActivityDeclaresEveryFieldItAssigns is the regression guard
// for the bug where buildView emitted "this.id = viewN;" for widget kinds that
// declareFields did not know about, producing Java that could not compile.
func TestRenderDynamicActivityDeclaresEveryFieldItAssigns(t *testing.T) {
	tree := ir.ColumnWidget{ID: "root", Children: []ir.Widget{
		ir.ScrollViewWidget{ID: "scroll", Children: []ir.Widget{
			ir.TextViewWidget{ID: "nested", Text: "inside a scroll view"},
		}},
		ir.CardViewWidget{ID: "card", Children: []ir.Widget{
			ir.ProgressBarWidget{ID: "bar"},
		}},
		ir.WebViewWidget{ID: "web", Src: "https://example.com"},
		ir.SwitchWidget{ID: "toggle", OnChangedFunc: "onToggle"},
		ir.XMLView{ID: "legacy", Layout: "my_layout"},
	}}

	src := RenderDynamicActivity(config.Default(), "MainActivity", tree)

	declared := map[string]bool{}
	for _, m := range fieldDeclRE.FindAllStringSubmatch(src, -1) {
		declared[m[1]] = true
	}
	var assigned []string
	seen := map[string]bool{}
	for _, m := range fieldAssignRE.FindAllStringSubmatch(src, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			assigned = append(assigned, m[1])
		}
	}
	if len(assigned) == 0 {
		t.Fatal("generated activity assigns no fields; the test is not exercising the generator")
	}
	for _, name := range assigned {
		if !declared[name] {
			t.Errorf("generated Java assigns to this.%s but never declares it", name)
		}
	}
}

// TestRenderDynamicActivityEscapesText guards the bug where widget text went
// into a Java string literal unescaped, so any quote broke the build.
func TestRenderDynamicActivityEscapesText(t *testing.T) {
	tree := ir.ColumnWidget{Children: []ir.Widget{
		ir.TextViewWidget{ID: "tv", Text: `He said "hello" and left`},
		ir.ButtonWidget{ID: "btn", Text: "back\\slash", OnClickFunc: "handler"},
		ir.TextFieldWidget{ID: "tf", Placeholder: `type "here"`},
	}}

	src := RenderDynamicActivity(config.Default(), "MainActivity", tree)

	if strings.Contains(src, `fromHtml("He said "`) {
		t.Error("widget text was interpolated into the Java literal without escaping")
	}
	if !strings.Contains(src, `He said \"hello\" and left`) {
		t.Errorf("expected the quote to be escaped, got:\n%s", src)
	}
	// Every string literal in the generated source must be well formed.
	for _, line := range strings.Split(src, "\n") {
		if !strings.Contains(line, "(") {
			continue
		}
		for _, lit := range javaStringLitRE.FindAllString(line, -1) {
			if strings.HasSuffix(lit, `\`) && !strings.HasSuffix(lit, `\\`) {
				t.Errorf("unterminated string literal in: %s", line)
			}
		}
	}
}

// TestRenderDynamicActivityWiresEventsByID makes sure the event name the Java
// side dispatches matches the one the parser registers in Go.
func TestRenderDynamicActivityWiresEventsByID(t *testing.T) {
	tree := ir.ColumnWidget{Children: []ir.Widget{
		ir.ButtonWidget{ID: "go2apk_button_1", Text: "press", OnClickFunc: "HandleClick"},
		ir.SwitchWidget{ID: "sw", OnChangedFunc: "OnToggle"},
		ir.TextFieldWidget{ID: "tf", OnChangedFunc: "OnEdit"},
	}}

	src := RenderDynamicActivity(config.Default(), "MainActivity", tree)

	for _, want := range []string{
		`NativeBridge.sendEvent("go2apk_button_1_onclick")`,
		`NativeBridge.sendCheckedEvent("sw_onchanged", isChecked)`,
		`NativeBridge.sendTextEvent("tf_onchanged", s.toString())`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated Java does not dispatch %s", want)
		}
	}
}

func TestRenderManifestEscapesAndInjects(t *testing.T) {
	cfg := config.Default()
	cfg.Name = `My "App" & Co`
	prog := &ir.Program{
		Permissions: []string{"android.permission.CAMERA"},
		Receivers:   []ir.BroadcastReceiverDecl{{Name: "on_battery", Action: "android.intent.action.BATTERY_LOW"}},
		HasVPN:      true,
	}

	manifest := RenderManifest(cfg, prog)

	if strings.Contains(manifest, `My "App" & Co`) {
		t.Error("app name was not escaped for XML")
	}
	if !strings.Contains(manifest, `My &quot;App&quot; &amp; Co`) {
		t.Errorf("expected XML-escaped app name, got:\n%s", manifest)
	}
	for _, want := range []string{
		`<uses-permission android:name="android.permission.CAMERA" />`,
		`<uses-permission android:name="android.permission.INTERNET" />`,
		`<receiver android:name=".Go2APKBroadcastReceiver"`,
		`<service android:name=".Go2ApkVpnService"`,
	} {
		if !strings.Contains(manifest, want) {
			t.Errorf("manifest is missing %s", want)
		}
	}
}

func TestRenderBuildGradleIncludesGeneratedDependencies(t *testing.T) {
	cfg := config.Default()
	cfg.GradleDependencies = []string{"androidx.appcompat:appcompat:1.7.0"}

	gradle := RenderBuildGradle(cfg, cardViewDependency)

	// The user's dependency and the one the generated code needs must both be
	// present, and neither may be duplicated.
	for _, want := range []string{
		`implementation "androidx.appcompat:appcompat:1.7.0"`,
		`implementation "androidx.cardview:cardview:1.0.0"`,
	} {
		if n := strings.Count(gradle, want); n != 1 {
			t.Errorf("%s appears %d times, want 1", want, n)
		}
	}
}

func TestRenderEventsGenRegistersEveryEvent(t *testing.T) {
	prog := &ir.Program{
		Events: map[string]ir.EventHandler{
			"btn_onclick":  {Expr: "HandleClick"},
			"sw_onchanged": {Expr: "OnToggle", Imports: []string{"strconv"}},
		},
		AndroidImport: "example.com/fork/android",
	}

	src := RenderEventsGen(prog)

	for _, want := range []string{
		GeneratedFileHeader,
		`"example.com/fork/android"`,
		// An inline handler copies its imports across so it still compiles.
		`"strconv"`,
		`android.RegisterEvent("btn_onclick", HandleClick)`,
		`android.RegisterEvent("sw_onchanged", OnToggle)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("events_gen.go is missing %s:\n%s", want, src)
		}
	}
}

func TestRenderEventsGenUsesTheRegistrationMatchingTheHandler(t *testing.T) {
	prog := &ir.Program{
		AndroidImport: "example.com/fork/android",
		Events: map[string]ir.EventHandler{
			"b_onclick":   {Expr: "onClick", Kind: ir.HandlerPlain},
			"s_onchanged": {Expr: "onToggle", Kind: ir.HandlerChecked},
			"t_onchanged": {Expr: "onEdit", Kind: ir.HandlerText},
		},
	}

	src := RenderEventsGen(prog)
	for _, want := range []string{
		`android.RegisterEvent("b_onclick", onClick)`,
		`android.RegisterCheckedEvent("s_onchanged", onToggle)`,
		`android.RegisterTextEvent("t_onchanged", onEdit)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("events_gen.go is missing %s:\n%s", want, src)
		}
	}
}

// TestRenderDynamicActivityDeliversWidgetValues makes sure the Java side sends
// the value the handler expects, which is the other half of the contract.
func TestRenderDynamicActivityDeliversWidgetValues(t *testing.T) {
	tree := ir.ColumnWidget{Children: []ir.Widget{
		ir.TextFieldWidget{ID: "field", OnChangedFunc: "onEdit"},
		ir.SwitchWidget{ID: "toggle", OnChangedFunc: "onToggle"},
		ir.ButtonWidget{ID: "btn", OnClickFunc: "onClick"},
	}}

	src := RenderDynamicActivity(config.Default(), "MainActivity", tree)

	for _, want := range []string{
		`NativeBridge.sendTextEvent("field_onchanged", s.toString())`,
		`NativeBridge.sendCheckedEvent("toggle_onchanged", isChecked)`,
		`NativeBridge.sendEvent("btn_onclick")`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated Java is missing %s", want)
		}
	}
}
