package config

// Config contains the project metadata used to generate Android artifacts.
type Config struct {
	Name        string
	Package     string
	Version     string
	MinSDK      int
	TargetSDK   int
	Orientation string
	Theme       string
}

// Default returns beginner-friendly Android application defaults.
func Default() Config {
	return Config{
		Name:        "Go2APKApp",
		Package:     "com.example.go2apkapp",
		Version:     "0.1.0",
		MinSDK:      23,
		TargetSDK:   35,
		Orientation: "unspecified",
		Theme:       "@style/AppTheme",
	}
}
