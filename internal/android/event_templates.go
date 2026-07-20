package android

import (
	"fmt"
	"strings"

	"github.com/Baba01hacker666/Go2APK/ir"
)

// RenderEventsGen generates a Go init file that registers event handlers dynamically.
func RenderEventsGen(prog *ir.Program) string {
	var b strings.Builder
	b.WriteString("package main\n\n")

	if prog != nil && len(prog.Events) > 0 {
		b.WriteString("import \"github.com/Baba01hacker666/Go2APK/android\"\n\n")
		b.WriteString("func init() {\n")
		for eventName, goFunc := range prog.Events {
			b.WriteString(fmt.Sprintf("\tandroid.RegisterEvent(\"%s\", %s)\n", eventName, goFunc))
		}
		b.WriteString("}\n")
	} else {
		b.WriteString("// No events to register\n")
	}

	return b.String()
}
