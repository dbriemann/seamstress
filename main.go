package main

import (
	"fmt"
	"math"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var (
	mainWin      *gtk.Window
	imgWin       *gtk.Window
	img          *gtk.Image
	screenWidth  int
	screenHeight int
)

func closeImg() {
	imgWin.Hide()
}

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	failOn(err)
	mainWin = win
	mainWin.SetTitle("Seamstress")
	mainWin.Connect("destroy", func() {
		gtk.MainQuit()
	})

	display, err := gdk.DisplayGetDefault()
	failOn(err)
	mon, err := display.GetPrimaryMonitor()
	failOn(err)
	screenWidth = mon.GetWorkarea().GetWidth()
	screenHeight = mon.GetWorkarea().GetHeight()

	win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	failOn(err)
	imgWin = win
	imgWin.SetTitle("Image Window")
	imgWin.Hide()
	imgWin.Connect("destroy", closeImg)
	// TODO: hide or undeletable

	img, err = gtk.ImageNew()
	failOn(err)
	imgWin.Add(img)
	imgWin.SetResizable(false)

	mainLayout, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	failOn(err)
	mainWin.Add(mainLayout)

	menu := makeMainMenu()
	mainLayout.PackStart(menu, false, false, 0)

	content := makeMainContent()
	mainLayout.PackStart(content, true, true, 0)

	mainWin.SetDefaultSize(400, 600)
	mainWin.ShowAll()
	mainWin.SetKeepAbove(true)
	gtk.Main()
}

type itemAction struct {
	name   string
	action func()
}

func makeMenu(root string, children ...itemAction) *gtk.MenuItem {
	rootItem, err := gtk.MenuItemNewWithLabel(root)
	failOn(err)
	menu, err := gtk.MenuNew()
	failOn(err)
	rootItem.SetSubmenu(menu)

	for _, c := range children {
		if c.name != "" {
			item, err := gtk.MenuItemNewWithLabel(c.name)
			failOn(err)
			item.Connect("activate", c.action)
			menu.Append(item)
		} else {
			item, err := gtk.SeparatorMenuItemNew()
			failOn(err)
			menu.Append(item)
		}
	}
	return rootItem
}

func nilfunc() {}

func makeMainMenu() *gtk.MenuBar {
	mainMenu, err := gtk.MenuBarNew()
	failOn(err)

	// File..
	fileMenu := makeMenu("File", itemAction{"Open", openImage}, itemAction{"Save", nilfunc}, itemAction{"", nilfunc}, itemAction{"Quit", gtk.MainQuit})
	mainMenu.Append(fileMenu)
	// Edit..
	editMenu := makeMenu("Edit", itemAction{"Mark", nilfunc}, itemAction{"Cut", nilfunc})
	mainMenu.Append(editMenu)
	// Help..
	helpMenu := makeMenu("Help", itemAction{"About", nilfunc})
	mainMenu.Append(helpMenu)

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

func carve(twidth int, theight int) {
	fmt.Println("carve", twidth-img.GetPixbuf().GetWidth(), theight-img.GetPixbuf().GetHeight())
}

func openImage() {
	dialog, err := gtk.FileChooserNativeDialogNew("Open Image", mainWin, gtk.FILE_CHOOSER_ACTION_OPEN, "_Open", "_Cancel")
	failOn(err)
	filter, err := gtk.FileFilterNew()
	failOn(err)
	filter.AddPattern("*.png")
	filter.AddPattern("*.jpg")
	filter.AddPattern("*.jpeg")
	dialog.AddFilter(filter)
	resp := dialog.Run()
	if gtk.ResponseType(resp) == gtk.RESPONSE_ACCEPT {
		chooser := dialog
		fname := chooser.GetFilename()
		fmt.Println(fname)
		// TODO - fail nicer
		pix, err := gdk.PixbufNewFromFile(fname)
		failOn(err)
		imgWidth := pix.GetWidth()
		imgHeight := pix.GetHeight()

		var wdiv, hdiv int
		if imgWidth > screenWidth {
			wdiv = int(math.Ceil(float64(imgWidth) / float64(screenWidth)))
		}
		if imgHeight > screenHeight {
			hdiv = int(math.Ceil(float64(imgHeight) / float64(screenHeight)))
		}

		if wdiv != 0 && wdiv > hdiv {
			pix, err = pix.ScaleSimple(imgWidth/wdiv, imgHeight/wdiv, gdk.INTERP_BILINEAR)
			failOn(err)
		} else if hdiv != 0 {
			pix, err = pix.ScaleSimple(imgWidth/hdiv, imgHeight/hdiv, gdk.INTERP_BILINEAR)
			failOn(err)
		}
		img.SetFromPixbuf(pix)

		// TODO : look into later -> cairo.CreateImageSurfaceForData()
		imgWin.ShowAll()
	}
}

func makeMainContent() *gtk.Grid {
	grid, err := gtk.GridNew()
	failOn(err)

	openBtn := makeMainButton("gtk-open")
	openBtn.Connect("clicked", openImage)
	saveBtn := makeMainButton("gtk-save")
	markBtn := makeMainButton("gtk-edit")
	cutBtn := makeMainButton("gtk-cut")
	undoBtn := makeMainButton("gtk-undo")
	redoBtn := makeMainButton("gtk-redo")
	plusYBtn := makeMainButton("gtk-add")
	plusXBtn := makeMainButton("gtk-add")
	minusYBtn := makeMainButton("gtk-remove")
	minusXBtn := makeMainButton("gtk-remove")

	grid.Attach(openBtn, 1, 1, 1, 1)
	grid.Attach(saveBtn, 2, 1, 1, 1)
	grid.Attach(markBtn, 1, 2, 1, 1)
	grid.Attach(cutBtn, 2, 2, 1, 1)
	grid.Attach(undoBtn, 1, 3, 1, 1)
	grid.Attach(redoBtn, 2, 3, 1, 1)
	l, err := gtk.LabelNew("")
	failOn(err)
	grid.Attach(l, 1, 4, 2, 1)

	scaleGrid, err := gtk.GridNew()
	failOn(err)
	scaleGrid.Attach(plusYBtn, 2, 1, 1, 1)
	scaleGrid.Attach(minusXBtn, 1, 2, 1, 1)
	scaleGrid.Attach(plusXBtn, 3, 2, 1, 1)
	scaleGrid.Attach(minusYBtn, 2, 3, 1, 1)

	sbtn, err := gtk.SpinButtonNewWithRange(1.0, 100.0, 1.0)
	failOn(err)
	scaleGrid.Attach(sbtn, 2, 2, 1, 1)

	grid.Attach(scaleGrid, 1, 5, 2, 2)

	return grid
}

func failOn(err error) {
	if err != nil {
		panic(err)
	}
}
