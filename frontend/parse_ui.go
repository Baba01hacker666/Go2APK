package frontend

import (
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"os"
	pathpkg "path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Baba01hacker666/Go2APK/ir"
	"golang.org/x/tools/go/packages"
)

// DefaultAndroidImport is assumed when the app source does not import the
// Go2APK android runtime under a different (vendored or forked) path.
const DefaultAndroidImport = "github.com/Baba01hacker666/Go2APK/android"

// UIProgram is everything the frontend can learn about an app by reading its
// source: the widget tree, the events that tie widgets to Go handlers, and the
// Android platform features the app declares.
type UIProgram struct {
	// Root is the widget tree passed to ui.Run, for single-page apps.
	Root ir.Widget
	// Pages are the screens passed to ui.RunApp, for multi-page apps.
	Pages []ir.Page
	// Events maps generated Android event names to their Go handlers.
	Events map[string]ir.EventHandler
	// Permissions are the manifest permissions the app declares.
	Permissions []string
	// Receivers are the broadcast receivers the app declares.
	Receivers []ir.BroadcastReceiverDecl
	// HasVPN reports whether the app starts a VPN service.
	HasVPN bool
	// AndroidImport is the import path the app uses for the Go2APK runtime.
	AndroidImport string
}

// widgetParser holds the state of a single extraction pass.
type widgetParser struct {
	events map[string]ir.EventHandler
	errs   []error

	// imports holds the imports of the file being inspected. An inline handler
	// is copied into a generated file, so the imports it relies on have to be
	// carried across with it.
	imports []fileImport

	// fset owns the positions of the files being parsed. It is needed to print
	// an inline func literal back into source form.
	fset *token.FileSet
	// info holds the type checker's view of the package, used to read the
	// signature of a handler given by name.
	info *types.Info

	root        ir.Widget
	pages       []ir.Page
	permissions []string
	receivers   []ir.BroadcastReceiverDecl
	hasVPN      bool

	// dir is the directory of the file currently being parsed. Relative paths
	// passed to ui.LoadFile resolve against it rather than against the
	// directory go2apk happened to be invoked from.
	dir string

	// autoID counts widgets that need a generated ID so that every generated
	// name stays unique and stable for the duration of the build.
	autoID int
}

// ExtractUI parses the AST of the given packages looking for ui.Run, ui.RunApp
// and the android package declarations, and turns them into a UIProgram.
//
// Only package main is inspected: that is where an application's entry point
// and therefore its UI declaration must live. Problems found in the source are
// reported together rather than one at a time, so a single run surfaces every
// mistake.
func ExtractUI(pkgs []*packages.Package) (UIProgram, error) {
	p := &widgetParser{events: make(map[string]ir.EventHandler)}

	for _, pkg := range pkgs {
		if pkg.Name != "main" {
			continue
		}
		for _, file := range pkg.Syntax {
			p.dir = ""
			p.fset = pkg.Fset
			p.info = pkg.TypesInfo
			if pkg.Fset != nil {
				if pos := pkg.Fset.Position(file.Pos()); pos.Filename != "" {
					p.dir = filepath.Dir(pos.Filename)
				}
			}
			p.inspect(file)
		}
	}

	androidImport := detectAndroidImport(pkgs)
	if err := p.err(); err != nil {
		return UIProgram{}, err
	}

	return UIProgram{
		Root:          p.root,
		Pages:         p.pages,
		Events:        p.events,
		Permissions:   p.permissions,
		Receivers:     p.receivers,
		HasVPN:        p.hasVPN,
		AndroidImport: androidImport,
	}, nil
}

func (p *widgetParser) inspect(file *ast.File) {
	p.imports = fileImports(file)
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
		if !ok {
			return true
		}
		switch id.Name {
		case "ui":
			switch sel.Sel.Name {
			case "Run":
				if len(call.Args) > 0 {
					p.root = p.widget(call.Args[0])
				}
				return false
			case "RunApp":
				for _, arg := range call.Args {
					if page := p.page(arg); page.Name != "" {
						p.pages = append(p.pages, page)
					}
				}
				return false
			}
		case "android":
			p.androidCall(sel.Sel.Name, call)
		}
		return true
	})
}

