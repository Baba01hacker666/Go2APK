# Media and Animations

Go2APK natively maps declarative Go UI widgets and API calls to Android's native multimedia and animation frameworks.

## Media Widgets

You can embed Images, Video, and Audio using the `ui.Image`, `ui.Video`, and `ui.Audio` widgets, which map directly to Android's `ImageView`, `VideoView`, and `MediaPlayer`.

### Image

The `Image` widget natively loads images from a URL using Android's `BitmapFactory` in a background thread and applies it to an `ImageView`.

```go
ui.Image{
    ID: "my_image",
    Src: "https://example.com/logo.png",
    CSS: "width: 200px; height: 200px;",
}
```

### Video

The `Video` widget maps to `VideoView` and begins playback automatically. It accepts a standard Android URI.

```go
ui.Video{
    ID: "my_video",
    Src: "file:///sdcard/Download/video.mp4",
    CSS: "width: 100%; height: 300px;",
}
```

### Audio

The `Audio` widget uses a headless `MediaPlayer`. It does not render anything on screen, but can be triggered to auto-play.

```go
ui.Audio{
    ID: "bg_music",
    Src: "https://example.com/soundtrack.mp3",
    AutoPlay: true,
}
```

## HTML Formatted Text

Standard widgets such as `TextView` and `Button` automatically support HTML formatting. Go2APK natively routes their `Text` values through `android.text.Html.fromHtml` using the `FROM_HTML_MODE_COMPACT` flag, enabling basic HTML tags (e.g., `<b>`, `<i>`, `<font color="...">`).

```go
ui.TextView{
    ID: "formatted_text",
    Text: "<b>Bold Text</b> and <i>Italic Text</i>",
}
```

## UI Animations

Go2APK provides direct bindings to Android's powerful `ObjectAnimator` framework. You can trigger arbitrary native property animations on your widgets natively from Go code with zero UI-thread blocking.

Use `android.Animate` and pass the widget ID, property name, target float value, and duration in milliseconds.

```go
import "github.com/Baba01hacker666/Go2APK/android"

func onButtonClick() {
    // Fades out the widget natively over 500ms
    android.Animate("my_image", "alpha", 0.0, 500)

    // Moves the widget 100 pixels to the right over 300ms
    android.Animate("my_image", "translationX", 100.0, 300)
}
```

### Supported Animation Properties
Because `android.Animate` hooks directly into `ObjectAnimator.ofFloat()`, any valid Android float property is supported:
- `"alpha"`: Fade in (1.0) / Fade out (0.0).
- `"translationX"`: Move horizontally.
- `"translationY"`: Move vertically.
- `"scaleX"`: Scale horizontally.
- `"scaleY"`: Scale vertically.
- `"rotation"`: Rotate clockwise (degrees).
