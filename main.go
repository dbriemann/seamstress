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

func makeMenu(root string, children ...string) *gtk.MenuItem {
	rootItem, err := gtk.MenuItemNewWithLabel(root)
	failOn(err)
	menu, err := gtk.MenuNew()
	failOn(err)
	rootItem.SetSubmenu(menu)

	for _, c := range children {
		if c != "" {
			item, err := gtk.MenuItemNewWithLabel(c)
			failOn(err)
			menu.Add(item)
		} else {
			item, err := gtk.SeparatorMenuItemNew()
			failOn(err)
			menu.Add(item)
		}
	}
	return rootItem
}

func makeMainMenu() *gtk.MenuBar {
	mainMenu, err := gtk.MenuBarNew()
	failOn(err)

	// File..
	fileMenu := makeMenu("File", "Open", "Save", "", "Quit")
	mainMenu.Add(fileMenu)
	// Edit..
	editMenu := makeMenu("Edit", "Mark", "Cut")
	mainMenu.Add(editMenu)
	// Help..
	helpMenu := makeMenu("Help", "About")
	mainMenu.Add(helpMenu)

	return mainMenu
}

func makeMainButton(iconName string) *gtk.Button {
	btn, err := gtk.ButtonNew()
	failOn(err)
	img, err := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_DIALOG)
	failOn(err)
	btn.Add(img)
	btn.SetHExpand(true)
	btn.SetVExpand(true)

	return btn
}

func makeMainContent() *gtk.Grid {
	grid, err := gtk.GridNew()
	failOn(err)

	openBtn := makeMainButton("gtk-open")
	saveBtn := makeMainButton("gtk-save")
	markBtn := makeMainButton("gtk-edit")
	cutBtn := makeMainButton("gtk-cut")
	undoBtn := makeMainButton("gtk-undo")
	redoBtn := makeMainButton("gtk-redo")

	grid.Attach(openBtn, 1, 1, 1, 1)
	grid.Attach(saveBtn, 2, 1, 1, 1)
	grid.Attach(markBtn, 1, 2, 1, 1)
	grid.Attach(cutBtn, 2, 2, 1, 1)
	grid.Attach(undoBtn, 1, 3, 1, 1)
	grid.Attach(redoBtn, 2, 3, 1, 1)
	return grid
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

	menu := makeMainMenu()
	mainLayout.PackStart(menu, false, false, 0)

	content := makeMainContent()
	mainLayout.PackStart(content, true, true, 0)

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
