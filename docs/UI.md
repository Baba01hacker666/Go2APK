# Go2APK UI Styling Guide

Go2APK translates declarative Go UI structs into native Android Views (`LinearLayout`, `TextView`, `Button`, `EditText`).

You can style these views using either the `Style` struct or standard inline `CSS`.

## Using Inline CSS

The most powerful way to style your app is by using the `CSS` field on widgets. This lets you write standard CSS rules that are parsed at build-time and translated directly into native Java Android calls. 

**Note: There is ZERO performance penalty at runtime. The CSS is converted to native Java methods (e.g. `setBackgroundColor`, `setPadding`) during the build process.**

### Supported CSS Properties

| CSS Property | Example | Maps to Android Java |
|---|---|---|
| `background-color` | `#FF5722`, `red`, `#000000` | `view.setBackgroundColor(Color.parseColor(...))` |
| `color` | `#FFFFFF`, `black` | `textView.setTextColor(Color.parseColor(...))` |
| `font-size` | `16px`, `24dp`, `18sp` | `textView.setTextSize(Float.parseFloat(...))` |
| `padding` | `16px` | `view.setPadding(16, 16, 16, 16)` |
| `border-radius` | `8px` | Generates a `GradientDrawable` with `setCornerRadius` |
| `text-align` | `center`, `left`, `right`, `start`, `end` | `textView.setGravity(Gravity.CENTER)` |
| `font-weight` | `bold` | `textView.setTypeface(Typeface.DEFAULT_BOLD)` |
| `opacity` | `0.5`, `1.0` | `view.setAlpha(0.5f)` |

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
	OnClick: func() {
		// ...
	},
}
```

## Previewing Your Design
Instead of building the entire Android APK (which takes time) or starting an emulator, you can instantly preview your layout!

Run this in your project root:
```bash
go2apk preview
```

This will generate a `preview.html` file in your directory. Open it in any web browser to see a Material Design approximation of how your app will look on an Android device!
