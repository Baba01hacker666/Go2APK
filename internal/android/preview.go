package android

import (
	"fmt"
	"strings"

	"github.com/Baba01hacker666/Go2APK/internal/config"
	"github.com/Baba01hacker666/Go2APK/ir"
)

// RenderPreviewHTML generates an HTML document that mimics the Android layout of the app.
func RenderPreviewHTML(cfg config.Config, prog *ir.Program) string {
	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	sb.WriteString(fmt.Sprintf("<title>%s Preview</title>\n", cfg.Name))
	sb.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")

	// Basic CSS reset and mimicking Android Material Design
	sb.WriteString("<style>\n")
	sb.WriteString(`
		* { box-sizing: border-box; }
		body { 
			font-family: Roboto, sans-serif; 
			margin: 0; 
			padding: 0; 
			background-color: #f5f5f5; 
			display: flex; 
			justify-content: center; 
			align-items: center; 
			min-height: 100vh;
		}
		.device-screen {
			width: 360px;
			height: 640px;
			background-color: white;
			box-shadow: 0 4px 12px rgba(0,0,0,0.15);
			border-radius: 12px;
			overflow: hidden;
			display: flex;
			flex-direction: column;
			position: relative;
		}
		.app-bar {
			background-color: #6200ea; /* Default purple */
			color: white;
			padding: 16px;
			font-size: 20px;
			font-weight: 500;
			box-shadow: 0 2px 4px rgba(0,0,0,0.2);
			z-index: 10;
		}
		.content {
			flex: 1;
			display: flex;
			flex-direction: column;
			overflow-y: auto;
		}
		/* Go2APK default styles */
		.go2apk-column { display: flex; flex-direction: column; align-items: center; }
		.go2apk-row { display: flex; flex-direction: row; align-items: center; justify-content: center; }
		.go2apk-button {
			background-color: #e0e0e0;
			border: none;
			border-radius: 4px;
			padding: 10px 16px;
			font-size: 14px;
			font-weight: 500;
			text-transform: uppercase;
			cursor: pointer;
			margin: 8px;
			transition: background-color 0.2s;
		}
		.go2apk-button:hover { background-color: #d5d5d5; }
		.go2apk-textview {
			margin: 8px;
			font-size: 14px;
		}
		.go2apk-textfield {
			margin: 8px;
			padding: 8px 12px;
			font-size: 16px;
			border: 1px solid #ccc;
			border-radius: 4px;
			outline: none;
		}
		.go2apk-textfield:focus { border-color: #6200ea; border-width: 2px; padding: 7px 11px; }
	`)
	sb.WriteString("</style>\n</head>\n<body>\n")

	sb.WriteString("<div class=\"device-screen\">\n")
	sb.WriteString(fmt.Sprintf("  <div class=\"app-bar\">%s</div>\n", cfg.Name))
	sb.WriteString("  <div class=\"content\">\n")

	if prog != nil && prog.UI != nil {
		renderWidgetHTML(&sb, prog.UI)
	}

	sb.WriteString("  </div>\n")
	sb.WriteString("</div>\n")
	sb.WriteString("</body>\n</html>\n")
	return sb.String()
}

