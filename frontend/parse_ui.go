package frontend

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/Baba01hacker666/Go2APK/ir"
	"golang.org/x/tools/go/packages"
)

// ExtractUI parses the AST for ui.Run and extracts the widget tree.
func ExtractUI(pkgs []*packages.Package) (ir.Widget, map[string]string) {
	events := make(map[string]string)
	var rootWidget ir.Widget

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				// Look for ui.Run
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				id, ok := sel.X.(*ast.Ident)
				if !ok || id.Name != "ui" || sel.Sel.Name != "Run" {
					return true
				}

				fmt.Println("Found ui.Run call!")

				// Found ui.Run, parse the first argument
				if len(call.Args) > 0 {
					rootWidget = parseWidget(call.Args[0], events)
				}
				return false
			})
		}
	}
	return rootWidget, events
}

func parseStyle(expr ast.Expr) ir.Style {
	style := ir.Style{}
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return style
	}
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key := kv.Key.(*ast.Ident).Name
		if bl, ok := kv.Value.(*ast.BasicLit); ok {
			switch key {
			case "BackgroundColor":
				style.BackgroundColor = strings.Trim(bl.Value, "\"")
			case "TextColor":
				style.TextColor = strings.Trim(bl.Value, "\"")
			case "TextSize":
				style.TextSize, _ = strconv.Atoi(bl.Value)
			case "Padding":
				style.Padding, _ = strconv.Atoi(bl.Value)
			case "Margin":
				style.Margin, _ = strconv.Atoi(bl.Value)
			case "Width":
				style.Width, _ = strconv.Atoi(bl.Value)
			case "Height":
				style.Height, _ = strconv.Atoi(bl.Value)
			case "Weight":
				v, _ := strconv.ParseFloat(bl.Value, 32)
				style.Weight = float32(v)
			}
		} else if sel, ok := kv.Value.(*ast.SelectorExpr); ok {
			// e.g., ui.MatchParent or ui.WrapContent
			if sel.Sel.Name == "MatchParent" {
				if key == "Width" {
					style.Width = -1
				}
				if key == "Height" {
					style.Height = -1
				}
			} else if sel.Sel.Name == "WrapContent" {
				if key == "Width" {
					style.Width = -2
				}
				if key == "Height" {
					style.Height = -2
				}
			}
		} else if un, ok := kv.Value.(*ast.UnaryExpr); ok {
			if un.Op == token.SUB {
				if bl, ok := un.X.(*ast.BasicLit); ok {
					val, _ := strconv.Atoi(bl.Value)
					switch key {
					case "Width":
						style.Width = -val
					case "Height":
						style.Height = -val
					}
				}
			}
		}
	}
	return style
}

func parseWidget(expr ast.Expr, events map[string]string) ir.Widget {
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	var typeName string
	switch t := compLit.Type.(type) {
	case *ast.SelectorExpr:
		typeName = t.Sel.Name
	case *ast.Ident:
		typeName = t.Name
	}

	switch typeName {
	case "Column":
		fmt.Println("Found Column")
		col := ir.ColumnWidget{}
		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key := kv.Key.(*ast.Ident).Name

			if key == "ID" {
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					col.ID = bl.Value[1 : len(bl.Value)-1]
				}
			} else if key == "Style" {
				col.Style = parseStyle(kv.Value)
			} else if key == "Children" {
				if cl, ok := kv.Value.(*ast.CompositeLit); ok {
					for _, childElt := range cl.Elts {
						childWidget := parseWidget(childElt, events)
						if childWidget != nil {
							col.Children = append(col.Children, childWidget)
						}
					}
				}
			}
		}
		return col
	case "Row":
		row := ir.RowWidget{}
		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key := kv.Key.(*ast.Ident).Name

			if key == "ID" {
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					row.ID = bl.Value[1 : len(bl.Value)-1]
				}
			} else if key == "Style" {
				row.Style = parseStyle(kv.Value)
			} else if key == "Children" {
				if cl, ok := kv.Value.(*ast.CompositeLit); ok {
					for _, childElt := range cl.Elts {
						childWidget := parseWidget(childElt, events)
						if childWidget != nil {
							row.Children = append(row.Children, childWidget)
						}
					}
				}
			}
		}
		return row
	case "TextView":
		tv := ir.TextViewWidget{}
		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key := kv.Key.(*ast.Ident).Name
			if key == "ID" {
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					tv.ID = bl.Value[1 : len(bl.Value)-1]
				}
			} else if key == "Text" {
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					tv.Text = bl.Value[1 : len(bl.Value)-1]
				}
			} else if key == "Style" {
				tv.Style = parseStyle(kv.Value)
			}
		}
		return tv
	case "Button":
		btn := ir.ButtonWidget{}
		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key := kv.Key.(*ast.Ident).Name
			if key == "ID" {
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					btn.ID = bl.Value[1 : len(bl.Value)-1]
				}
			} else if key == "Text" {
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					btn.Text = bl.Value[1 : len(bl.Value)-1]
				}
			} else if key == "Style" {
				btn.Style = parseStyle(kv.Value)
			} else if key == "OnClick" {
				if id, ok := kv.Value.(*ast.Ident); ok {
					btn.OnClickFunc = id.Name
					if btn.ID != "" {
						events[btn.ID+"_onclick"] = id.Name
					}
				}
			}
		}
		return btn
	case "TextField":
		tf := ir.TextFieldWidget{}
		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key := kv.Key.(*ast.Ident).Name
			if key == "ID" {
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					tf.ID = bl.Value[1 : len(bl.Value)-1]
				}
			} else if key == "Placeholder" {
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					tf.Placeholder = bl.Value[1 : len(bl.Value)-1]
				}
			} else if key == "Style" {
				tf.Style = parseStyle(kv.Value)
			} else if key == "OnChanged" {
				if id, ok := kv.Value.(*ast.Ident); ok {
					tf.OnChangedFunc = id.Name
					if tf.ID != "" {
						events[tf.ID+"_onchanged"] = id.Name
					}
				}
			}
		}
		return tf
	}
	return nil
}
