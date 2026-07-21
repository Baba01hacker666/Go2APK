//go:build !android || !cgo

package android

import (
	"io"
	"net/http"
)

func updateTextNative(id string, text string) {}

func setProperty(id, name, value string) {}

func navigate(target string) {}

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

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
