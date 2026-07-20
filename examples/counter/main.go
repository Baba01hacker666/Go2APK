package main

import (
	"fmt"
	"github.com/Baba01hacker666/Go2APK/ui"
)

var count int = 0

// Increment handles the button click event.
// In the future transpiler, this Go function will be exposed over JNI
// or directly compiled to Java.
func Increment() {
	count++
	fmt.Printf("Count is now: %d\n", count)
	// TODO: Need a way to update UI state, e.g. ui.Update("counter_text", count)
}

func main() {
	// Declarative UI layout
	ui.Run(
		ui.Column{
			ID: "main_layout",
			Children: []ui.Widget{
				ui.TextView{
					ID:   "title",
					Text: "Counter App",
				},
				ui.TextView{
					ID:   "counter_text",
					Text: "0",
				},
				ui.Button{
					ID:      "increment_button",
					Text:    "Increment",
					OnClick: Increment,
				},
			},
		},
	)
}
