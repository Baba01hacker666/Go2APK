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

// ExtractUI parses the AST for ui.Run and extracts the widget tree, permissions,
// and broadcast receiver declarations.
func ExtractUI(pkgs []*packages.Package) (ir.Widget, map[string]string, []string, []ir.BroadcastReceiverDecl) {
	events := make(map[string]string)
	var rootWidget ir.Widget
	var permissions []string
	var receivers []ir.BroadcastReceiverDecl

	for _, pkg := range pkgs {
		if pkg.Name != "main" {
			continue
		}
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				id, ok := sel.X.(*ast.Ident)
				if !ok || (id.Name != "ui" && id.Name != "android") {
					return true
				}

				if id.Name == "ui" && sel.Sel.Name == "Run" {
					fmt.Println("Found ui.Run call!")
					if len(call.Args) > 0 {
						rootWidget = parseWidget(call.Args[0], events)
					}
					return false
				}

				if id.Name == "android" {
					switch sel.Sel.Name {
					case "Permission":
						// android.Permission("android.permission.XYZ")
						if len(call.Args) == 1 {
							if perm := stringLit(call.Args[0]); perm != "" {
								permissions = append(permissions, perm)
							}
						}

					case "BroadcastReceiver":
						// android.BroadcastReceiver(action, eventName, handler)
						if len(call.Args) >= 2 {
							action := stringLit(call.Args[0])
							eventName := stringLit(call.Args[1])
							if action != "" && eventName != "" {
								receivers = append(receivers, ir.BroadcastReceiverDecl{
									Name:     eventName,
									Action:   action,
									Exported: false,
								})
							}
						}
					}
				}

				return true
			})
		}
	}
	return rootWidget, events, permissions, receivers
}

// stringLit extracts the string value from a basic string literal AST node.
func stringLit(expr ast.Expr) string {
	bl, ok := expr.(*ast.BasicLit)
	if !ok || bl.Kind != token.STRING {
		return ""
	}
	return strings.Trim(bl.Value, "\"")
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
		} else if id, ok := kv.Value.(*ast.Ident); ok && (id.Name == "true" || id.Name == "false") {
			fields[key] = id.Name == "true"
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
	css, _ := fields["CSS"].(string)

	switch typeName {
	case "Column":
		col := ir.ColumnWidget{ID: id, Style: style, CSS: css}
		if children, ok := fields["Children"].([]ir.Widget); ok {
			col.Children = children
		}
		return col
	case "Row":
		row := ir.RowWidget{ID: id, Style: style, CSS: css}
		if children, ok := fields["Children"].([]ir.Widget); ok {
			row.Children = children
		}
		return row
	case "TextView":
		tv := ir.TextViewWidget{ID: id, Style: style, CSS: css}
		tv.Text, _ = fields["Text"].(string)
		return tv
	case "Button":
		btn := ir.ButtonWidget{ID: id, Style: style, CSS: css}
		btn.Text, _ = fields["Text"].(string)
		if onClick, ok := fields["OnClick"].(string); ok {
			btn.OnClickFunc = onClick
			if id != "" {
				events[id+"_onclick"] = onClick
			}
		}
		return btn
	case "TextField":
		tf := ir.TextFieldWidget{ID: id, Style: style, CSS: css}
		tf.Placeholder, _ = fields["Placeholder"].(string)
		if onChanged, ok := fields["OnChanged"].(string); ok {
			tf.OnChangedFunc = onChanged
			if id != "" {
				events[id+"_onchanged"] = onChanged
			}
		}
		return tf
	case "Image":
		img := ir.ImageWidget{ID: id, Style: style, CSS: css}
		img.Src, _ = fields["Src"].(string)
		return img
	case "Audio":
		aud := ir.AudioWidget{ID: id}
		aud.Src, _ = fields["Src"].(string)
		if autoPlay, ok := fields["AutoPlay"].(bool); ok {
			aud.AutoPlay = autoPlay
		}
		return aud
	case "Video":
		vid := ir.VideoWidget{ID: id, Style: style, CSS: css}
		vid.Src, _ = fields["Src"].(string)
		return vid
	case "WebView":
		wv := ir.WebViewWidget{ID: id, Style: style, CSS: css}
		wv.Src, _ = fields["Src"].(string)
		wv.HTML, _ = fields["HTML"].(string)
		return wv
	case "ScrollView":
		sv := ir.ScrollViewWidget{ID: id, Style: style, CSS: css}
		if children, ok := fields["Children"].([]ir.Widget); ok {
			sv.Children = children
		}
		return sv
	case "CardView":
		cv := ir.CardViewWidget{ID: id, Style: style, CSS: css}
		if children, ok := fields["Children"].([]ir.Widget); ok {
			cv.Children = children
		}
		return cv
	case "ProgressBar":
		pb := ir.ProgressBarWidget{ID: id, Style: style, CSS: css}
		return pb
	case "Switch":
		sw := ir.SwitchWidget{ID: id, Style: style, CSS: css}
		if checked, ok := fields["Checked"].(bool); ok {
			sw.Checked = checked
		}
		if onChanged, ok := fields["OnChanged"].(string); ok {
			if id != "" {
				events[id+"_onchanged"] = onChanged
			}
		}
		return sw
	}
	return nil
}
