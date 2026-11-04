package android

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Baba01hacker666/Go2APK/ir"
)

// RenderEventsGen generates a Go init file that registers event handlers dynamically.
func RenderEventsGen(prog *ir.Program) string {
	var b strings.Builder
	b.WriteString(GeneratedFileHeader)
	b.WriteString("\npackage main\n\n")

	if prog != nil && len(prog.Events) > 0 {
		importPath := prog.AndroidImport
		if importPath == "" {
			importPath = "github.com/Baba01hacker666/Go2APK/android"
		}
		b.WriteString(fmt.Sprintf("import %q\n\n", importPath))
		b.WriteString("func init() {\n")
		// Sorted so that regenerating an unchanged project produces a byte
		// identical file instead of churning on Go's randomised map order.
		names := make([]string, 0, len(prog.Events))
		for eventName := range prog.Events {
			names = append(names, eventName)
		}
		sort.Strings(names)
		for _, eventName := range names {
			b.WriteString(fmt.Sprintf("\tandroid.RegisterEvent(%s, %s)\n", javaStringLiteral(eventName), prog.Events[eventName]))
		}
		b.WriteString("}\n")
	} else {
		b.WriteString("// This app declares no events.\n")
	}

	return b.String()
}
