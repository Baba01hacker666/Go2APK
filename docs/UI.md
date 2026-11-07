# Building a UI

Go2APK reads a declarative Go widget tree at build time and generates the
equivalent native Android Views — `LinearLayout`, `TextView`, `EditText`,
`CardView` and so on. Nothing is interpreted on the device: the tree is turned
into Java once, when you build.

```go
package main

import "github.com/Baba01hacker666/Go2APK/ui"

func main() {
	ui.Run(ui.Column{
		Style: ui.Style{Width: ui.MatchParent, Padding: 24},
		Children: []ui.Widget{
			ui.TextView{ID: "title", Text: "Hello"},
			ui.Button{ID: "go", Text: "Go", Style: ui.Style{Width: ui.WrapContent}},
		},
	})
}
```

For a multi-page app, use `ui.RunApp`. The first page is the launcher activity:

```go
ui.RunApp(
	ui.Page{Name: "MainActivity", Root: mainScreen()},
	ui.Page{Name: "SettingsActivity", Root: settingsScreen()},
)
```

---

## Widget reference

| Widget | Android view | Notes |
| --- | --- | --- |
| `Column` | `LinearLayout` (vertical) | `Children` |
| `Row` | `LinearLayout` (horizontal) | `Children` |
| `TextView` | `TextView` | `Text`, accepts HTML |
| `Button` | `Button` | `Text`, `OnClick` |
| `TextField` | `EditText` | `Placeholder`, `OnChanged` |
| `Switch` | `Switch` | `Checked`, `OnChanged` |
| `Image` | `ImageView` | `Src`, loaded off the UI thread |
| `Video` | `VideoView` | `Src` |
| `Audio` | `MediaPlayer` | `Src`, `AutoPlay`. Headless — renders nothing |
| `WebView` | `WebView` | `Src` (URL) or `HTML` (raw markup) |
| `ScrollView` | `ScrollView` | `Children` |
| `CardView` | AndroidX `CardView` | `Children`, elevation and rounded corners |
| `ProgressBar` | `ProgressBar` | |
| `XMLView` | `View` | Inflates your own layout resource by name |

Every widget takes an `ID`, and most take a `Style` and a `CSS` string. See
[Media and animations](MEDIA_AND_ANIMATIONS.md) for the image, video and audio
widgets, and `android.Animate`.

### Sizes

`Style.Width` and `Style.Height` are `ui.MatchParent` (-1), `ui.WrapContent`
(-2), or a size in dp. An unset size falls back to the widget's default, which
is why a `Button` in a `Column` needs a width — see
[Layout warnings](#layout-warnings).

### Style struct

| Field | Type | Meaning |
| --- | --- | --- |
| `BackgroundColor` | `string` | e.g. `"#FF0000"` |
| `TextColor` | `string` | e.g. `"#FFFFFF"` |
| `TextSize` | `int` | sp |
| `Padding` | `int` | dp |
| `Margin` | `int` | dp |
| `Width`, `Height` | `int` | `MatchParent`, `WrapContent`, or dp |
| `Weight` | `float32` | Share of leftover space in a `Row` or `Column` |
| `Animation` | `ui.Animation` | `Type`, `Duration`, `Delay`, `Loop` |

`ui.Animation.Type` is one of `fade_in`, `slide_up`, `bounce`, `pulse`.

---

## Styling with inline CSS

The `CSS` field takes standard CSS declarations that are parsed at build time
and translated into native Java calls. There is no runtime cost and no CSS
engine on the device.

```go
ui.Button{
	ID: "submit_btn",
	Text: "Submit",
	CSS: `
		background-color: #6200ea;
		color: #ffffff;
		font-size: 16px;
		padding: 12px;
		border-radius: 8px;
		font-weight: bold;
	`,
	OnClick: onSubmit,
}
```

You can also load CSS or HTML from a file with `ui.LoadFile`, which is inlined at
build time.

### Supported CSS properties

| CSS property | Example | Maps to Android Java |
|---|---|---|
| `background-color` | `#FF5722`, `red` | A `GradientDrawable` used as the view background |
| `color` | `#FFFFFF`, `black` | `setTextColor(Color.parseColor(...))` |
| `font-size` | `16px`, `24dp`, `18sp` | `setTextSize(Float.parseFloat(...))` |
| `padding` | `16px` | `setPadding(16, 16, 16, 16)` |
| `border-radius` | `8px` | `GradientDrawable.setCornerRadius` |
| `border` | `1px solid black` | `GradientDrawable.setStroke` |
| `box-shadow`, `elevation` | `4px` | `setElevation(...)` |
| `text-align` | `center`, `left`, `right`, `start`, `end` | `setGravity(...)` |
| `font-weight` | `bold` | `setTypeface(Typeface.DEFAULT_BOLD)` |
| `opacity` | `0.5`, `1.0` | `setAlpha(0.5f)` |

`color`, `font-size`, `text-align` and `font-weight` only apply to widgets that
display text. Percentages are not supported; use `px`, `dp` or `sp`.

---

## Widget IDs and events

An `ID` becomes a Java field name and part of the event name dispatched over
JNI, so it must be a valid Java identifier: letters, digits and underscores, not
starting with a digit, and not a Java reserved word. `go2apk check` rejects
anything else rather than emitting Java that will not compile.

Handlers may be a named function or an inline `func() {}`:

```go
ui.Button{ID: "save", Text: "Save", OnClick: save}
ui.Button{ID: "clear", Text: "Clear", OnClick: func() { reset() }}
```

Some widgets pass a value to the handler. The shape of the handler decides how
the event is delivered, so use the signature the widget expects:

| Widget | Field | Handler | Receives |
| --- | --- | --- | --- |
| `Button` | `OnClick` | `func()` | — |
| `Switch` | `OnChanged` | `func(checked bool)` | the new switch state |
| `TextField` | `OnChanged` | `func(text string)` | the new field contents |

```go
ui.Switch{ID: "notify", OnChanged: func(checked bool) { enabled = checked }}
ui.TextField{ID: "search", Placeholder: "Search", OnChanged: func(text string) { filter(text) }}
```

A widget with a handler but no `ID` is given a generated one, because the event
cannot be dispatched without a name to dispatch it by. That works, but if you
want to update the widget from Go with `android.UpdateText` you need to set an
`ID` yourself.

## Updating the UI from Go

```go
android.UpdateText("title", "New text")
current := android.GetText("search")
android.Animate("title", "alpha", 0.0, 500)
```

All three take a widget ID. `UpdateText` and `Animate` are dispatched onto the
Android UI thread by the generated `NativeBridge`, so they are safe to call from
any goroutine. `GetText` reads the view directly, so call it from an event
handler, which already runs on the UI thread.

---

## Previewing your design

```bash
go2apk preview
```

This writes `preview.html`, a Material Design approximation of how your app will
look on a device. It is much faster than building an APK or starting an
emulator, and it is a layout aid only: it does not run your Go code.

## Layout warnings

`go2apk check` and every build report widgets that would be laid out with zero
width. A `Button`, for example, has no intrinsic width, so inside a `Column` it
collapses unless you set one:

```go
ui.Button{Text: "Go", Style: ui.Style{Width: ui.WrapContent}}
```

Inside a `Row`, a zero width is fine as long as the widget has a weight: that is
the usual `Width: 0, Weight: 1` pattern, which lets Android divide the space
between the row's children.
