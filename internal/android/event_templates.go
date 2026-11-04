package android

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Baba01hacker666/Go2APK/ir"
)

// RenderEventsGen generates a Go init file that registers event handlers
// dynamically.
//
// Event names are sorted so that regenerating an unchanged project produces a
// byte-identical file instead of churning on Go's randomised map order.
func RenderEventsGen(prog *ir.Program) string {
	var b strings.Builder
	b.WriteString(GeneratedFileHeader)
	b.WriteString("\n\npackage main\n\n")

	if prog != nil && len(prog.Events) > 0 {
		importPath := prog.AndroidImport
		if importPath == "" {
			importPath = "github.com/Baba01hacker666/Go2APK/android"
		}

		// Inline handlers are copied out of the app's source, so they may need
		// imports the generated file would not otherwise have.
		imports := []string{importPath}
		seen := map[string]bool{importPath: true}
		names := make([]string, 0, len(prog.Events))
		for eventName := range prog.Events {
			names = append(names, eventName)
		}
		sort.Strings(names)
		for _, eventName := range names {
			for _, imp := range prog.Events[eventName].Imports {
				if imp != "" && !seen[imp] {
					seen[imp] = true
					imports = append(imports, imp)
				}
			}
		}

		if len(imports) == 1 {
			b.WriteString(fmt.Sprintf("import %q\n\n", imports[0]))
		} else {
			b.WriteString("import (\n")
			for _, imp := range imports {
				b.WriteString(fmt.Sprintf("\t%q\n", imp))
			}
			b.WriteString(")\n\n")
		}

		b.WriteString("func init() {\n")
		for _, eventName := range names {
			handler := prog.Events[eventName]
			b.WriteString(fmt.Sprintf("\tandroid.%s(%s, %s)\n",
				registerFunc(handler.Kind), javaStringLiteral(eventName), handler.Expr))
		}
		b.WriteString("}\n")
	} else {
		b.WriteString("// This app declares no events.\n")
	}

	return b.String()
}

// registerFunc is the registration call that matches what the handler expects to
// receive. Each kind has its own entry point because the value has to cross JNI
// alongside the event name.
func registerFunc(kind ir.HandlerKind) string {
	switch kind {
	case ir.HandlerChecked:
		return "RegisterCheckedEvent"
	case ir.HandlerText:
		return "RegisterTextEvent"
	}
	return "RegisterEvent"
}
