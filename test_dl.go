package main

import (
	"fmt"
	"github.com/Baba01hacker666/Go2APK/internal/gradle"
)

func main() {
	bin, err := gradle.EnsureGradle("/tmp/testapp")
	fmt.Println(bin, err)
}