func (p *widgetParser) androidCall(name string, call *ast.CallExpr) {
	switch name {
	case "Permission":
		if len(call.Args) == 1 {
			if perm, ok := p.stringValue(call.Args[0], "android.Permission"); ok && perm != "" {
				p.permissions = append(p.permissions, perm)
			}
		}
	case "BroadcastReceiver":
		if len(call.Args) < 2 {
			p.errorf("android.BroadcastReceiver requires (action, eventName, handler)")
			return
		}
		action, okAction := p.stringValue(call.Args[0], "android.BroadcastReceiver action")
		eventName, okEvent := p.stringValue(call.Args[1], "android.BroadcastReceiver eventName")
		if !okAction || !okEvent {
			return
		}
		if action == "" || eventName == "" {
			p.errorf("android.BroadcastReceiver requires a non-empty action and event name")
			return
		}
		p.receivers = append(p.receivers, ir.BroadcastReceiverDecl{
			Name:     eventName,
			Action:   action,
			Exported: false,
		})
	case "StartVPN":
		p.hasVPN = true
	}
}

func (p *widgetParser) page(expr ast.Expr) ir.Page {
	var page ir.Page
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return page
	}
	f := p.fields(compLit)
	name, _ := f.str("Name")
	root := f.widget("Root")
	if name != "" && !javaClassNameRE.MatchString(name) {
		p.errorf("page name %q is not a valid Android activity name; use letters, digits and underscores, starting with a letter", name)
		return ir.Page{}
	}
	page.Name = name
	page.Root = root
	return page
}

// widget turns one composite literal into an IR widget. The literal is
// recognised by its type name only, which mirrors how the ui package is meant
// to be used and keeps the parser independent of type checking.
func (p *widgetParser) widget(expr ast.Expr) ir.Widget {
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	typeName := ""
	switch t := compLit.Type.(type) {
	case *ast.SelectorExpr:
		typeName = t.Sel.Name
	case *ast.Ident:
		typeName = t.Name
	}
	if typeName == "" {
		return nil
	}
	f := p.fields(compLit)
	return p.buildWidget(typeName, f)
}

func (p *widgetParser) buildWidget(typeName string, f widgetFields) ir.Widget {
	// Resolve the ID first: handlers are registered under a name derived from
	// it, and it may have to be generated before they are recorded.
	id := f.id(typeName)

	switch typeName {
	case "Column":
		children, _ := f.children("Children")
		return ir.ColumnWidget{ID: id, Style: f.style(), CSS: f.css(), Children: children}
	case "Row":
		children, _ := f.children("Children")
		return ir.RowWidget{ID: id, Style: f.style(), CSS: f.css(), Children: children}
	case "TextView":
		text, _ := f.str("Text")
		return ir.TextViewWidget{ID: id, Text: text, Style: f.style(), CSS: f.css()}
	case "Button":
		text, _ := f.str("Text")
		return ir.ButtonWidget{ID: id, Text: text, Style: f.style(), CSS: f.css(),
			OnClickFunc: f.handler("OnClick", typeName, "onclick")}
	case "TextField":
		placeholder, _ := f.str("Placeholder")
		return ir.TextFieldWidget{ID: id, Placeholder: placeholder, Style: f.style(), CSS: f.css(),
			OnChangedFunc: f.handler("OnChanged", typeName, "onchanged")}
	case "Image":
		src, _ := f.str("Src")
		return ir.ImageWidget{ID: id, Src: src, Style: f.style(), CSS: f.css()}
	case "Audio":
		src, _ := f.str("Src")
		autoPlay, _ := f.bool("AutoPlay")
		return ir.AudioWidget{ID: id, Src: src, AutoPlay: autoPlay}
	case "Video":
		src, _ := f.str("Src")
		return ir.VideoWidget{ID: id, Src: src, Style: f.style(), CSS: f.css()}
	case "XMLView":
		layout, _ := f.str("Layout")
		return ir.XMLView{ID: id, Layout: layout, Style: f.style(), CSS: f.css()}
	case "WebView":
		src, _ := f.str("Src")
		html, _ := f.str("HTML")
		return ir.WebViewWidget{ID: id, Src: src, HTML: html, Style: f.style(), CSS: f.css()}
	case "ScrollView":
		children, _ := f.children("Children")
		return ir.ScrollViewWidget{ID: id, Style: f.style(), CSS: f.css(), Children: children}
	case "CardView":
		children, _ := f.children("Children")
		return ir.CardViewWidget{ID: id, Style: f.style(), CSS: f.css(), Children: children}
	case "ProgressBar":
		return ir.ProgressBarWidget{ID: id, Style: f.style(), CSS: f.css()}
	case "Switch":
		checked, _ := f.bool("Checked")
		return ir.SwitchWidget{ID: id, Checked: checked, Style: f.style(), CSS: f.css(),
			OnChangedFunc: f.handler("OnChanged", typeName, "onchanged")}
	}
	p.errorf("unknown widget ui.%s; if this is a typo it would otherwise be dropped silently", typeName)
	return nil
}

