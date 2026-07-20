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
				// Button Rows
				ui.Row{
					Style: ui.Style{Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_clear", Text: "C", Style: ui.Style{BackgroundColor: "#F38BA8", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("C") }},
						ui.Button{ID: "btn_lp", Text: "(", Style: ui.Style{BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("(") }},
						ui.Button{ID: "btn_rp", Text: ")", Style: ui.Style{BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput(")") }},
						ui.Button{ID: "btn_div", Text: "/", Style: ui.Style{BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("/") }},
					},
				},
				ui.Row{
					Style: ui.Style{Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_7", Text: "7", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("7") }},
						ui.Button{ID: "btn_8", Text: "8", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("8") }},
						ui.Button{ID: "btn_9", Text: "9", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("9") }},
						ui.Button{ID: "btn_mul", Text: "*", Style: ui.Style{BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("*") }},
					},
				},
				ui.Row{
					Style: ui.Style{Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_4", Text: "4", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("4") }},
						ui.Button{ID: "btn_5", Text: "5", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("5") }},
						ui.Button{ID: "btn_6", Text: "6", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("6") }},
						ui.Button{ID: "btn_sub", Text: "-", Style: ui.Style{BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("-") }},
					},
				},
				ui.Row{
					Style: ui.Style{Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_1", Text: "1", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("1") }},
						ui.Button{ID: "btn_2", Text: "2", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("2") }},
						ui.Button{ID: "btn_3", Text: "3", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("3") }},
						ui.Button{ID: "btn_add", Text: "+", Style: ui.Style{BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("+") }},
					},
				},
				ui.Row{
					Style: ui.Style{Width: ui.MatchParent, Margin: 4},
					Children: []ui.Widget{
						ui.Button{ID: "btn_0", Text: "0", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("0") }},
						ui.Button{ID: "btn_dot", Text: ".", Style: ui.Style{BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput(".") }},
						ui.Button{ID: "btn_del", Text: "DEL", Style: ui.Style{BackgroundColor: "#F38BA8", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("DEL") }},
						ui.Button{ID: "btn_eq", Text: "=", Style: ui.Style{BackgroundColor: "#A6E3A1", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 1.0}, OnClick: func() { handleInput("=") }},
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
}
