# Go2APK Feature Expansion Plan

## 1. Optimize Startup and Binary Size (ProGuard & Minify)
- [ ] Add `-ldflags="-s -w"` to the `go build -buildmode=c-shared` command in `internal/builder/builder.go`.
- [ ] Enable `minifyEnabled true` and `shrinkResources true` by default for release builds in `internal/android/gradle_templates.go`.

## 2. HTML Text Support
- [ ] Update `internal/android/dynamic_templates.go`: `TextView.setText` and `updateWidgetText` should parse HTML if the text contains HTML tags (e.g. `<b>`, `<i>`) using `android.text.Html.fromHtml(..., android.text.Html.FROM_HTML_MODE_COMPACT)`.
- [ ] Update `android/jni_android.go` and Java bridge to properly handle HTML updates.

## 3. Media Widgets (Image, Audio, Video)
- [ ] Add `Image`, `Audio`, and `Video` structs to `ui/ui.go`.
- [ ] Add corresponding `ImageWidget`, `AudioWidget`, `VideoWidget` to `ir/ir.go`.
- [ ] Update `frontend/parse_ui.go` to parse these new widgets and their `Src` properties.
- [ ] Update `internal/android/dynamic_templates.go` to generate Android `ImageView`, `VideoView`, and `MediaPlayer`.
- [ ] Update `internal/android/preview.go` to support `<img>`, `<audio>`, and `<video>` tags.

## 4. Animation Framework
- [ ] Add `android.Animate(id string, property string, to float32, durationMs int)` to `android/android.go`.
- [ ] Add JNI bindings in `android/jni_android.go` and `android/jni_stub.go`.
- [ ] Update `NativeBridge.java` to invoke Android `ObjectAnimator` for properties like `"alpha"`, `"translationX"`, `"scaleX"`, etc.