// widgetFields gives typed, error-reporting access to the fields of a widget
// composite literal.
type widgetFields struct {
	p    *widgetParser
	kv   map[string]ast.Expr
	kind string
	// resolvedID caches the ID so that a generated ID is handed out once per
	// widget even though several fields consult it.
	resolvedID string
	idChecked  bool
}

func (p *widgetParser) fields(compLit *ast.CompositeLit) widgetFields {
	f := widgetFields{p: p, kv: make(map[string]ast.Expr, len(compLit.Elts))}
	switch t := compLit.Type.(type) {
	case *ast.SelectorExpr:
		f.kind = t.Sel.Name
	case *ast.Ident:
		f.kind = t.Name
	}
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		f.kv[key.Name] = kv.Value
	}
	return f
}

// id resolves the widget's ID, generating one when the widget needs an identity
// but the author did not supply one.
//
// An ID is required whenever the widget carries an event handler: the handler
// is dispatched by name over JNI, and that name is derived from the ID. The
// generator also emits the ID as a Java field name, so it must be a valid Java
// identifier.
func (f *widgetFields) id(kind string) string {
	if f.idChecked {
		return f.resolvedID
	}
	f.idChecked = true
	raw, ok := f.p.stringValue(f.kv["ID"], "ID")
	if ok && raw != "" {
		if !javaIdentRE.MatchString(raw) {
			f.p.errorf("widget ID %q is not a valid Java identifier; IDs are emitted as Java field names, so use letters, digits and underscores starting with a letter", raw)
			return ""
		}
		if javaReserved[raw] {
			f.p.errorf("widget ID %q is a Java reserved word and cannot be used as a field name", raw)
			return ""
		}
		f.resolvedID = raw
		return f.resolvedID
	}
	if raw != "" && !ok {
		return ""
	}
	if f.needsID() {
		f.p.autoID++
		f.resolvedID = fmt.Sprintf("go2apk_%s_%d", strings.ToLower(kind), f.p.autoID)
	}
	return f.resolvedID
}

// needsID reports whether the widget declares an event handler and therefore
// must have an identity even if the author did not give it one.
func (f *widgetFields) needsID() bool {
	for _, key := range []string{"OnClick", "OnChanged"} {
		if _, ok := f.kv[key]; ok {
			return true
		}
	}
	return false
}

func (f *widgetFields) str(key string) (string, error) {
	expr, ok := f.kv[key]
	if !ok {
		return "", nil
	}
	v, ok := f.p.stringValue(expr, key)
	if !ok {
		return "", fmt.Errorf("%s must be a string literal", key)
	}
	return v, nil
}

func (f *widgetFields) bool(key string) (bool, error) {
	expr, ok := f.kv[key]
	if !ok {
		return false, nil
	}
	id, ok := expr.(*ast.Ident)
	if !ok {
		f.p.errorf("%s must be the literal true or false", key)
		return false, fmt.Errorf("%s must be a boolean literal", key)
	}
	switch id.Name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	f.p.errorf("%s must be the literal true or false, got %q", key, id.Name)
	return false, fmt.Errorf("%s must be a boolean literal", key)
}

func (f *widgetFields) style() ir.Style {
	expr, ok := f.kv["Style"]
	if !ok {
		return ir.Style{}
	}
	return f.p.parseStyle(expr)
}

func (f *widgetFields) css() string {
	v, _ := f.str("CSS")
	return v
}

func (f *widgetFields) children(key string) ([]ir.Widget, error) {
	expr, ok := f.kv[key]
	if !ok {
		return nil, nil
	}
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		f.p.errorf("%s must be a slice literal of ui.Widget values", key)
		return nil, fmt.Errorf("%s must be a slice literal", key)
	}
	var children []ir.Widget
	for _, elt := range compLit.Elts {
		if w := f.p.widget(elt); w != nil {
			children = append(children, w)
		}
	}
	return children, nil
}

