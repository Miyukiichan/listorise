package handlers

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/Miyukiichan/listorise/model"
)

func Config() model.Config {
	file, err := os.Open("conf.json")
	var config model.Config = model.Config{
		DatabasePath: "listorise.db",
		Port: "6573",
	}
	if (err == nil) {
		decoder := json.NewDecoder(file)
		file.Close()
		c := model.Config{}
		err = decoder.Decode(&c)
		if (err == nil) {
			if (config.DatabasePath != "") {
				config.DatabasePath = c.DatabasePath
			}
			if (c.Port != "") {
				_, err := strconv.Atoi(c.Port)
				if (err != nil) {
					log.Fatal(err)
				}
				config.Port = c.Port
			}
		}
	}
	return config
}