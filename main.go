package main

import (
	"flag"
	"fmt"
	"github.com/sthetz/tetanus/botapi"
	"github.com/sthetz/tetanus/config"
)

const BotName = "Tetanus"

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

	bot := botapi.New(config.APItoken())
	bot.Listen()
}
