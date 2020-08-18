package main

import (
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	// Create Gtk Application, change appID to your application domain name reversed.
	const appID = "io.github.dbriemann"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	failOn(err)

	// Application signals available
	// startup -> sets up the application when it first starts
	// activate -> shows the default first window of the application (like a new document). This corresponds to the application being launched by the desktop environment.
	// open -> opens files and shows them in a new window. This corresponds to someone trying to open a document (or documents) using the application from the file browser, or similar.
	// shutdown ->  performs shutdown tasks
	// Setup Gtk Application callback signals
	application.Connect("activate", func() { onActivate(application) })
	// Run Gtk application
	os.Exit(application.Run(os.Args))
}

func makeMenu() *gtk.MenuBar {
	mainMenu, err := gtk.MenuBarNew()
	failOn(err)
	// File..
	fileMenu, err := gtk.MenuItemNewWithLabel("File")
	failOn(err)
	mainMenu.Add(fileMenu)
	// Edit..
	editMenu, err := gtk.MenuItemNewWithLabel("Edit")
	failOn(err)
	mainMenu.Add(editMenu)
	// Help
	helpMenu, err := gtk.MenuItemNewWithLabel("Help")
	failOn(err)
	mainMenu.Add(helpMenu)

	return mainMenu
}

func newWindow(application *gtk.Application) *gtk.ApplicationWindow {
	// Create ApplicationWindow
	win, err := gtk.ApplicationWindowNew(application)
	failOn(err)

	// Set ApplicationWindow Properties
	win.SetTitle("Seamstress")

	mainLayout, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	failOn(err)
	win.Add(mainLayout)

	menu := makeMenu()
	mainLayout.PackStart(menu, false, false, 0)

	win.SetDefaultSize(400, 600)

	return win
}

// Callback signal from Gtk Application
func onActivate(application *gtk.Application) {
	win := newWindow(application)

	win.ShowAll()
}

func failOn(err error) {
	if err != nil {
		panic(err)
	}
}