func (f *widgetFields) widget(key string) ir.Widget {
	expr, ok := f.kv[key]
	if !ok {
		return nil
	}
	return f.p.widget(expr)
}

// fileImport is an import from the file being parsed.
type fileImport struct {
	path string
	// name is the identifier the import is referenced by: its alias, or the
	// last element of its path.
	name string
}

// fileImports extracts the imports of a file, ignoring blank and dot imports
// whose identifiers cannot be matched against source text.
func fileImports(file *ast.File) []fileImport {
	var out []fileImport
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		name := ""
		if spec.Name != nil {
			name = spec.Name.Name
		} else {
			name = packageIdent(path)
		}
		if name == "" || name == "_" || name == "." {
			continue
		}
		out = append(out, fileImport{path: path, name: name})
	}
	return out
}

// packageIdent is the identifier a package is referenced by in source, which is
// the last path element for an unaliased import.
func packageIdent(path string) string {
	base := pathpkg.Base(path)
	// Module paths may end in a major version suffix, which is not part of the
	// package identifier.
	if i := strings.LastIndexByte(base, 'v'); i > 0 {
		if _, err := strconv.Atoi(base[i+1:]); err == nil {
			base = base[:i]
		}
	}
	return base
}

// importsUsedBy returns the imports referenced as a package qualifier in code,
// so a handler copied into a generated file keeps compiling.
func (p *widgetParser) importsUsedBy(code string) []string {
	var out []string
	for _, imp := range p.imports {
		if regexp.MustCompile(`\b` + regexp.QuoteMeta(imp.name) + `\.`).MatchString(code) {
			out = append(out, imp.path)
		}
	}
	return out
}

// handler records a Go callback used as an event handler.
//
// The callback is stored as a Go expression rather than a name so that both a
// named function and an inline func literal work; the generated events_gen.go
// passes it straight to android.RegisterEvent.
func (f *widgetFields) handler(key, kind, suffix string) string {
	expr, ok := f.kv[key]
	if !ok {
		return ""
	}

	var callback ir.EventHandler
	switch v := expr.(type) {
	case *ast.Ident:
		callback.Expr = v.Name
		kind, ok := f.p.handlerKind(v)
		if !ok {
			f.p.errorf("%s must be a func(), func(bool) or func(string)", key)
			return ""
		}
		callback.Kind = kind
	case *ast.FuncLit:
		rendered, err := f.p.renderExpr(v)
		if err != nil {
			f.p.errorf("cannot use this %s handler: %v", key, err)
			return ""
		}
		callback.Expr = rendered
		callback.Imports = f.p.importsUsedBy(rendered)
		kind, ok := handlerLitKind(v)
		if !ok {
			f.p.errorf("%s must be a func(), func(bool) or func(string)", key)
			return ""
		}
		callback.Kind = kind
	default:
		f.p.errorf("%s must be the name of a Go function or an inline func(){}", key)
		return ""
	}

	widgetID := f.resolvedID
	if widgetID == "" {
		f.p.errorf("%s on a %s requires a usable widget ID so the event can be dispatched", key, kind)
		return ""
	}
	f.p.events[widgetID+"_"+suffix] = callback
	return callback.Expr
}

// handlerKind reports what the named function expects to receive.
func (p *widgetParser) handlerKind(ident *ast.Ident) (ir.HandlerKind, bool) {
	if p.info == nil {
		// Without type information, fall back to the plain form; the Go
		// compiler still rejects a signature the widget's field does not match.
		return ir.HandlerPlain, true
	}
	obj := p.info.ObjectOf(ident)
	if obj == nil {
		return ir.HandlerPlain, true
	}
	return signatureKind(obj.Type())
}

// handlerLitKind reports what an inline func literal expects to receive.
func handlerLitKind(lit *ast.FuncLit) (ir.HandlerKind, bool) {
	params := lit.Type.Params
	if params == nil || params.NumFields() == 0 {
		return ir.HandlerPlain, true
	}
	if params.NumFields() != 1 || len(params.List[0].Names) > 1 {
		return ir.HandlerPlain, false
	}
	id, isIdent := params.List[0].Type.(*ast.Ident)
	if !isIdent {
		return ir.HandlerPlain, false
	}
	switch id.Name {
	case "bool":
		return ir.HandlerChecked, true
	case "string":
		return ir.HandlerText, true
	}
	return ir.HandlerPlain, false
}

