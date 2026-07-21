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
	// Run a multi-page app
	ui.RunApp(
		ui.Page{
			Name: "MainActivity",
			Root: ui.Column{
				ID: "main_layout",
				Style: ui.Style{
					BackgroundColor: "#1E1E2E",
					Padding:         24,
					Width:           ui.MatchParent,
					Height:          ui.MatchParent,
				},
				Children: []ui.Widget{
					ui.TextView{
						ID:   "welcome_text",
						Text: "Welcome to Go2APK Multi-Page Demo",
						Style: ui.Style{
							TextColor: "#CBA6F7",
							TextSize:  24,
							Margin:    16,
						},
					},
					ui.Button{
						ID: "btn_calc",
						Text: "Open Calculator",
						Style: ui.Style{
							BackgroundColor: "#89B4FA",
							TextColor: "#11111B",
							TextSize: 20,
							Margin: 16,
						},
						OnClick: func() {
							android.Navigate("CalculatorActivity")
						},
					},
					ui.Button{
						ID: "btn_http",
						Text: "Test HTTP GET",
						Style: ui.Style{
							BackgroundColor: "#A6E3A1",
							TextColor: "#11111B",
							TextSize: 20,
							Margin: 16,
						},
						OnClick: onTestHttp,
					},
					ui.TextView{
						ID: "http_result",
						Text: "Result will appear here",
						Style: ui.Style{
							TextColor: "#CDD6F4",
							TextSize: 14,
							Margin: 16,
						},
					},
				},
			},
		},
		ui.Page{
			Name: "CalculatorActivity",
			Root: ui.Column{
				ID: "calc_layout",
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
							ui.Button{ID: "btn_eq", Text: "=", Style: ui.Style{Width: 0, BackgroundColor: "#A6E3A1", TextColor: "#11111B", TextSize: 28, Margin: 8, Weight: 2.0}, OnClick: onEq},
						},
					},
				},
			},
		},
	)
}

func onTestHttp() {
	android.Permission("android.permission.INTERNET")
	res, err := android.HTTPGet("https://httpbin.org/get")
	if err != nil {
		android.UpdateText("http_result", "Error: "+err.Error())
	} else {
		// Truncate to first 100 characters for display
		if len(res) > 100 {
			res = res[:100] + "..."
		}
		android.UpdateText("http_result", res)
	}
}

// Handlers for CalculatorActivity

func appendInput(val string) {
	if currentInput == "0" && val != "." {
		currentInput = val
	} else {
		currentInput += val
	}
	android.UpdateText("display", currentInput)
	android.Animate("display", "alpha", 0.5, 50)
	android.Animate("display", "alpha", 1.0, 50)
}

func setOp(op string) {
	if currentInput != "" {
		if previousInput != "" {
			calculate()
		} else {
			previousInput = currentInput
		}
		currentInput = ""
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

func updateDisplay(text string) {
	fmt.Println("Display:", text)
	android.UpdateText("display", text)
}

func onClear() {
	currentInput = "0"
	previousInput = ""
	operator = ""
	updateDisplay("0")
}

func on0()   { appendInput("0") }
func on1()   { appendInput("1") }
func on2()   { appendInput("2") }
func on3()   { appendInput("3") }
func on4()   { appendInput("4") }
func on5()   { appendInput("5") }
func on6()   { appendInput("6") }
func on7()   { appendInput("7") }
func on8()   { appendInput("8") }
func on9()   { appendInput("9") }
func onDot() { appendInput(".") }

func onAdd() { setOp("+") }
func onSub() { setOp("-") }
func onMul() { setOp("*") }
func onDiv() { setOp("/") }

func onEq() { calculate() }

func onLP() {}
func onRP() {}
