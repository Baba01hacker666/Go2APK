package main

import (
	"fmt"
	"strconv"

	"github.com/Baba01hacker666/Go2APK/ui"
)

var (
	currentInput  string
	previousInput string
	operator      string
)

func main() {
	ui.Run(
		ui.Column{
			ID: "main_layout",
			Children: []ui.Widget{
				ui.TextView{
					ID:   "display",
					Text: "0",
				},
				ui.Row{
					ID: "row_1",
					Children: []ui.Widget{
						ui.Button{
							ID:   "btn_7",
							Text: "7",
							OnClick: func() {
								appendDigit("7")
							},
						},
						ui.Button{
							ID:   "btn_8",
							Text: "8",
							OnClick: func() {
								appendDigit("8")
							},
						},
						ui.Button{
							ID:   "btn_9",
							Text: "9",
							OnClick: func() {
								appendDigit("9")
							},
						},
						ui.Button{
							ID:   "btn_div",
							Text: "/",
							OnClick: func() {
								setOperator("/")
							},
						},
					},
				},
				ui.Row{
					ID: "row_2",
					Children: []ui.Widget{
						ui.Button{
							ID:   "btn_4",
							Text: "4",
							OnClick: func() {
								appendDigit("4")
							},
						},
						ui.Button{
							ID:   "btn_5",
							Text: "5",
							OnClick: func() {
								appendDigit("5")
							},
						},
						ui.Button{
							ID:   "btn_6",
							Text: "6",
							OnClick: func() {
								appendDigit("6")
							},
						},
						ui.Button{
							ID:   "btn_mul",
							Text: "*",
							OnClick: func() {
								setOperator("*")
							},
						},
					},
				},
				ui.Row{
					ID: "row_3",
					Children: []ui.Widget{
						ui.Button{
							ID:   "btn_1",
							Text: "1",
							OnClick: func() {
								appendDigit("1")
							},
						},
						ui.Button{
							ID:   "btn_2",
							Text: "2",
							OnClick: func() {
								appendDigit("2")
							},
						},
						ui.Button{
							ID:   "btn_3",
							Text: "3",
							OnClick: func() {
								appendDigit("3")
							},
						},
						ui.Button{
							ID:   "btn_sub",
							Text: "-",
							OnClick: func() {
								setOperator("-")
							},
						},
					},
				},
				ui.Row{
					ID: "row_4",
					Children: []ui.Widget{
						ui.Button{
							ID:   "btn_clear",
							Text: "C",
							OnClick: func() {
								currentInput = ""
								previousInput = ""
								operator = ""
								updateDisplay("0")
							},
						},
						ui.Button{
							ID:   "btn_0",
							Text: "0",
							OnClick: func() {
								appendDigit("0")
							},
						},
						ui.Button{
							ID:   "btn_eq",
							Text: "=",
							OnClick: func() {
								calculate()
							},
						},
						ui.Button{
							ID:   "btn_add",
							Text: "+",
							OnClick: func() {
								setOperator("+")
							},
						},
					},
				},
			},
		},
	)
}

func appendDigit(d string) {
	currentInput += d
	updateDisplay(currentInput)
}

func setOperator(op string) {
	if currentInput != "" {
		if previousInput != "" {
			calculate()
		} else {
			previousInput = currentInput
			currentInput = ""
		}
	}
	operator = op
}

func calculate() {
	if previousInput == "" || currentInput == "" || operator == "" {
		return
	}
	p, err1 := strconv.ParseFloat(previousInput, 64)
	c, err2 := strconv.ParseFloat(currentInput, 64)
	if err1 != nil || err2 != nil {
		updateDisplay("Error")
		return
	}

	var res float64
	switch operator {
	case "+":
		res = p + c
	case "-":
		res = p - c
	case "*":
		res = p * c
	case "/":
		if c == 0 {
			updateDisplay("Div by 0")
			currentInput = ""
			previousInput = ""
			operator = ""
			return
		}
		res = p / c
	}

	ans := fmt.Sprintf("%g", res)
	updateDisplay(ans)
	previousInput = ans
	currentInput = ""
	operator = ""
}

// In the generic architecture, we might need a way to dynamically update UI.
// For now, since the transpiler parses AST, maybe it doesn't support SetText yet.
// Wait, I should check how the dynamic UI generation in `MainActivity.java` works.
func updateDisplay(text string) {
	// TODO: Need a way to interact with UI at runtime.
	fmt.Println("Display:", text)
}
