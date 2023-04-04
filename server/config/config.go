package config

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var cfg *ini.File

func init() {

	var err error

	cfg, err = ini.Load("./config/config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return
	}
}

func GetConfig(section string, key string) string {
	return cfg.Section(section).Key(key).String()
}
