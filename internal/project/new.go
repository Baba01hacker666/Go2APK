package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewApp creates a beginner-friendly Go2APK app with UI and logic separated.
func NewApp(root, name string, force bool) error {
	if strings.TrimSpace(name) == "" {
		name = "my-go-app"
	}
	dir := filepath.Join(root, name)
	if filepath.IsAbs(name) {
		dir = name
	}
	files := map[string]string{
		"go.mod":      fmt.Sprintf("module %s\n\ngo 1.25.0\n\nrequire github.com/Baba01hacker666/Go2APK v0.0.0\n\nreplace github.com/Baba01hacker666/Go2APK => %s\n", moduleName(name), filepath.ToSlash(root)),
		"go2apk.yaml": fmt.Sprintf("name: %s\npackage: com.example.%s\nversion: 0.1.0\nmin_sdk: 23\ntarget_sdk: 36\norientation: unspecified\ntheme: @style/AppTheme\nsource: .\nobfuscate: false\n# Add Android Maven dependencies here, for example:\n# android_dependencies: com.google.android.material:material:1.12.0\n", titleName(name), packageSuffix(name)),
		"main.go": `package main

import "github.com/Baba01hacker666/Go2APK/ui"

func main() {
	ui.Run(ui.Build(HomeScreen))
}
`,
		"home_ui.go": `package main

import "github.com/Baba01hacker666/Go2APK/ui"

func HomeScreen() ui.Widget {
	return ui.Column{
		ID: "home",
		Style: ui.Style{Padding: 24, Width: ui.MatchParent, Height: ui.MatchParent},
		Children: []ui.Widget{
			ui.TextView{ID: "title", Text: Title(), Style: ui.Style{TextSize: 24, Margin: 12}},
			ui.Button{ID: "hello_button", Text: "Tap me", OnClick: OnHelloTapped},
		},
	}
}
`,
		"logic.go": `package main

import "github.com/Baba01hacker666/Go2APK/android"

func Title() string { return "Hello from Go2APK" }

func OnHelloTapped() {
	android.Toast("Hello from separated logic!")
	android.UpdateText("title", "You tapped the button")
}
`,
	}
	for name, contents := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if _, err := os.Stat(path); err == nil && !force {
			return fmt.Errorf("%s already exists; use --force", path)
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			return err
		}
	}
	fmt.Printf("created Go2APK app in %s\n", dir)
	return nil
}

func moduleName(name string) string { return "example.com/" + packageSuffix(name) }
func titleName(name string) string  { return strings.Title(strings.ReplaceAll(name, "-", " ")) }
func packageSuffix(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "goapp"
	}
	return b.String()
}
