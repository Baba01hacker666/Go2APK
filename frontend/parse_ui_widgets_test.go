package frontend

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/Baba01hacker666/Go2APK/ir"
)

func TestParseCompositeImage(t *testing.T) {
	// Construct an AST for: ui.Image{Src: "http://example.com/img.png"}
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Image"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Src"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"http://example.com/img.png"`},
			},
		},
	}
	events := make(map[string]string)
	widget := parseWidget(comp, events)

	img, ok := widget.(ir.ImageWidget)
	if !ok {
		t.Fatalf("Expected ImageWidget, got %T", widget)
	}
	if img.Src != "http://example.com/img.png" {
		t.Errorf("Expected Src to be http://example.com/img.png, got %s", img.Src)
	}
}

func TestParseCompositeVideo(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Video"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Src"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"file:///sdcard/vid.mp4"`},
			},
		},
	}
	events := make(map[string]string)
	widget := parseWidget(comp, events)

	vid, ok := widget.(ir.VideoWidget)
	if !ok {
		t.Fatalf("Expected VideoWidget, got %T", widget)
	}
	if vid.Src != "file:///sdcard/vid.mp4" {
		t.Errorf("Expected Src to be file:///sdcard/vid.mp4, got %s", vid.Src)
	}
}

func TestParseCompositeAudio(t *testing.T) {
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
	events := make(map[string]string)
	widget := parseWidget(comp, events)

	aud, ok := widget.(ir.AudioWidget)
	if !ok {
		t.Fatalf("Expected AudioWidget, got %T", widget)
	}
	if aud.Src != "http://example.com/audio.mp3" {
		t.Errorf("Expected Src to be http://example.com/audio.mp3, got %s", aud.Src)
	}
	if !aud.AutoPlay {
		t.Errorf("Expected AutoPlay to be true, got false")
	}
}

func TestParseCompositeScrollView(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "ScrollView"}},
	}
	events := make(map[string]string)
	widget := parseWidget(comp, events)

	_, ok := widget.(ir.ScrollViewWidget)
	if !ok {
		t.Fatalf("Expected ScrollViewWidget, got %T", widget)
	}
}

func TestParseCompositeCardView(t *testing.T) {
	comp := &ast.CompositeLit{
		Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "CardView"}},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "CSS"},
				Value: &ast.BasicLit{Kind: token.STRING, Value: `"box-shadow: 4px;"`},
			},
		},
	}
	events := make(map[string]string)
	widget := parseWidget(comp, events)

	cv, ok := widget.(ir.CardViewWidget)
	if !ok {
		t.Fatalf("Expected CardViewWidget, got %T", widget)
	}
	if cv.CSS != "box-shadow: 4px;" {
		t.Errorf("Expected CSS to be parsed, got %s", cv.CSS)
	}
}
