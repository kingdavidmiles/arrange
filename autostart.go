package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/emersion/go-autostart"
)

func AutoStart() {
	// get exec name
	path, err := os.Executable()
	handleErr(err)

	icon, err := filepath.Abs("assets/icon.png")
	handleErr(err, "icon not found")

	app := &autostart.App{
		Name:        "arrange",
		DisplayName: "Arrange",
		Icon:        icon,
		Exec:        []string{path, "--watch"},
	}

	if err := app.Enable(); err != nil {
		log.Fatal(err)
	}
}
