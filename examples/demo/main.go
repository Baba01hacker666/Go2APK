package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"strconv"

	"github.com/Baba01hacker666/Go2APK/android"
	"github.com/Baba01hacker666/Go2APK/ui"
)

var (
	currentExpression string
)

func main() {
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
						ID:   "btn_calc",
						Text: "Open Calculator",
						Style: ui.Style{
							BackgroundColor: "#89B4FA",
							TextColor:       "#11111B",
							TextSize:        20,
							Margin:          16,
						},
						OnClick: onOpenCalc,
					},
				},
			},
		},
		ui.Page{
			Name: "BasicCalculatorActivity",
			Root: ui.Column{
				ID: "calc_layout",
				Style: ui.Style{
					BackgroundColor: "#1E1E2E",
					Padding:         24,
					Width:           ui.MatchParent,
					Height:          ui.MatchParent,
				},
				Children: []ui.Widget{
					ui.Button{
						ID:   "btn_back",
						Text: "Back to Home",
						Style: ui.Style{
							BackgroundColor: "#F38BA8",
							TextColor:       "#11111B",
							TextSize:        16,
							Margin:          8,
							Width:           ui.MatchParent,
						},
						OnClick: onBack,
					},
					ui.TextView{
						ID:   "display",
						Text: "0",
						Style: ui.Style{
							TextColor:       "#CBA6F7",
							TextSize:        48,
							BackgroundColor: "#181825",
							Padding:         32,
							Margin:          16,
							Width:           ui.MatchParent,
						},
					},
					ui.Button{
						ID:      "btn_adv",
						Text:    "Advanced Mode",
						Style:   ui.Style{Width: ui.MatchParent, BackgroundColor: "#F5C2E7", TextColor: "#11111B", TextSize: 20, Margin: 8},
						OnClick: onAdvMode,
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
		},
		ui.Page{
			Name: "AdvancedCalculatorActivity",
			Root: ui.Column{
				ID: "adv_layout",
				Style: ui.Style{
					BackgroundColor: "#1E1E2E",
					Padding:         24,
					Width:           ui.MatchParent,
					Height:          ui.MatchParent,
				},
				Children: []ui.Widget{
					ui.Button{
						ID:   "btn_back_adv",
						Text: "Back to Home",
						Style: ui.Style{
							BackgroundColor: "#F38BA8",
							TextColor:       "#11111B",
							TextSize:        16,
							Margin:          8,
							Width:           ui.MatchParent,
						},
						OnClick: onBack,
					},
					ui.TextView{
						ID:   "display",
						Text: "0",
						Style: ui.Style{
							TextColor:       "#CBA6F7",
							TextSize:        48,
							BackgroundColor: "#181825",
							Padding:         32,
							Margin:          16,
							Width:           ui.MatchParent,
						},
					},
					ui.Button{
						ID:      "btn_bas",
						Text:    "Basic Mode",
						Style:   ui.Style{Width: ui.MatchParent, BackgroundColor: "#F5C2E7", TextColor: "#11111B", TextSize: 20, Margin: 8},
						OnClick: onBasMode,
					},
					ui.Row{
						Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
						Children: []ui.Widget{
							ui.Button{ID: "btn_sin", Text: "sin(", Style: ui.Style{Width: 0, BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onSin},
							ui.Button{ID: "btn_cos", Text: "cos(", Style: ui.Style{Width: 0, BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onCos},
							ui.Button{ID: "btn_tan", Text: "tan(", Style: ui.Style{Width: 0, BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onTan},
							ui.Button{ID: "btn_clear", Text: "C", Style: ui.Style{Width: 0, BackgroundColor: "#F38BA8", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onClear},
						},
					},
					ui.Row{
						Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
						Children: []ui.Widget{
							ui.Button{ID: "btn_lp", Text: "(", Style: ui.Style{Width: 0, BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onLP},
							ui.Button{ID: "btn_rp", Text: ")", Style: ui.Style{Width: 0, BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onRP},
							ui.Button{ID: "btn_sqrt", Text: "sqrt(", Style: ui.Style{Width: 0, BackgroundColor: "#89B4FA", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onSqrt},
							ui.Button{ID: "btn_div", Text: "/", Style: ui.Style{Width: 0, BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onDiv},
						},
					},
					ui.Row{
						Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
						Children: []ui.Widget{
							ui.Button{ID: "btn_7", Text: "7", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on7},
							ui.Button{ID: "btn_8", Text: "8", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on8},
							ui.Button{ID: "btn_9", Text: "9", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on9},
							ui.Button{ID: "btn_mul", Text: "*", Style: ui.Style{Width: 0, BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onMul},
						},
					},
					ui.Row{
						Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
						Children: []ui.Widget{
							ui.Button{ID: "btn_4", Text: "4", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on4},
							ui.Button{ID: "btn_5", Text: "5", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on5},
							ui.Button{ID: "btn_6", Text: "6", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on6},
							ui.Button{ID: "btn_sub", Text: "-", Style: ui.Style{Width: 0, BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onSub},
						},
					},
					ui.Row{
						Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
						Children: []ui.Widget{
							ui.Button{ID: "btn_1", Text: "1", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on1},
							ui.Button{ID: "btn_2", Text: "2", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on2},
							ui.Button{ID: "btn_3", Text: "3", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on3},
							ui.Button{ID: "btn_add", Text: "+", Style: ui.Style{Width: 0, BackgroundColor: "#F9E2AF", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onAdd},
						},
					},
					ui.Row{
						Style: ui.Style{Height: 0, Weight: 1, Width: ui.MatchParent, Margin: 4},
						Children: []ui.Widget{
							ui.Button{ID: "btn_0", Text: "0", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: on0},
							ui.Button{ID: "btn_dot", Text: ".", Style: ui.Style{Width: 0, BackgroundColor: "#313244", TextColor: "#CDD6F4", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onDot},
							ui.Button{ID: "btn_del", Text: "DEL", Style: ui.Style{Width: 0, BackgroundColor: "#F38BA8", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onDel},
							ui.Button{ID: "btn_eq", Text: "=", Style: ui.Style{Width: 0, BackgroundColor: "#A6E3A1", TextColor: "#11111B", TextSize: 24, Margin: 4, Weight: 1.0}, OnClick: onEq},
						},
					},
				},
			},
		},
	)
}

func onOpenCalc() {
	android.Navigate("BasicCalculatorActivity")
}

func onBack() {
	android.Navigate("MainActivity")
}

func onAdvMode() {
	android.Navigate("AdvancedCalculatorActivity")
}

func onBasMode() {
	android.Navigate("BasicCalculatorActivity")
}

func updateDisplay() {
	disp := currentExpression
	if disp == "" {
		disp = "0"
	}
	android.UpdateText("display", disp)
	android.Animate("display", "alpha", 0.5, 50)
	android.Animate("display", "alpha", 1.0, 50)
}

func appendInput(val string) {
	currentExpression += val
	updateDisplay()
}

func onClear() {
	currentExpression = ""
	updateDisplay()
}

func onDel() {
	if len(currentExpression) > 0 {
		currentExpression = currentExpression[:len(currentExpression)-1]
	}
	updateDisplay()
}

func onEq() {
	if currentExpression == "" {
		return
	}
	expr, err := parser.ParseExpr(currentExpression)
	if err != nil {
		android.UpdateText("display", "Error")
		currentExpression = ""
		return
	}

	res, err := evaluate(expr)
	if err != nil {
		android.UpdateText("display", "Error: "+err.Error())
		currentExpression = ""
		return
	}

	currentExpression = fmt.Sprintf("%g", res)
	updateDisplay()
}

func evaluate(expr ast.Expr) (float64, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.INT || e.Kind == token.FLOAT {
			return strconv.ParseFloat(e.Value, 64)
		}
		return 0, fmt.Errorf("unsupported literal")
	case *ast.ParenExpr:
		return evaluate(e.X)
	case *ast.BinaryExpr:
		left, err := evaluate(e.X)
		if err != nil {
			return 0, err
		}
		right, err := evaluate(e.Y)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			if right == 0 {
				return 0, fmt.Errorf("div by 0")
			}
			return left / right, nil
		}
		return 0, fmt.Errorf("unsupported op")
	case *ast.CallExpr:
		ident, ok := e.Fun.(*ast.Ident)
		if !ok {
			return 0, fmt.Errorf("unsupported call")
		}
		if len(e.Args) != 1 {
			return 0, fmt.Errorf("expected 1 arg")
		}
		arg, err := evaluate(e.Args[0])
		if err != nil {
			return 0, err
		}
		switch ident.Name {
		case "sin":
			return math.Sin(arg), nil
		case "cos":
			return math.Cos(arg), nil
		case "tan":
			return math.Tan(arg), nil
		case "sqrt":
			return math.Sqrt(arg), nil
		}
		return 0, fmt.Errorf("unknown func")
	}
	return 0, fmt.Errorf("unsupported expr")
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

func onAdd() { appendInput("+") }
func onSub() { appendInput("-") }
func onMul() { appendInput("*") }
func onDiv() { appendInput("/") }

func onLP() { appendInput("(") }
func onRP() { appendInput(")") }

func onSin()  { appendInput("sin(") }
func onCos()  { appendInput("cos(") }
func onTan()  { appendInput("tan(") }
func onSqrt() { appendInput("sqrt(") }
