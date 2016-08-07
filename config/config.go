package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var userConfig = struct {
	APItoken string `json:"api_token"`
	GroupID  int64  `json:"group_id"`
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

func APItoken() string {
	return userConfig.APItoken
}

func GroupID() int64 {
	return userConfig.GroupID
}

func NoImages(id int) bool {
	for _, muted := range userConfig.NoImages {
		if id == muted {
			return true
		}
	}
	return false
}
