package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var userConfig = struct {
	BinPath  string `json:"telegram_cli_path"`
	GroupID  int    `json:"group_id"`
	NoImages []int  `json:"no_images"`
}{}

func Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Unable to read config file: %v", err.Error())
	}
	if err := json.Unmarshal(data, &userConfig); err != nil {
		return fmt.Errorf("Malformed JSON in config: %v", err.Error())
	}

	return nil
}

func GroupID() int {
	return userConfig.GroupID
}

func BinPath() string {
	return userConfig.BinPath
}

func NoImages(id int) bool {
	for _, muted := range userConfig.NoImages {
		if id == muted {
			return true
		}
	}
	return false
}
