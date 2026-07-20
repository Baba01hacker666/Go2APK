package main

import (
	"fmt"
	"os"

	"github.com/Baba01hacker666/Go2APK/internal/builder"
	"github.com/Baba01hacker666/Go2APK/internal/doctor"
	"github.com/Baba01hacker666/Go2APK/internal/project"
	"github.com/Baba01hacker666/Go2APK/internal/sdk"
	"github.com/Baba01hacker666/Go2APK/internal/workflow"
)

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
		return project.Init(root)
	case "check":
		return builder.Check(root)
	case "build":
		return builder.Build(root, buildOptions(args[1:]))
	case "release":
		return builder.Release(root, buildOptions(args[1:]))
	case "clean":
		return builder.Clean(root)
	case "doctor":
		return doctor.Check(os.Stdout)
	case "workflow":
		if len(args) > 1 && args[1] == "init" {
			return workflow.Init(root)
		}
		return fmt.Errorf("unknown workflow command; use: go2apk workflow init")
	case "sdk":
		if len(args) > 1 && args[1] == "install" {
			return sdk.Install(root, os.Stdout)
		}
		return fmt.Errorf("unknown sdk command; use: go2apk sdk install")
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func buildOptions(args []string) builder.Options {
	var opts builder.Options
	for _, arg := range args {
		if arg == "--obfuscate" {
			opts.Obfuscate = true
		}
	}
	return opts
}

func usage() {
	fmt.Println(`Go2APK converts Go projects into Android APK scaffolding.

Usage:
  go2apk init           Generate go2apk.yaml, Android templates, and helper scripts
  go2apk check          Check declarative UI syntax and logic
  go2apk build [--obfuscate]
                         Validate inputs and run a debug Gradle build when available
  go2apk release [--obfuscate]
                         Validate inputs and run an obfuscated release Gradle build when requested
  go2apk sdk install    Install Android command-line SDK tools into .go2apk/android-sdk
  go2apk workflow init  Generate GitHub Actions workflows
  go2apk doctor         Check Go, Java, Gradle, and Android SDK tooling
  go2apk clean          Remove generated dist artifacts`)
}
