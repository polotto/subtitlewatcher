package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"subtitlewatcher/folderwatcher"
	"subtitlewatcher/subtitledb"
)

var FileFormats = []string{".avi", ".mkv", ".mp4", ".m4v", ".mov", ".mpg", ".wmv"}
var languages = []string{"pt", "en"}
var localWatcher *watcher.Watcher

func downloadSubtitle(filePath string, onSuccess func(), onError func(err error)) {
	err := subtitledb.Get(languages, filePath)
	if err != nil {
		onError(err)
	} else {
		onSuccess()
	}
}

func subtitleWatcherStart(folderPath string, onSuccess func(folderPath string), onError func(err error)) {
	localWatcher = folderwatcher.New()
	err := folderwatcher.Watch(localWatcher, FileFormats, folderPath, func(filePath string) {
		fmt.Printf(filePath)
		err := subtitledb.Get(languages, filePath)
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

func subtitleWatcherStop(onSuccess func()) {
	folderwatcher.Stop(localWatcher)
	onSuccess()
}
