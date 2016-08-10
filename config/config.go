package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var cnf = struct {
	BinPath       string   `json:"telegram_cli_path"`
	ManagedGroups []jGroup `json:"managed_groups"`
	MGroupsMap    map[int64]*Group
}{}

func Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Unable to read config file: %v", err.Error())
	}
	if err := json.Unmarshal(data, &cnf); err != nil {
		return fmt.Errorf("Malformed JSON in config: %v", err.Error())
	}

	cnf.MGroupsMap = map[int64]*Group{}
	for _, group := range cnf.ManagedGroups {
		cnf.MGroupsMap[group.ID] = group.Init()
	}

	return nil
}

func ManagedGroup(id int64) *Group {
	return cnf.MGroupsMap[id]
}

func BinPath() string {
	return cnf.BinPath
}
