package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Baba01hacker666/Go2APK/internal/builder"
	"github.com/Baba01hacker666/Go2APK/internal/config"
	"github.com/Baba01hacker666/Go2APK/internal/doctor"
	"github.com/Baba01hacker666/Go2APK/internal/project"
	"github.com/Baba01hacker666/Go2APK/internal/sdk"
	"github.com/Baba01hacker666/Go2APK/internal/workflow"
)

// loadConfig reads the project configuration from the working directory.
func loadConfig(root string) (config.Config, error) {
	cfg, err := config.Load(filepath.Join(root, "go2apk.yaml"))
	if err != nil {
		return config.Config{}, err
	}
	cfg.ResolveSecrets()
	return cfg, nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "go2apk:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		usage()
		return nil
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	switch args[0] {
	case "init":
		force := false
		if len(args) > 1 && args[1] == "--force" {
			force = true
		}
		return project.Init(root, force)
	case "check":
		return builder.Check(root)
	case "generate":
		return builder.Generate(root)
	case "build":
		opts, err := buildOptions(args[1:])
		if err != nil {
			return err
		}
		return builder.Build(root, opts)
	case "release":
		opts, err := buildOptions(args[1:])
		if err != nil {
			return err
		}
		return builder.Release(root, opts)
	case "clean":
		return builder.Clean(root)
	case "doctor":
		// Doctor runs before a project may exist, so a missing go2apk.yaml is
		// not an error here: fall back to the defaults.
		cfg, err := loadConfig(root)
		if err != nil {
			cfg = config.Default()
		}
		return doctor.Check(os.Stdout, cfg)
	case "workflow":
		if len(args) > 1 && args[1] == "init" {
			cfg, err := loadConfig(root)
			if err != nil {
				return err
			}
			return workflow.Init(root, cfg)
		}
		return fmt.Errorf("unknown workflow command; use: go2apk workflow init")
	case "preview":
		return builder.Preview(root)
	case "sdk":
		if len(args) > 1 && args[1] == "install" {
			cfg, err := loadConfig(root)
			if err != nil {
				return err
			}
			return sdk.Install(root, cfg, os.Stdout)
		}
		return fmt.Errorf("unknown sdk command; use: go2apk sdk install")
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func buildOptions(args []string) (builder.Options, error) {
	var opts builder.Options
	for _, arg := range args {
		switch arg {
		case "--obfuscate":
			opts.Obfuscate = true
		default:
			// Rejecting unknown flags keeps a typo from silently turning into a
			// build with the wrong settings.
			return opts, fmt.Errorf("unknown flag %q", arg)
		}
	}
	return opts, nil
}

func usage() {
	fmt.Print(`Go2APK turns a declarative Go UI application into a native Android APK.

Usage:
  go2apk init [--force]   Write go2apk.yaml, the Android project, scripts and CI
  go2apk check            Parse the app and report layout warnings
  go2apk generate         Regenerate the Android project without running Gradle
  go2apk preview          Render preview.html to check the layout in a browser
  go2apk build [--obfuscate]
                          Build a debug APK into dist/debug
  go2apk release [--obfuscate]
                          Build a release APK and AAB into dist/release
  go2apk sdk install      Write the Android SDK installer scripts
  go2apk workflow init    Regenerate the GitHub Actions workflows
  go2apk doctor           Report which build prerequisites are present
  go2apk clean            Remove generated dist artifacts

Run each command from the directory holding go2apk.yaml.
`)
}
