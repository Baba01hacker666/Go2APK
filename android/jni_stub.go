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

// StartVPN requests Android VPN permissions. If granted, it starts the VPN
// service with the given config. The established TUN file descriptor (fd)
// is sent back to the onEstablished callback.
func StartVPN(config VpnConfig, onEstablished func(fd int)) {
	// Stub implementation
}
