package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var arrangeDescription = ` Arrange is a commmand tool that groups files in a directory into folders according to their file extensions.
 for example files with extension .jpg, .png & .gif are grouped into Pictures
 and files with extension .mp4, .avi, and .mkv are grouped into Videos
`

type Arrange struct {
	Folder     string
	Extensions []string `json:"extensions"`
}

// Create an array to represent all kind of files we can meet during our
// transversing
var stuff = []Arrange{
	{
		Folder:     "videos",
		Extensions: []string{"mp4", "avi,", "mkv", "webm"},
	},
	{
		Folder:     "images",
		Extensions: []string{"jpg", "jpeg", "png", "gif", "webp"},
	},
	{
		Folder:     "vectors",
		Extensions: []string{"svg"},
	},
	{
		Folder:     "musics",
		Extensions: []string{"mp3", "wav", "oob"},
	},
	{
		Folder:     "documents",
		Extensions: []string{"pdf", "odt", "txt", "doc", "docx"},
	},
	{
		Folder:     "applications",
		Extensions: []string{"exe", "appimage"},
	},
	{
		Folder:     "archieves",
		Extensions: []string{"zip", "xz", "gz"},
	},
	{
		Folder:     "debs",
		Extensions: []string{"deb", "sh"},
	},
}

func main() {
	if runtime.GOOS == "windows" {
		panic("We currently don't support windows at the moment!")
	}

	//	get users home dir
	homedir, err := os.UserHomeDir()
	handleErr(err)

	downloadsDir := fmt.Sprintf("%s/Downloads", homedir)

	// process comand line argument/flags
	folderPtr := flag.String("path", downloadsDir, "path to arrange")
	watchPtr := flag.Bool("watch", false, "watch for changes in this path")
	flag.Parse()

	var Usage = func() {
		fmt.Fprintf(os.Stderr, "%s %s:\n\n", arrangeDescription, os.Args[0])

		flag.PrintDefaults()
	}

	if cap(os.Args) <= 1 {
		Usage()
		os.Exit(1)
	}

	// convert flags to normal strings
	var (
		folder      = *folderPtr
		shouldWatch = *watchPtr
	)

	// check the folder existence
	_, err = os.Stat(folder)
	if os.IsNotExist(err) {
		handleErr(err, "Folder does not exist.")
	}

	// walk through the files in the `folder`
	f, err := ioutil.ReadDir(folder)
	handleErr(err)

	// create a variable to hold the number of files moved
	noOfFiles := 0

	// loop through all the files in the directory
	for _, file := range f {
		// if the file is a directory, skip it ..
		// .. we don;t care about dirs
		if file.IsDir() {
			continue
		}

		// for each file in the directory,
		// loop through our list of stuff and move the file
		LoopAndMove(folder, file, &noOfFiles)
	}

	if noOfFiles > 0 {
		msg := fmt.Sprintf("Moved %d files", noOfFiles)
		// notify user that we're done witth moving
		NotifySys("", msg)
	}

	if shouldWatch {
		// block and watch the folder
		WatchDir(folder)
	}

}

func LoopAndMove(folder string, file fs.FileInfo, noOfFilesPtr *int) {

	// loop through our stuff
	for _, s := range stuff {
		// extract the file extension of the current file
		fileExt := filepath.Ext(file.Name())
		// remove . from the file extension
		fileExt = strings.Replace(fileExt, ".", "", 1)
		fileExt = strings.ToLower(fileExt)

		// if the file extension exists in our list of extensions ..
		// for this current stuff, print it
		if Contains(s.Extensions, fileExt) {
			// the folder plus the folder in which we're moving the file
			destinationPath := fmt.Sprintf("%s/%s", folder, s.Folder)

			// the folder plus th filename in the looped dir
			fileName := fmt.Sprintf("%s/%s", folder, file.Name())

			// increment the counter for every moved file
			*noOfFilesPtr++

			MoveFile(destinationPath, fileName)
		}
	}

}

// Creates a folder to move the file or moves the file if folder exist
func MoveFile(folder, file string) {
	log.Printf("Moving '%s', to '%s' \n", file, folder)

	// create a destination directory if its not existing
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.Mkdir(folder, 0777)
		handleErr(err)
	}

	//finally, append the final path to include the filename
	destFolder := fmt.Sprintf("%s/%s", folder, filepath.Base(file))

	// move the file from the parent dir to the folder
	err := os.Rename(file, destFolder)
	handleErr(err)
}
