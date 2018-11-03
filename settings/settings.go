package settings

import (
	"os"

	"io/ioutil"

	"encoding/json"
)

//Settings represents the settings from the settings.json file
type Settings struct {
	Port         int    `json:"port"`
	BeveragePath string `json:"beveragepath"`
	AccountPath  string `json:"accountpath"`
	UserPath     string `json:"userpath"`
	WebappPath   string `json:"webapppath"`
	WebappRoute  string `json:"webapproute"`
	LoggingLevel string `json:"logginglevel"`
}

var (
	//S singleton kinda
	S = Settings{Port: 8080, BeveragePath: "beverages.json", WebappRoute: "/app", WebappPath: "webapp", LoggingLevel: "DEBUG"}
)

//ReadSettings reads the settings file and stores the values in settings.S
func ReadSettings() {
	file, err := os.Open("settings.json")
	if err != nil {
		println("Couldnt open settings file")
		return
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		println("Couldnt read settings file")
		return
	}
	if err := json.Unmarshal(bytes, &S); err != nil {
		println("Couldnt parse settings file")
		println(err.Error())
		return
	}
}
