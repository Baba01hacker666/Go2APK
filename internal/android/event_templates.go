package android

import (
	"fmt"
	"strings"

	"github.com/Baba01hacker666/Go2APK/ir"
)

// RenderEventsGen generates a Go init file that registers event handlers dynamically.
func RenderEventsGen(prog *ir.Program) string {
	var b strings.Builder
	b.WriteString("package main\n\nimport \"github.com/Baba01hacker666/Go2APK/ui\"\n\nfunc init() {\n")
	if prog != nil && prog.Events != nil {
		for eventName, goFunc := range prog.Events {
			b.WriteString(fmt.Sprintf("\tui.RegisterEvent(\"%s\", %s)\n", eventName, goFunc))
		}
	}
	b.WriteString("}\n")
	return b.String()
}
