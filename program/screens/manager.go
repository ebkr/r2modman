package screens

import (
	"github.com/ebkr/r2modman/program/modfetch"
	"github.com/gotk3/gotk3/gtk"
)

// ManagerScreen : The main screen
type ManagerScreen struct {
	screen
}

// Show : Show the manager screen
func (manager *ManagerScreen) Show() {
	if manager.window == nil {
		manager.init("r2modman")
		manager.create()
		manager.window.Connect("destroy", func() {
			gtk.MainQuit()
		})
	}
	manager.window.ShowAll()
}

func (manager *ManagerScreen) create() {

	manager.window.SetSizeRequest(200, 400)

	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	buttonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)

	scrollFrame, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)

	scrollFrame.Add(scrollBox)

	delete, _ := gtk.ButtonNewWithLabel("Delete")
	update, _ := gtk.ButtonNewWithLabel("Update")
	play, _ := gtk.ButtonNewWithLabel("Play Risk of Rain 2")

	buttonBox.PackStart(delete, false, false, 2)
	buttonBox.PackStart(update, false, false, 10)
	buttonBox.PackEnd(play, false, false, 0)

	mainBox.PackStart(scrollFrame, true, true, 2)
	mainBox.PackEnd(buttonBox, false, true, 0)
	manager.window.Add(mainBox)

	mods := modfetch.GetMods()
	mods = append(mods, modfetch.Mod{Name: "MyMod", URL: "http://..."})
	mods = append(mods, modfetch.Mod{Name: "MyMod2", URL: "http://..."})
	mods = append(mods, modfetch.Mod{Name: "MyMod3", URL: "http://..."})
	for _, mod := range mods {
		cb, _ := gtk.CheckButtonNewWithLabel(mod.Name)
		scrollBox.PackStart(cb, false, false, 2)
	}

	mainBox.SetBorderWidth(10)
}
