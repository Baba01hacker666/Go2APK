//go:build android

// Package main is a Go-only Android GUI demo app used by Go2APK builds.
package main

import (
	"image/color"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

type theme struct {
	top    color.RGBA
	accent color.RGBA
}

var themes = []theme{
	{top: color.RGBA{R: 0x21, G: 0x96, B: 0xF3, A: 0xFF}, accent: color.RGBA{R: 0xFF, G: 0xC1, B: 0x07, A: 0xFF}},
	{top: color.RGBA{R: 0x67, G: 0x3A, B: 0xB7, A: 0xFF}, accent: color.RGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}},
	{top: color.RGBA{R: 0x00, G: 0x96, B: 0x88, A: 0xFF}, accent: color.RGBA{R: 0xFF, G: 0x57, B: 0x22, A: 0xFF}},
}

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var width, height int
		selected := 0
		lastFrame := time.Now()

		for event := range a.Events() {
			switch event := a.Filter(event).(type) {
			case lifecycle.Event:
				if event.To == lifecycle.StageDead {
					return
				}
				if ctx, ok := event.DrawContext.(gl.Context); ok {
					glctx = ctx
				}
			case size.Event:
				width, height = event.WidthPx, event.HeightPx
				if glctx != nil {
					glctx.Viewport(0, 0, width, height)
				}
			case touch.Event:
				if event.Type == touch.TypeEnd {
					selected = (selected + 1) % len(themes)
					a.Send(paint.Event{})
				}
			case paint.Event:
				if glctx == nil || event.External {
					continue
				}
				draw(glctx, themes[selected], width, height, time.Since(lastFrame))
				lastFrame = time.Now()
				a.Publish()
			}
		}
	})
}

func draw(glctx gl.Context, active theme, width, height int, frameAge time.Duration) {
	// This Go-only GUI intentionally avoids Java views. It paints a touch-responsive
	// OpenGL surface managed by gomobile; tapping the app cycles the interface theme.
	pulse := float32((frameAge.Milliseconds()%1000)+300) / 1300
	bg := mix(active.top, color.RGBA{R: 0x12, G: 0x12, B: 0x12, A: 0xFF}, 0.22)
	if width > height {
		bg = mix(bg, active.accent, 0.12*pulse)
	}
	glctx.ClearColor(channel(bg.R), channel(bg.G), channel(bg.B), 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
}

func mix(a, b color.RGBA, amount float32) color.RGBA {
	inv := 1 - amount
	return color.RGBA{
		R: uint8(float32(a.R)*inv + float32(b.R)*amount),
		G: uint8(float32(a.G)*inv + float32(b.G)*amount),
		B: uint8(float32(a.B)*inv + float32(b.B)*amount),
		A: 0xFF,
	}
}

func channel(v uint8) float32 { return float32(v) / 255 }
