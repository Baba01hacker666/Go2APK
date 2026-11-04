# Go2APK UI Styling Guide

Go2APK translates declarative Go UI structs into native Android Views (`LinearLayout`, `TextView`, `Button`, `EditText`).

You can style these views using either the `Style` struct or standard inline `CSS`.

## Using Inline CSS

The most powerful way to style your app is by using the `CSS` field on widgets. This lets you write standard CSS rules that are parsed at build-time and translated directly into native Java Android calls. 

There is no runtime cost: the CSS is parsed at build time and turned into native
Java calls in the generated Activity, so nothing is interpreted on the device.

### Supported CSS Properties

| CSS Property | Example | Maps to Android Java |
|---|---|---|
| `background-color` | `#FF5722`, `red`, `#000000` | A `GradientDrawable` used as the view background |
| `color` | `#FFFFFF`, `black` | `textView.setTextColor(Color.parseColor(...))` |
| `font-size` | `16px`, `24dp`, `18sp` | `textView.setTextSize(Float.parseFloat(...))` |
| `padding` | `16px` | `view.setPadding(16, 16, 16, 16)` |
| `border-radius` | `8px` | `GradientDrawable.setCornerRadius` |
| `border` | `1px solid black` | `GradientDrawable.setStroke` |
| `box-shadow`, `elevation` | `4px` | `view.setElevation(...)` |
| `text-align` | `center`, `left`, `right`, `start`, `end` | `textView.setGravity(...)` |
| `font-weight` | `bold` | `textView.setTypeface(Typeface.DEFAULT_BOLD)` |
| `opacity` | `0.5`, `1.0` | `view.setAlpha(0.5f)` |

`color`, `font-size`, `text-align` and `font-weight` only apply to widgets that
display text. Percentages are not supported; use `px`, `dp` or `sp`.

### Example
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

## Previewing Your Design
Instead of building the entire Android APK (which takes time) or starting an emulator, you can instantly preview your layout!

Run this in your project root:
```bash
go2apk preview
```

This will generate a `preview.html` file in your directory. Open it in any web browser to see a Material Design approximation of how your app will look on an Android device!

## Widget IDs and events

An `ID` is used as a Java field name and as part of the event name that is
dispatched over JNI, so it must be a valid Java identifier: letters, digits and
underscores, not starting with a digit, and not a Java reserved word. `go2apk
check` rejects anything else rather than emitting Java that will not compile.

Event handlers may be a named function or an inline `func() {}`:

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
