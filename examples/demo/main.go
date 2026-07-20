package main

import (
	"fmt"
	"strconv"

	"github.com/Baba01hacker666/Go2APK/android"
	"github.com/Baba01hacker666/Go2APK/ui"
)

var (
	currentInput  string
	previousInput string
	operator      string
)

func main() {
	// A beautiful, modern dark-themed calculator UI
	ui.Run(
		ui.Column{
			ID: "main_layout",
			Style: ui.Style{
				BackgroundColor: "#1E1E2E", // Dark violet-blue background
				Padding:         24,
				Width:           ui.MatchParent,
				Height:          ui.MatchParent,
			},
			Children: []ui.Widget{
				ui.TextView{
					ID:   "display",
					Text: "0",
					Style: ui.Style{
						TextColor:       "#CBA6F7", // Pastel purple text
						TextSize:        64,
						BackgroundColor: "#181825", // Slightly darker container
						Padding:         32,
						Margin:          16,
						Width:           ui.MatchParent,
					},
				},
				ui.Row{
					Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_clear", Text: "C", Style: ui.Style{Width: 0, BackgroundColor: "#F38BA8", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onClear},
						ui.Button{ID: "btn_lp", Text: "(", Style: ui.Style{Width: 0, BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onLP},
						ui.Button{ID: "btn_rp", Text: ")", Style: ui.Style{Width: 0, BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onRP},
						ui.Button{ID: "btn_div", Text: "/", Style: ui.Style{Width: 0, BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onDiv},
					},
				},
				ui.Row{
					Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_7", Text: "7", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on7},
						ui.Button{ID: "btn_8", Text: "8", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on8},
						ui.Button{ID: "btn_9", Text: "9", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on9},
						ui.Button{ID: "btn_mul", Text: "*", Style: ui.Style{Width: 0, BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onMul},
					},
				},
				ui.Row{
					Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_4", Text: "4", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on4},
						ui.Button{ID: "btn_5", Text: "5", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on5},
						ui.Button{ID: "btn_6", Text: "6", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on6},
						ui.Button{ID: "btn_sub", Text: "-", Style: ui.Style{Width: 0, BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onSub},
					},
				},
				ui.Row{
					Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_1", Text: "1", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on1},
						ui.Button{ID: "btn_2", Text: "2", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on2},
						ui.Button{ID: "btn_3", Text: "3", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on3},
						ui.Button{ID: "btn_add", Text: "+", Style: ui.Style{Width: 0, BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onAdd},
					},
				},
				ui.Row{
					Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_0", Text: "0", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: on0},
						ui.Button{ID: "btn_dot", Text: ".", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onDot},
						ui.Button{ID: "btn_del", Text: "DEL", Style: ui.Style{Width: 0, BackgroundColor: "#F38BA8", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onDel},
						ui.Button{ID: "btn_eq", Text: "=", Style: ui.Style{Width: 0, BackgroundColor: "#A6E3A1", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: onEq},
					},
				},
			},
		},
	)
}

func handleInput(in string) {
	switch in {
	case "C":
		currentInput = ""
		previousInput = ""
		operator = ""
		updateDisplay("0")
	case "DEL":
		if len(currentInput) > 0 {
			currentInput = currentInput[:len(currentInput)-1]
		}
		if currentInput == "" {
			updateDisplay("0")
		} else {
			updateDisplay(currentInput)
		}
	case "+", "-", "*", "/":
		if currentInput != "" {
			if previousInput != "" {
				calculate()
			} else {
				previousInput = currentInput
				currentInput = ""
			}
		}
		operator = in
	case "=":
		calculate()
	default: // numbers and .
		currentInput += in
		updateDisplay(currentInput)
	}
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

func updateDisplay(text string) {
	fmt.Println("Display:", text)
	android.UpdateText("display", text)
}

var (
	onClear = func() { handleInput("C") }
	onLP    = func() { handleInput("(") }
	onRP    = func() { handleInput(")") }
	onDiv   = func() { handleInput("/") }
	on7     = func() { handleInput("7") }
	on8     = func() { handleInput("8") }
	on9     = func() { handleInput("9") }
	onMul   = func() { handleInput("*") }
	on4     = func() { handleInput("4") }
	on5     = func() { handleInput("5") }
	on6     = func() { handleInput("6") }
	onSub   = func() { handleInput("-") }
	on1     = func() { handleInput("1") }
	on2     = func() { handleInput("2") }
	on3     = func() { handleInput("3") }
	onAdd   = func() { handleInput("+") }
	on0     = func() { handleInput("0") }
	onDot   = func() { handleInput(".") }
	onDel   = func() { handleInput("DEL") }
	onEq    = func() { handleInput("=") }
)
