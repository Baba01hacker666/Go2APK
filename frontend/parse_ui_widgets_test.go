package frontend

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/Baba01hacker666/Go2APK/ir"
)

// newTestParser returns a parser ready to convert a single composite literal.
func newTestParser() *widgetParser {
	return &widgetParser{events: make(map[string]ir.EventHandler)}
}

// parseWidget converts one composite literal such as
// ui.Image{Src: "http://example.com/img.png"} into an IR widget.
func parseWidget(t *testing.T, comp *ast.CompositeLit) ir.Widget {
	t.Helper()
	p := newTestParser()
	w := p.widget(comp)
	if err := p.err(); err != nil {
		t.Fatalf("parsing widget: %v", err)
	}
	return w
}

func TestParseWidgetImage(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Image"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Src"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"http://example.com/img.png"`},
			},
		},
	}

	img, ok := parseWidget(t, comp).(ir.ImageWidget)
	if !ok {
		t.Fatalf("expected ImageWidget")
	}
	if img.Src != "http://example.com/img.png" {
		t.Errorf("Src = %q, want http://example.com/img.png", img.Src)
	}
}

func TestParseWidgetVideo(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Video"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Src"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"file:///sdcard/vid.mp4"`},
			},
		},
	}

	vid, ok := parseWidget(t, comp).(ir.VideoWidget)
	if !ok {
		t.Fatalf("expected VideoWidget")
	}
	if vid.Src != "file:///sdcard/vid.mp4" {
		t.Errorf("Src = %q, want file:///sdcard/vid.mp4", vid.Src)
	}
}

func TestParseWidgetAudio(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Audio"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Src"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"http://example.com/audio.mp3"`},
			},
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "AutoPlay"},
				Value: &ast.Ident{Name: "true"},
			},
		},
	}

	aud, ok := parseWidget(t, comp).(ir.AudioWidget)
	if !ok {
		t.Fatalf("expected AudioWidget")
	}
	if aud.Src != "http://example.com/audio.mp3" {
		t.Errorf("Src = %q, want http://example.com/audio.mp3", aud.Src)
	}
	if !aud.AutoPlay {
		t.Error("AutoPlay = false, want true")
	}
}

func TestParseWidgetScrollView(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "ScrollView"}},
	}

	if _, ok := parseWidget(t, comp).(ir.ScrollViewWidget); !ok {
		t.Fatalf("expected ScrollViewWidget")
	}
}

func TestParseWidgetCardView(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "CardView"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "CSS"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"box-shadow: 4px;"`},
			},
		},
	}

	cv, ok := parseWidget(t, comp).(ir.CardViewWidget)
	if !ok {
		t.Fatalf("expected CardViewWidget")
	}
	if cv.CSS != "box-shadow: 4px;" {
		t.Errorf("CSS = %q, want box-shadow: 4px;", cv.CSS)
	}
}

func TestParseWidgetWebView(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "WebView"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Src"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"https://google.com"`},
			},
		},
	}

	wv, ok := parseWidget(t, comp).(ir.WebViewWidget)
	if !ok {
		t.Fatalf("expected WebViewWidget")
	}
	if wv.Src != "https://google.com" {
		t.Errorf("Src = %q, want https://google.com", wv.Src)
	}
}

func TestParseWidgetXMLView(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "XMLView"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Layout"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"my_layout"`},
			},
		},
	}

	xv, ok := parseWidget(t, comp).(ir.XMLView)
	if !ok {
		t.Fatalf("expected XMLView")
	}
	if xv.Layout != "my_layout" {
		t.Errorf("Layout = %q, want my_layout", xv.Layout)
	}
}

func TestParseWidgetStyleDimensions(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Column"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Style"},
				Value: &ast.CompositeLit{
					Elts: []ast.Expr{
						&ast.KeyValueExpr{Key: &ast.Ident{Name: "Width"}, Value: &ast.SelectorExpr{Sel: &ast.Ident{Name: "MatchParent"}}},
						&ast.KeyValueExpr{Key: &ast.Ident{Name: "Height"}, Value: &ast.SelectorExpr{Sel: &ast.Ident{Name: "WrapContent"}}},
						&ast.KeyValueExpr{Key: &ast.Ident{Name: "Padding"}, Value: &ast.BasicLit{Kind: token.INT, Value: "12"}},
						&ast.KeyValueExpr{Key: &ast.Ident{Name: "Weight"}, Value: &ast.BasicLit{Kind: token.FLOAT, Value: "1.5"}},
					},
				},
			},
		},
	}

	col, ok := parseWidget(t, comp).(ir.ColumnWidget)
	if !ok {
		t.Fatalf("expected ColumnWidget")
	}
	if col.Style.Width != -1 {
		t.Errorf("Width = %d, want -1 (MatchParent)", col.Style.Width)
	}
	if col.Style.Height != -2 {
		t.Errorf("Height = %d, want -2 (WrapContent)", col.Style.Height)
	}
	if col.Style.Padding != 12 {
		t.Errorf("Padding = %d, want 12", col.Style.Padding)
	}
	if col.Style.Weight != 1.5 {
		t.Errorf("Weight = %v, want 1.5", col.Style.Weight)
	}
}

func TestParseWidgetReportsBadStyleValue(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Column"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Style"},
				Value: &ast.CompositeLit{
					Elts: []ast.Expr{
						&ast.KeyValueExpr{Key: &ast.Ident{Name: "Padding"}, Value: &ast.BasicLit{Kind: token.STRING, Value: `"12px"`}},
					},
				},
			},
		},
	}

	p := newTestParser()
	p.widget(comp)
	if err := p.err(); err == nil {
		t.Fatal("expected an error for a non-numeric Padding")
	}
}

func TestParseWidgetUnknownTypeIsReported(t *testing.T) {
	// A widget kind the parser does not know about must be reported rather
	// than silently dropped from the layout.
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Textvew"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Text"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"typo"`},
			},
		},
	}

	p := newTestParser()
	if w := p.widget(comp); w != nil {
		t.Fatalf("expected no widget, got %#v", w)
	}
	err := p.err()
	if err == nil {
		t.Fatal("expected an error for an unknown widget type")
	}
	if !strings.Contains(err.Error(), "unknown widget") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseWidgetRejectsReservedWordID(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "TextView"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "ID"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"class"`},
			},
		},
	}

	p := newTestParser()
	p.widget(comp)
	if err := p.err(); err == nil {
		t.Fatal("expected an error for a Java reserved word used as an ID")
	}
}
