package android

import "encoding/json"

// VpnConfig represents the configuration for the Android VpnService.
type VpnConfig struct {
	Address     string `json:"address,omitempty"`     // e.g., "10.0.0.2"
	Prefix      int    `json:"prefix,omitempty"`      // e.g., 24
	Route       string `json:"route,omitempty"`       // e.g., "0.0.0.0"
	RoutePrefix int    `json:"routePrefix,omitempty"` // e.g., 0
	MTU         int    `json:"mtu,omitempty"`         // e.g., 1500
	DNS         string `json:"dns,omitempty"`         // e.g., "8.8.8.8"
	Session     string `json:"session,omitempty"`     // Name of the VPN session
}

func (c VpnConfig) toJSON() string {
	b, _ := json.Marshal(c)
	return string(b)
}
