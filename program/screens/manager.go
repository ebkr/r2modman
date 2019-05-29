package screens

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"

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

	manager.window.SetDefaultSize(400, 250)

	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	mainBox.SetBorderWidth(10)

	stack, _ := gtk.StackNew()
	stackSwitcher, _ := gtk.StackSwitcherNew()
	stackSwitcher.SetStack(stack)

	boxInstalled, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	boxAvailable, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	stack.AddTitled(boxInstalled, "installed", "Installed Mods")
	stack.AddTitled(boxAvailable, "available", "Available Mods")

	scrollWindowInstalled, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollWindowAvailable, _ := gtk.ScrolledWindowNew(nil, nil)
	boxInstalled.PackStart(scrollWindowInstalled, true, true, 2)
	boxAvailable.PackStart(scrollWindowAvailable, true, true, 2)

	listInstalled, _ := gtk.ListBoxNew()
	listAvailable, _ := gtk.ListBoxNew()
	scrollWindowInstalled.Add(listInstalled)
	scrollWindowAvailable.Add(listAvailable)

	buttonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
	local, _ := gtk.ButtonNewWithLabel("Install Local Zip")
	play, _ := gtk.ButtonNewWithLabel("Play Risk of Rain 2")

	buttonBox.PackStart(local, false, true, 2)
	buttonBox.PackStart(play, false, true, 2)

	mainBox.PackEnd(buttonBox, false, true, 2)

	mainBox.PackStart(stackSwitcher, false, true, 0)
	mainBox.PackStart(stack, true, true, 0)
	manager.window.Add(mainBox)

	manager.updateMods(listInstalled)
	//go manager.downloadThunderstoreList()

	// Events
	local.Connect("clicked", func() {
		fileChooser, _ := gtk.FileChooserNativeDialogNew("Select Mods", manager.window, gtk.FILE_CHOOSER_ACTION_OPEN, "Install Mod", "Cancel")
		fileFilter, _ := gtk.FileFilterNew()
		fileFilter.AddPattern("*.zip")
		fileChooser.SetFilter(fileFilter)
		fileChooser.Show()
		fileChooser.Connect("response", func() {
			filePath := strings.Split(fileChooser.GetFilename(), "\\")
			fileNameSplit := strings.Split(filePath[len(filePath)-1], ".")
			fileName := strings.Join(fileNameSplit[0:len(fileNameSplit)-1], ".")
			res := modfetch.Unzip(fileName, fileChooser.GetFilename())
			val, exists := res["manifest.json"]
			if exists {
				mod := modfetch.MakeModFromManifest(val, "")
				mods := modfetch.GetMods()
				mods = append(mods, mod)
				modfetch.UpdateMods(mods)
				manager.updateMods(listInstalled)
			}
		})
	})

	go manager.downloadThunderstoreList(listAvailable)

}

func (manager *ManagerScreen) updateMods(listBox *gtk.ListBox) {
	listBox.GetChildren().Foreach(func(child interface{}) {
		listBox.Remove(child.(gtk.IWidget))
	})
	mods := modfetch.GetMods()
	for _, mod := range mods {
		row, _ := gtk.ListBoxRowNew()
		rowBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
		row.Add(rowBox)
		_, err := os.Stat(mod.Path + "/icon.png")
		if !os.IsNotExist(err) {
			// Load image from Pixbuf
			pbloader, _ := gdk.PixbufLoaderNew()
			pbloader.SetSize(32, 32)
			// Open file
			file, _ := os.Open(mod.Path + "/icon.png")
			fileData, _ := ioutil.ReadAll(file)
			file.Close()
			// Initialise PixbufLoader from file data
			pbloader.Write(fileData)
			pix, _ := pbloader.GetPixbuf()
			// Create image
			icon, _ := gtk.ImageNewFromPixbuf(pix)
			rowBox.PackStart(icon, false, false, 5)
		}

		name, _ := gtk.LabelNew(mod.Name)
		rowBox.PackStart(name, false, false, 2)

		delete, _ := gtk.ButtonNewFromIconName("window-close-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
		delete.SetProperty("relief", gtk.RELIEF_NONE)
		rowBox.PackEnd(delete, false, false, 2)
		listBox.Add(row)

		// Events
		delete.Connect("clicked", func() {
			refreshedMods := modfetch.GetMods()
			index := -1
			for i, a := range refreshedMods {
				if strings.Compare(mod.Name, a.Name) == 0 {
					err := os.RemoveAll(mod.Path)
					if err != nil {
						fmt.Println(err.Error())
					}
					index = i
				}
			}
			if index >= 0 {
				refreshedMods = append(refreshedMods[:index], refreshedMods[index+1:]...)
				modfetch.UpdateMods(refreshedMods)
				manager.updateMods(listBox)
			}
		})

		// End of loop
	}
	manager.window.ShowAll()
}

func (manager *ManagerScreen) downloadThunderstoreList(listBox *gtk.ListBox) {
	modfetch.ThunderstoreGenerateList()
	listBox.GetChildren().Foreach(func(child interface{}) {
		listBox.Remove(child.(gtk.IWidget))
	})
	for _, mod := range modfetch.ThunderstoreGetAll() {
		row, _ := gtk.ListBoxRowNew()
		box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
		row.Add(box)
		listBox.Add(row)

		image, _ := gtk.ImageNew()
		box.PackStart(image, false, false, 5)
		if len(mod.Versions[0].Icon) > 0 {
			image.SetFromPixbuf(pixbufFromURL(mod.Versions[0].Icon))
		}

		label, _ := gtk.LabelNew(mod.Name)
		box.PackStart(label, false, false, 2)

		download, _ := gtk.ButtonNewFromIconName("emblem-downloads-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
		box.PackEnd(download, false, false, 2)
	}
	listBox.ShowAll()
}

func pixbufFromURL(url string) *gdk.Pixbuf {
	pbloader, _ := gdk.PixbufLoaderNew()
	pbloader.SetSize(32, 32)
	stream, _ := http.Get(url)
	r, _ := ioutil.ReadAll(stream.Body)
	pbloader.Write(r)
	stream.Body.Close()
	pix, _ := pbloader.GetPixbuf()
	return pix
}
