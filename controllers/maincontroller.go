package maincontroller

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"strings"
	"subtitlewatcher/folderwatcher"
	"subtitlewatcher/messenger"
	"subtitlewatcher/settings"
	"subtitlewatcher/subtitle"
	"subtitlewatcher/systemcontext"
)

var langsSubtitlesAvailable []string

// https://www.loc.gov/standards/iso639-2/php/code_list.php
var langsSubtitleChosen = []string{"pob", "eng"}

var FileFormats = []string{".avi", ".mkv", ".mp4", ".m4v", ".mov", ".mpg", ".wmv"}
var localWatcher *watcher.Watcher

func DownloadSubtitle(errorMsg string, filePath string, onSuccess func(), onError func(err error)) {
	err := subtitle.Get(langsSubtitleChosen, filePath, errorMsg)
	if err != nil {
		onError(err)
	} else {
		onSuccess()
	}
}

func SubtitleWatcherStart(errorMsg string, folderPath string, onSuccess func(folderPath string),
	onError func(err error), onUpdate func(found bool, log string)) {
	localWatcher = folderwatcher.New()
	err := folderwatcher.Watch(localWatcher, FileFormats, folderPath, func(filePath string) {
		fmt.Printf(filePath)
		err := subtitle.Get(langsSubtitleChosen, filePath, errorMsg)
		if err != nil {
			log.Print(err)
			onUpdate(false, err.Error())
		} else {
			onUpdate(true, filePath)
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
	return langsSubtitlesAvailable
}

func Select(selection string, index int) {
	splitSelection := strings.Split(selection, " - ")[0]
	langsSubtitleChosen[index] = splitSelection

	settings.WriteLanguages(langsSubtitleChosen)
}

func FindSelection(index int) string {
	for _, lang := range langsSubtitlesAvailable {
		if strings.Split(lang, " - ")[0] == langsSubtitleChosen[index] {
			return lang
		}
	}
	return ""
}

func SaveSettings() error {
	return settings.WriteConfig()
}

func LoadSettings() error {
	langs := messenger.LanguagesSubtitles()

	for _, lang := range langs {
		langsSubtitlesAvailable = append(langsSubtitlesAvailable, lang["code"]+" - "+lang["language"])
	}

	err := settings.ReadConfig()
	if err != nil {
		return err
	}

	copy(langsSubtitleChosen, settings.Loaded.LangsSubtitle)
	return nil
}

func CheckSystemContextAvailable() bool {
	return systemcontext.Available()
}

func AddSystemContext(onSuccess func(), onError func(err error)) {
	if err := systemcontext.Configure(); err != nil {
		onError(err)
	} else {
		onSuccess()
	}
}
