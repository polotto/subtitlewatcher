package settings

import (
	"encoding/json"
	"subtitlewatcher/ioutil"
)

type Settings struct {
	LangsSubtitle []string
}

var Loaded = Settings{LangsSubtitle: []string{}}
var settingsDir =  ioutil.UserHome() + "/subtitlewatcher"
var settingsOutput = settingsDir + "/settings.json"

func WriteLanguages(langsSutitle []string) {
	Loaded.LangsSubtitle = []string{}
	Loaded.LangsSubtitle = append(Loaded.LangsSubtitle, langsSutitle...)
}

func WriteConfig() error {
	file, err := json.MarshalIndent(Loaded, "", "	")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(settingsOutput, file)
}

func ReadConfig() error {
	file, err := ioutil.ReadFile(settingsOutput)
	if err != nil {
		return ioutil.MakeDirectoryIfNotExists(settingsDir)
	}

	err = json.Unmarshal(file, &Loaded)
	if err != nil {
		return err
	}

	return nil
}