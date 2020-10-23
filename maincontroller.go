package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"os/user"
	"subtitlewatcher/folderwatcher"
	"subtitlewatcher/subtitledb"
)

func userInfo() (string, string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return "", "", err
	}
	return  usr.HomeDir, usr.Name, nil
}

var FileFormats = []string{".avi", ".mkv", ".mp4", ".m4v", ".mov", ".mpg", ".wmv"}
var languages = []string{"pt", "en"}
var watchStarted = false
var localWatcher *watcher.Watcher

func downloadSubtitle(filePath string, onSuccess func(), onError func(err error)) {
	err := subtitledb.Get(languages, filePath)
	if err != nil {
		onError(err)
	} else {
		onSuccess()
	}
}

func subtitleWatcher(enabled func(home string, name string), disabled func(), onError func(err error)) {
	watchStarted = !watchStarted
	if watchStarted {
		home, name, err := userInfo()
		if err != nil {
			onError(err)
		} else {
			localWatcher = folderwatcher.New()
			err := folderwatcher.Watch(localWatcher, FileFormats, home, func(filePath string) {
				fmt.Printf(filePath)
				err := subtitledb.Get(languages, filePath)
				if err != nil {
					log.Print(err)
				}
			})
			if err != nil {
				onError(err)
			} else {
				enabled(home, name)
			}
		}
	} else {
		folderwatcher.Stop(localWatcher)
		disabled()
	}
}