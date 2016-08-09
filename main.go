package main

import (
	"flag"
	"fmt"
	"github.com/sthetz/tetanus/cli-wrapper"
	"github.com/sthetz/tetanus/config"
	"github.com/sthetz/tetanus/modules"
)

// 171773961 j
// 122081242 s

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
	defer wrapper.Stop()

	wrapper.AddHandler(modules.NoImages)
	wrapper.AddHandler(modules.Topic.Set)
	wrapper.AddHandler(modules.Topic.Guard)
	wrapper.AddHandler(modules.PicUpdater.Update)
	wrapper.Listen()
}