// signatureKind reports what a func type is called with.
func signatureKind(t types.Type) (ir.HandlerKind, bool) {
	sig, isSig := t.Underlying().(*types.Signature)
	if !isSig {
		return ir.HandlerPlain, false
	}
	if sig.Params().Len() == 0 {
		return ir.HandlerPlain, true
	}
	if sig.Params().Len() != 1 {
		return ir.HandlerPlain, false
	}
	basic, isBasic := sig.Params().At(0).Type().Underlying().(*types.Basic)
	if !isBasic {
		return ir.HandlerPlain, false
	}
	switch basic.Kind() {
	case types.Bool:
		return ir.HandlerChecked, true
	case types.String, types.UntypedString:
		return ir.HandlerText, true
	}
	return ir.HandlerPlain, false
}

// renderExpr turns an expression back into Go source, so an inline func literal
// written in the UI declaration can be copied into the generated file.
func (p *widgetParser) renderExpr(expr ast.Expr) (string, error) {
	if p.fset == nil {
		return "", errors.New("no file set available")
	}
	var buf strings.Builder
	if err := printer.Fprint(&buf, p.fset, expr); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// stringValue evaluates a string expression, resolving ui.LoadFile calls
// relative to the file being parsed.
func (p *widgetParser) stringValue(expr ast.Expr, what string) (string, bool) {
	if expr == nil {
		return "", true
	}
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind != token.STRING {
			p.errorf("%s must be a string literal", what)
			return "", false
		}
		s, err := strconv.Unquote(v.Value)
		if err != nil {
			p.errorf("%s is not a valid string literal: %v", what, err)
			return "", false
		}
		return s, true
	case *ast.CallExpr:
		if sel, ok := v.Fun.(*ast.SelectorExpr); ok {
			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "ui" && sel.Sel.Name == "LoadFile" {
				if len(v.Args) != 1 {
					p.errorf("ui.LoadFile takes exactly one argument")
					return "", false
				}
				path, ok := p.stringValue(v.Args[0], "ui.LoadFile path")
				if !ok || path == "" {
					return "", false
				}
				return p.loadFile(path, what)
			}
		}
	}
	p.errorf("%s must be a string literal or ui.LoadFile(...)", what)
	return "", false
}

// loadFile reads an asset referenced through ui.LoadFile. Relative paths are
// resolved against the directory of the Go file that references them.
func (p *widgetParser) loadFile(path, what string) (string, bool) {
	resolved := path
	if !filepath.IsAbs(resolved) && p.dir != "" {
		resolved = filepath.Join(p.dir, path)
	}
	b, err := os.ReadFile(resolved)
	if err != nil {
		p.errorf("cannot read %s (%s): %v", what, resolved, err)
		return "", false
	}
	return string(b), true
}

func (p *widgetParser) parseStyle(expr ast.Expr) ir.Style {
	style := ir.Style{}
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		p.errorf("Style must be a ui.Style literal")
		return style
	}
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch key.Name {
		case "BackgroundColor", "TextColor":
			v, ok := p.stringValue(kv.Value, "Style."+key.Name)
			if ok {
				if key.Name == "BackgroundColor" {
					style.BackgroundColor = v
				} else {
					style.TextColor = v
				}
			}
		case "TextSize", "Padding", "Margin":
			v, ok := p.intValue(kv.Value, "Style."+key.Name)
			if ok {
				switch key.Name {
				case "TextSize":
					style.TextSize = v
				case "Padding":
					style.Padding = v
				case "Margin":
					style.Margin = v
				}
			}
		case "Width", "Height":
			v, ok := p.dimensionValue(kv.Value, "Style."+key.Name)
			if ok {
				if key.Name == "Width" {
					style.Width = v
				} else {
					style.Height = v
				}
			}
		case "Weight":
			v, ok := p.floatValue(kv.Value, "Style."+key.Name)
			if ok {
				style.Weight = float32(v)
			}
		case "Animation":
			comp, ok := kv.Value.(*ast.CompositeLit)
			if !ok {
				p.errorf("Style.Animation must be a ui.Animation literal")
				continue
			}
			style.Animation = p.parseAnimation(comp)
		default:
			p.errorf("unknown ui.Style field %q", key.Name)
		}
	}
	return style
}

func (p *widgetParser) intValue(expr ast.Expr, what string) (int, bool) {
	n, err := p.numberValue(expr, what)
	if err != nil {
		p.errorf("%s must be an integer literal: %v", what, err)
		return 0, false
	}
	return int(n), true
}

