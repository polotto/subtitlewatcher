package maincontroller

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"strings"
	"subtitlewatcher/folderwatcher"
	"subtitlewatcher/ioutil"
	"subtitlewatcher/messenger"
	"subtitlewatcher/subtitle"
)

var langsAvailable []string
// https://www.loc.gov/standards/iso639-2/php/code_list.php
var langsSubtitle = []string{"pb", "en"}
var settingsOutput = "./settings.txt"

var FileFormats = []string{".avi", ".mkv", ".mp4", ".m4v", ".mov", ".mpg", ".wmv"}
var localWatcher *watcher.Watcher

func DownloadSubtitle(errorMsg string, filePath string, onSuccess func(), onError func(err error)) {
	err := subtitle.Get(langsSubtitle, filePath, errorMsg)
	if err != nil {
		onError(err)
	} else {
		onSuccess()
	}
}

func SubtitleWatcherStart(errorMsg string, folderPath string, onSuccess func(folderPath string), onError func(err error)) {
	localWatcher = folderwatcher.New()
	err := folderwatcher.Watch(localWatcher, FileFormats, folderPath, func(filePath string) {
		fmt.Printf(filePath)
		err := subtitle.Get(langsSubtitle, filePath, errorMsg)
		if err != nil {
			log.Print(err)
		}
	})
	if err != nil {
		onError(err)
	} else {
		onSuccess(folderPath)
	}
}

func SubtitleWatcherStop(onSuccess func()) {
	folderwatcher.Stop(localWatcher)

	if onSuccess != nil {
		onSuccess()
	}
}

func Languages() []string {
	return langsAvailable
}

func Select(selection string, index int) {
	splitSelection := strings.Split(selection, " - ")[0]
	langsSubtitle[index] = splitSelection
}

func FindSelection(index int) string {
	for _, lang := range langsAvailable {
		if strings.Split(lang, " - ")[0] == langsSubtitle[index] {
			return lang
		}
	}
	return ""
}

func SaveSettings() {
	_ = ioutil.WriteLines(langsSubtitle, settingsOutput)
}

func LoadSettings() {
	langs := messenger.Languages()

	for _, lang := range langs {
		langsAvailable = append(langsAvailable, lang["code"] + " - " + lang["language"])
	}

	settings, err := ioutil.ReadLines(settingsOutput)
	if err != nil {
		return
	}

	copy(langsSubtitle, settings)
}

func remove(slice []string, search string) []string {
	index := 0
	for i, s := range slice {
		if s == search {
			index = i
			break
		}
	}

	return append(slice[:index], slice[index+1:]...)
}