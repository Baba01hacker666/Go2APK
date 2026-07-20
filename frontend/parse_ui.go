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
		if pkg.Name != "main" {
			continue
		}
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

func parseComposite(compLit *ast.CompositeLit, events map[string]string) map[string]interface{} {
	fields := make(map[string]interface{})
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key := kv.Key.(*ast.Ident).Name
		if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
			fields[key] = bl.Value[1 : len(bl.Value)-1]
		} else if key == "Style" {
			fields["Style"] = parseStyle(kv.Value)
		} else if key == "Children" {
			var children []ir.Widget
			if cl, ok := kv.Value.(*ast.CompositeLit); ok {
				for _, childElt := range cl.Elts {
					if childWidget := parseWidget(childElt, events); childWidget != nil {
						children = append(children, childWidget)
					}
				}
			}
			fields["Children"] = children
		} else if key == "OnClick" || key == "OnChanged" {
			if id, ok := kv.Value.(*ast.Ident); ok {
				fields[key] = id.Name
			}
		}
	}
	return fields
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

	fields := parseComposite(compLit, events)
	id, _ := fields["ID"].(string)
	style, _ := fields["Style"].(ir.Style)

	switch typeName {
	case "Column":
		col := ir.ColumnWidget{ID: id, Style: style}
		if children, ok := fields["Children"].([]ir.Widget); ok {
			col.Children = children
		}
		return col
	case "Row":
		row := ir.RowWidget{ID: id, Style: style}
		if children, ok := fields["Children"].([]ir.Widget); ok {
			row.Children = children
		}
		return row
	case "TextView":
		tv := ir.TextViewWidget{ID: id, Style: style}
		tv.Text, _ = fields["Text"].(string)
		return tv
	case "Button":
		btn := ir.ButtonWidget{ID: id, Style: style}
		btn.Text, _ = fields["Text"].(string)
		if onClick, ok := fields["OnClick"].(string); ok {
			btn.OnClickFunc = onClick
			if id != "" {
				events[id+"_onclick"] = onClick
			}
		}
		return btn
	case "TextField":
		tf := ir.TextFieldWidget{ID: id, Style: style}
		tf.Placeholder, _ = fields["Placeholder"].(string)
		if onChanged, ok := fields["OnChanged"].(string); ok {
			tf.OnChangedFunc = onChanged
			if id != "" {
				events[id+"_onchanged"] = onChanged
			}
		}
		return tf
	}
	return nil
}