// dimensionValue accepts plain pixel values, the ui.MatchParent / ui.WrapContent
// constants, and negative literals, all of which the ui package uses for sizing.
func (p *widgetParser) dimensionValue(expr ast.Expr, what string) (int, bool) {
	if sel, ok := selectorName(expr); ok {
		switch sel {
		case "MatchParent":
			return -1, true
		case "WrapContent":
			return -2, true
		}
		p.errorf("%s must be an integer, ui.MatchParent or ui.WrapContent", what)
		return 0, false
	}
	if un, ok := expr.(*ast.UnaryExpr); ok && un.Op == token.SUB {
		n, err := p.numberValue(un.X, what)
		if err != nil {
			p.errorf("%s must be an integer literal: %v", what, err)
			return 0, false
		}
		return -int(n), true
	}
	n, err := p.numberValue(expr, what)
	if err != nil {
		p.errorf("%s must be an integer literal: %v", what, err)
		return 0, false
	}
	return int(n), true
}

func (p *widgetParser) floatValue(expr ast.Expr, what string) (float64, bool) {
	n, err := p.numberValue(expr, what)
	if err != nil {
		p.errorf("%s must be a numeric literal: %v", what, err)
		return 0, false
	}
	return n, true
}

func (p *widgetParser) numberValue(expr ast.Expr, what string) (float64, error) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || (lit.Kind != token.INT && lit.Kind != token.FLOAT) {
		return 0, fmt.Errorf("expected a numeric literal")
	}
	return strconv.ParseFloat(lit.Value, 64)
}

func (p *widgetParser) parseAnimation(compLit *ast.CompositeLit) ir.Animation {
	var anim ir.Animation
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch key.Name {
		case "Type":
			v, ok := p.stringValue(kv.Value, "Animation.Type")
			if ok {
				anim.Type = v
			}
		case "Duration", "Delay":
			v, ok := p.intValue(kv.Value, "Animation."+key.Name)
			if ok {
				if key.Name == "Duration" {
					anim.Duration = v
				} else {
					anim.Delay = v
				}
			}
		case "Loop":
			v, ok := p.boolValue(kv.Value, "Animation.Loop")
			if ok {
				anim.Loop = v
			}
		default:
			p.errorf("unknown ui.Animation field %q", key.Name)
		}
	}
	return anim
}

func (p *widgetParser) boolValue(expr ast.Expr, what string) (bool, bool) {
	id, ok := expr.(*ast.Ident)
	if !ok || (id.Name != "true" && id.Name != "false") {
		p.errorf("%s must be the literal true or false", what)
		return false, false
	}
	return id.Name == "true", true
}

func (p *widgetParser) errorf(format string, args ...any) {
	p.errs = append(p.errs, fmt.Errorf(format, args...))
}

func (p *widgetParser) err() error {
	return errors.Join(p.errs...)
}

func selectorName(expr ast.Expr) (string, bool) {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	return sel.Sel.Name, true
}

// detectAndroidImport finds the import path the app uses for the Go2APK android
// runtime, so generated code refers to the same package the app is built
// against (which may be a fork or a vendored copy).
func detectAndroidImport(pkgs []*packages.Package) string {
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, imp := range file.Imports {
				path, err := strconv.Unquote(imp.Path.Value)
				if err != nil {
					continue
				}
				if path == DefaultAndroidImport || strings.HasSuffix(path, "/Go2APK/android") {
					return path
				}
			}
		}
	}
	return DefaultAndroidImport
}

var (
	javaIdentRE     = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	javaClassNameRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	javaReserved    = map[string]bool{
		"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true,
		"case": true, "catch": true, "char": true, "class": true, "const": true,
		"continue": true, "default": true, "do": true, "double": true, "else": true,
		"enum": true, "extends": true, "final": true, "finally": true, "float": true,
		"for": true, "goto": true, "if": true, "implements": true, "import": true,
		"instanceof": true, "int": true, "interface": true, "long": true, "native": true,
		"new": true, "package": true, "private": true, "protected": true, "public": true,
		"return": true, "short": true, "static": true, "strictfp": true, "super": true,
		"switch": true, "synchronized": true, "this": true, "throw": true, "throws": true,
		"transient": true, "try": true, "void": true, "volatile": true, "while": true,
		"true": true, "false": true, "null": true,
	}
)
