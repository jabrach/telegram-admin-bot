package main

import (
	"flag"
	"fmt"
	"github.com/sthetz/tetanus/cli-wrapper"
	"github.com/sthetz/tetanus/config"
	"github.com/sthetz/tetanus/modules"
)

func main() {
	var configPath = flag.String("C", "", "Path to config")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("Usage: ./tetanus -C path/to/config.json")
		return
	}
	if err := config.Load(*configPath); err != nil {
		panic(err)
	}

	wrapper := cli.New()
	wrapper.AddHandler(modules.NoImages)
	wrapper.Listen()
}
