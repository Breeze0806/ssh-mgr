package main

import (
	"flag"
	"fmt"

	"github.com/Breeze0806/ssh-mgr/api/cmdline"
)

var (
	configFile = flag.String("c", "config.json", "config file")
)

func main() {
	flag.Parse()
	defer cmdline.Input()
	e, err := NewEnvironment(*configFile)
	if err != nil {
		fmt.Println("NewEnvironment fail error:", err)
		return
	}
	if err = e.Build(); err != nil {
		fmt.Println("Build fail error:", err)
		return
	}
	if err = e.Run(); err != nil {
		fmt.Println("Run fail error:", err)
		return
	}
}
