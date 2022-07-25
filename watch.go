package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/radovskyb/watcher"
)

func WatchDir(dir string) {
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)

	// Only notify rename and create events.
	w.FilterOps(watcher.Create)

	go func() {
		for {
			select {
			case event := <-w.Event:
				// get the base dir of the file
				folder := filepath.Dir(event.Path)
				noOfFiles := 0

				LoopAndMove(folder, event.FileInfo, &noOfFiles)

				// notifiy the user that this single file has been moved
				msg := fmt.Sprintf("Moved %s into %s", event.FileInfo.Name(), folder)
				NotifySys("", msg)

			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(dir); err != nil {
		log.Fatalln(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Second * 3); err != nil {
		log.Fatalln(err)
	}
}
