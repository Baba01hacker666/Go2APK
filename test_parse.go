package main

import (
	"encoding/json"
	"fmt"

	"github.com/Baba01hacker666/Go2APK/frontend"
)

func main() {
	f := frontend.New()
	prog, err := f.BuildIR("/root/Go2APK/examples/demo")
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(prog.UI, "", "  ")
	fmt.Println(string(b))
}
