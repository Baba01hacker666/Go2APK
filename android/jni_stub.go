//go:build !android || !cgo

package android

func updateTextNative(id string, text string) {}

func getTextNative(id string) string { return "" }

func requestPermissionNative(permission string, onResult func(granted bool)) {}

func startActivityNative(intent Intent) {}

func startActivityForResultNative(intent Intent, onResult func(resultCode int, data map[string]string)) {
}

func sendBroadcastNative(action string) {}

func sendBroadcastWithExtrasNative(action string, extras map[string]string) {}

func animateNative(id string, property string, to float32, durationMs int) {}
