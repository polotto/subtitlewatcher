package folderwatcher

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"regexp"
	"strings"
	"time"
)

func New() *watcher.Watcher {
	return watcher.New()
}

func Watch(w *watcher.Watcher, fileFormats []string, folderPath string, fileChanged func(filePath string)) error {
	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(100)

	// Only notify rename and move events.
	w.FilterOps(watcher.Create)

	// Only files that match the regular expression during file listings
	// will be watched.
	regexStr := strings.Join(fileFormats, "|")
	r := regexp.MustCompile("(" + regexStr + ")")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))
	w.IgnoreHiddenFiles(true)

	if err := w.AddRecursive(folderPath); err != nil {
		log.Fatalln(err)
		return err
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event)
				fileChanged(event.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	go func() {
		// Start the watching process - it'll check for changes every 100ms.
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()

	return nil
}

func Stop(folderWatcher *watcher.Watcher) {
	folderWatcher.Wait()
	folderWatcher.Close()
}