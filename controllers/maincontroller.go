package maincontroller

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"subtitlewatcher/folderwatcher"
	"subtitlewatcher/subtitle"
)

// https://www.loc.gov/standards/iso639-2/php/code_list.php
var languages = []string{"pob", "eng"}
var FileFormats = []string{".avi", ".mkv", ".mp4", ".m4v", ".mov", ".mpg", ".wmv"}
var localWatcher *watcher.Watcher

func DownloadSubtitle(errorMsg string, filePath string, onSuccess func(), onError func(err error)) {
	err := subtitle.Get(languages, filePath, errorMsg)
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
		err := subtitle.Get(languages, filePath, errorMsg)
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