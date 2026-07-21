package android

// Permission declares an Android manifest permission statically.
// Go2APK will automatically inject it into the generated AndroidManifest.xml.
// If it is a dangerous permission (like CAMERA or LOCATION), you must also request it at runtime from
// the user via RequestPermission.
//
// Example:
//
//	android.Permission("android.permission.CAMERA")
//	android.Permission("android.permission.ACCESS_FINE_LOCATION")
func Permission(name string) {
	// Statically analyzed — no runtime behaviour on Android.
	// Stub for non-Android builds.
}

// RequestPermission asks the user to grant a dangerous runtime permission.
// onResult is called with true if the permission was granted, false otherwise.
func RequestPermission(permission string, onResult func(granted bool)) {
	requestPermissionNative(permission, onResult)
}

// ──────────────────────────────────────────────────────────────────────────────
// Intents
// ──────────────────────────────────────────────────────────────────────────────

// Intent represents an Android Intent that can launch activities or services.
type Intent struct {
	// Action is the intent action string, e.g. "android.intent.action.VIEW"
	Action string
	// Data is the URI data for the intent, e.g. "https://example.com"
	Data string
	// Extras are key-value string pairs passed with the intent.
	Extras map[string]string
	// Package is an explicit package name for targeting a specific app (optional).
	Package string
}

// StartActivity launches an Android activity using the given intent.
//
// Example:
//
//	android.StartActivity(android.Intent{
//	    Action: "android.intent.action.VIEW",
//	    Data:   "https://example.com",
//	})
func StartActivity(intent Intent) {
	startActivityNative(intent)
}

// ──────────────────────────────────────────────────────────────────────────────
// Broadcast Receivers
// ──────────────────────────────────────────────────────────────────────────────

var eventRegistry = make(map[string]func())

// RegisterEvent binds an Android event name to a Go handler function.
// This is called automatically by the generated events_gen.go file.
func RegisterEvent(name string, handler func()) {
	eventRegistry[name] = handler
}

// handleEvent is called by the JNI bridge when an event occurs.
// It looks up the registered Go function and executes it.
//
//export handleEvent
func handleEvent(name string) {
	if handler, ok := eventRegistry[name]; ok {
		handler()
	}
}

// BroadcastReceiver declares a system broadcast receiver statically.
// Go2APK will inject this receiver into the AndroidManifest.xml and generate the
// Java stub to route the action back to this callback.
//
// Example:
//
//	android.BroadcastReceiver("android.intent.action.BATTERY_LOW", "on_battery_low", func() {
//	    // Handle battery low
//	})
func BroadcastReceiver(action string, eventName string, callback func()) {
	// Statically analyzed. At runtime we just register the event callback.
	RegisterEvent(eventName, callback)
}

// SendBroadcast sends a simple system broadcast with the given action string.
func SendBroadcast(action string) {
	sendBroadcastNative(action)
}

// SendBroadcastWithExtras sends a broadcast with a map of extra string data.
func SendBroadcastWithExtras(action string, extras map[string]string) {
	sendBroadcastWithExtrasNative(action, extras)
}

// GetText returns the text content of a widget by ID.
func GetText(id string) string {
	return getTextNative(id)
}

// UpdateText sets the text content of a widget by ID.
func UpdateText(id, text string) {
	updateTextNative(id, text)
}

// Animate animates a UI widget. property can be "alpha", "translationX", "translationY", "scaleX", "scaleY".
func Animate(id string, property string, to float32, durationMs int) {
	animateNative(id, property, to, durationMs)
}
