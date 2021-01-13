package main

import (
	"log"
	"github.com/fsnotify/fsnotify"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"fyne.io/fyne"
	"io/ioutil"
)

var listeningFile string
var watcher* fsnotify.Watcher
var label* widget.Label
var scroller* widget.ScrollContainer


func main() {
	SetupListener()
	SetupWindow()

	defer watcher.Close()
}


func UpdateFileDisplay() {
	dat, err := ioutil.ReadFile(listeningFile)

	if err != nil {
		label.SetText("Cannot read file.")
	} else {
		label.SetText(string(dat))
		scroller.ScrollToBottom()
	}

}

func ChangeListeningFile(filePath string) int {
	// returns 
	//	-1 if the previous file couldnt be removed
	//  0  if the file couldnt be found
	//  >0 if everything went ok
	log.Println(listeningFile)

	if listeningFile != "" {
		err := watcher.Remove(listeningFile)

		if err != nil {
			log.Println(err)
			return -1 
		}
	}
	err := watcher.Add(filePath)
	if err != nil {
		log.Println(err)
		listeningFile = ""
		return 0 
	}
	listeningFile = filePath
	return 1
}

func SetupWindow() {
	a := app.New()
	w := a.NewWindow("LiveFileObserver")

	label = widget.NewLabel("empty")
	input := widget.NewEntry()

	scroller = widget.NewScrollContainer(label)

	w.SetContent(widget.NewVBox(
		input,
		widget.NewButton("Look for file", func() {
			status := ChangeListeningFile(input.Text)

			switch (status) {
				case -1:
					label.SetText("An internal error occured! Restart the program and if the error persists, contact support.")
				case 0:
					label.SetText("File could not be found.")
				default:
					UpdateFileDisplay()
			}
		}),
		scroller,
	))

	scroller.SetMinSize(fyne.NewSize(1920, 500))
	w.ShowAndRun()
}


func SetupListener() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	
	go func() {
		for {
			select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						UpdateFileDisplay()
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
			}
		}
	}()
}
