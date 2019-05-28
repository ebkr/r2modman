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

	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	mainBox.SetBorderWidth(10)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	scroll.Add(scrollBox)

	grid, _ := gtk.GridNew()
	grid.Attach(scroll, 0, 0, 8, 200)

	removeMod, _ := gtk.ButtonNewWithLabel("Delete")
	updateMod, _ := gtk.ButtonNewWithLabel("Update")
	play, _ := gtk.ButtonNewWithLabel("Play RoR2")
	grid.Attach(removeMod, 0, 200, 2, 1)
	grid.Attach(updateMod, 2, 200, 2, 1)
	grid.Attach(play, 6, 200, 2, 1)

	grid.SetColumnSpacing(10)

	mainBox.PackStart(grid, true, true, 0)

	mods := modfetch.GetMods()
	mods = append(mods, modfetch.Mod{Name: "MyMod", URL: "http://..."})
	mods = append(mods, modfetch.Mod{Name: "MyMod2", URL: "http://..."})
	mods = append(mods, modfetch.Mod{Name: "MyMod3", URL: "http://..."})
	for _, mod := range mods {
		cb, _ := gtk.CheckButtonNewWithLabel(mod.Name)
		scrollBox.PackStart(cb, false, false, 2)
	}

	manager.window.Add(mainBox)
}