func renderWidgetHTML(sb *strings.Builder, w ir.Widget) {
	switch v := w.(type) {
	case ir.ColumnWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "column")
		sb.WriteString(fmt.Sprintf("    <div class=\"go2apk-column\" style=%q>\n", inlineCSS))
		for _, child := range v.Children {
			renderWidgetHTML(sb, child)
		}
		sb.WriteString("    </div>\n")
	case ir.RowWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "row")
		sb.WriteString(fmt.Sprintf("    <div class=\"go2apk-row\" style=%q>\n", inlineCSS))
		for _, child := range v.Children {
			renderWidgetHTML(sb, child)
		}
		sb.WriteString("    </div>\n")
	case ir.TextViewWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "textview")
		sb.WriteString(fmt.Sprintf("    <div class=\"go2apk-textview\" style=%q>%s</div>\n", inlineCSS, v.Text))
	case ir.ButtonWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "button")
		sb.WriteString(fmt.Sprintf("    <button class=\"go2apk-button\" style=%q>%s</button>\n", inlineCSS, v.Text))
	case ir.TextFieldWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "textfield")
		sb.WriteString(fmt.Sprintf("    <input type=\"text\" class=\"go2apk-textfield\" placeholder=%q style=%q>\n", v.Placeholder, inlineCSS))
	case ir.ImageWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "image")
		sb.WriteString(fmt.Sprintf("    <img src=%q style=%q>\n", v.Src, inlineCSS))
	case ir.AudioWidget:
		autoPlayStr := ""
		if v.AutoPlay {
			autoPlayStr = " autoplay"
		}
		sb.WriteString(fmt.Sprintf("    <audio src=%q controls%s></audio>\n", v.Src, autoPlayStr))
	case ir.VideoWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "video")
		sb.WriteString(fmt.Sprintf("    <video src=%q style=%q controls></video>\n", v.Src, inlineCSS))
	case ir.WebViewWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "webview")
		if v.Src != "" {
			sb.WriteString(fmt.Sprintf("    <iframe src=%q class=\"go2apk-webview\" style=%q></iframe>\n", v.Src, inlineCSS))
		} else if v.HTML != "" {
			sb.WriteString(fmt.Sprintf("    <iframe srcdoc=%q class=\"go2apk-webview\" style=%q></iframe>\n", v.HTML, inlineCSS))
		}
	case ir.ScrollViewWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "scrollview")
		// Give it a generic overflow-y: auto since it's a ScrollView
		sb.WriteString(fmt.Sprintf("    <div class=\"go2apk-scrollview\" style=\"overflow-y: auto; %s\">\n", inlineCSS))
		for _, child := range v.Children {
			renderWidgetHTML(sb, child)
		}
		sb.WriteString("    </div>\n")
	case ir.CardViewWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "cardview")
		// Mimic material card styling
		sb.WriteString(fmt.Sprintf("    <div class=\"go2apk-cardview\" style=\"box-shadow: 0 4px 8px rgba(0,0,0,0.1); border-radius: 8px; %s\">\n", inlineCSS))
		for _, child := range v.Children {
			renderWidgetHTML(sb, child)
		}
		sb.WriteString("    </div>\n")
	case ir.ProgressBarWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "progressbar")
		// Just a simple HTML5 indeterminate progress element equivalent (a spinner or progress bar)
		sb.WriteString(fmt.Sprintf("    <progress class=\"go2apk-progressbar\" style=%q></progress>\n", inlineCSS))
	case ir.SwitchWidget:
		inlineCSS := generateInlineCSS(v.Style, v.CSS, "switch")
		checkedAttr := ""
		if v.Checked {
			checkedAttr = "checked"
		}
		sb.WriteString(fmt.Sprintf("    <input type=\"checkbox\" class=\"go2apk-switch\" %s style=%q>\n", checkedAttr, inlineCSS))
	}
}

func generateInlineCSS(style ir.Style, customCSS string, widgetType string) string {
	var sb strings.Builder

	// Apply layout params
	if style.Width == -1 {
		sb.WriteString("width: 100%; ")
	}
	if style.Height == -1 {
		sb.WriteString("height: 100%; ")
	}
	if style.Weight > 0 {
		sb.WriteString(fmt.Sprintf("flex: %f; ", style.Weight))
	}

	// Apply padding and margin
	if style.Padding != 0 {
		sb.WriteString(fmt.Sprintf("padding: %dpx; ", style.Padding))
	}
	if style.Margin != 0 {
		sb.WriteString(fmt.Sprintf("margin: %dpx; ", style.Margin))
	}

	// Apply style properties
	if style.BackgroundColor != "" {
		sb.WriteString(fmt.Sprintf("background-color: %s; ", style.BackgroundColor))
	}
	if style.TextColor != "" {
		sb.WriteString(fmt.Sprintf("color: %s; ", style.TextColor))
	}
	if style.TextSize != 0 {
		sb.WriteString(fmt.Sprintf("font-size: %dpx; ", style.TextSize))
	}

	// Append custom user CSS
	if customCSS != "" {
		sb.WriteString(customCSS)
	}

	return sb.String()
}
