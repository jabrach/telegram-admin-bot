package main

import (
	"flag"
	"fmt"
	"github.com/jabrach/telegram-admin-bot/cli"
	"github.com/jabrach/telegram-admin-bot/config"
	"github.com/jabrach/telegram-admin-bot/modules"
)

// 171773961 j
// 122081242 s

func main() {
	var configPath = flag.String("C", "", "Path to config")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("Usage: ./telegram-admin-bot -C path/to/config.json")
		return
	}
	if err := config.Load(*configPath); err != nil {
		panic(err)
	}

	wrapper := cli.New()
	defer wrapper.Stop()

	// wrapper.AddHandler(modules.Log)
	wrapper.AddHandler(modules.Mute)
	wrapper.AddHandler(modules.Topic.Set)
	wrapper.AddHandler(modules.Topic.Guard)
	wrapper.AddHandler(modules.PicUpdater.Update)
	wrapper.Start()
}
