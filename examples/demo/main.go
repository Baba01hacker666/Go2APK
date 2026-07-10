//go:build android

// Package main is a tiny Go mobile demo app used by the Go2APK CI workflow.
package main

import (
	"image/color"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/gl"
)

var background = color.RGBA{R: 0x21, G: 0x96, B: 0xF3, A: 0xFF}

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
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
				if glctx != nil {
					glctx.Viewport(0, 0, event.WidthPx, event.HeightPx)
				}
			case paint.Event:
				if glctx == nil || event.External {
					continue
				}
				draw(glctx)
				a.Publish()
			}
		}
	})
}

func draw(glctx gl.Context) {
	glctx.ClearColor(float32(background.R)/255, float32(background.G)/255, float32(background.B)/255, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
}
