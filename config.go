package main

import (
	"MiiContestChannelServer/webpanel"
	"encoding/xml"
	"os"
)

func GetConfig() webpanel.Config {
	data, err := os.ReadFile("config.xml")
	checkError(err)

	var config webpanel.Config
	err = xml.Unmarshal(data, &config)
	checkError(err)

	return config
}
