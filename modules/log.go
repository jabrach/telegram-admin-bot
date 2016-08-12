package modules

import (
	"encoding/json"
	"github.com/jabrach/telegram-admin-bot/cli"
	"log"
)

func Log(msg *cli.Message, wrapper *cli.Wrapper) {
	if marshaled, err := json.MarshalIndent(msg.Data, "", "  "); err == nil {
		log.Printf("%s", marshaled)
	} else {
		log.Fatalln(err.Error())
	}
}
