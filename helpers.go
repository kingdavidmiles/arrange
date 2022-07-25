package main

import (
	_ "embed"
	"log"

	"github.com/gen2brain/beeep"
)

// Desktop notifications
func NotifySys(t, msg string) {
	if t == "" {
		t = "Arrange Notifier"
	}
	if err := beeep.Notify(t, msg, "assets/information.png"); err != nil {
		handleErr(err)
	}
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// handle & exit of an error if met
func handleErr(e error, msg ...string) {
	if e != nil {
		log.Fatalln(e, msg)
	}
}
