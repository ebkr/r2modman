package screens

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"

	"github.com/ebkr/r2modman/program/modfetch"
	"github.com/ebkr/r2modman/program/modinstall"
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
		manager.window.SetKeepAbove(false)
		manager.window.SetPosition(gtk.WIN_POS_CENTER)
		manager.window.Connect("destroy", func() {
			gtk.MainQuit()
		})
	}
	manager.window.ShowAll()
}

func (manager *ManagerScreen) create() {

	manager.window.SetDefaultSize(400, 300)

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
	local, _ := gtk.ButtonNewFromIconName("list-add-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	play, _ := gtk.ButtonNewWithLabel("Play Risk of Rain 2")
	getUpdates, _ := gtk.ButtonNewFromIconName("view-refresh-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)

	buttonBox.PackStart(local, false, true, 2)
	buttonBox.PackStart(getUpdates, false, true, 2)
	buttonBox.PackStart(play, false, true, 2)

	getUpdates.Connect("clicked", func() {
		manager.updateMods(listInstalled)
	})

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

	_, exists := os.Open("./program/path.txt")
	if os.IsNotExist(exists) {
		play.SetLabel("Locate Risk of Rain 2")
	}

	play.Connect("clicked", func() {
		// Symlink
		file, exists := os.Open("./program/path.txt")
		if os.IsNotExist(exists) {
			tempf, creationErr := os.Create("./program/path.txt")
			if creationErr != nil {
				return
			}
			file = tempf
			selector, _ := gtk.FileChooserNativeDialogNew("Select Risk of Rain 2 Location", manager.window, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER, "Select Folder", "Cancel")
			selector.Show()
			selector.Connect("response", func() {
				folder := selector.GetFilename()
				fmt.Println("Folder:", folder)
				file.Write([]byte(folder))
				file.Close()
				play.SetLabel("Play Risk of Rain 2")
			})
		} else {
			data, _ := ioutil.ReadAll(file)
			defer file.Close()
			gamePath := string(data)
			fmt.Println(gamePath)
			_, searchErr := exec.LookPath(gamePath + "/Risk of Rain 2.exe")

			if searchErr != nil {
				fmt.Println("Path Error:", searchErr.Error())
				os.Remove("./program/path.txt")
			}
			pluginPath := gamePath + "/BepInEx/plugins/"
			os.MkdirAll(pluginPath, 0777)

			mods := modfetch.GetMods()
			prepareErr := modinstall.PrepareInstall()
			if prepareErr != nil {
				fmt.Println(prepareErr.Error())
				return
			}
			for _, mod := range mods {
				modinstall.InstallMod(&mod, gamePath+"/BepInEx/")
			}
			go func() {
				time.Sleep(10)
				cmd := exec.Command("powershell", "-c", "start steam://run/632360")
				cmd.Run()
			}()
		}
	})

	go manager.downloadThunderstoreList(listAvailable, listInstalled)

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

		//delete, _ := gtk.ButtonNewFromIconName("window-close-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
		delete, _ := gtk.ButtonNewFromIconName("emblem-system-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
		delete.SetProperty("relief", gtk.RELIEF_NONE)
		rowBox.PackEnd(delete, false, false, 2)

		// Thunderstore Update Integration
		if strings.Compare(mod.Uuid4, "") == 0 {
			syncThunderstoreAlert, _ := gtk.ButtonNewFromIconName("dialog-warning-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
			func(m modfetch.Mod, button *gtk.Button) {
				button.Connect("clicked", func() {
					refreshedMods := modfetch.GetMods()
					modfetch.ThunderstoreLocalToOnline(&m)
					for i, a := range refreshedMods {
						if strings.Compare(a.Name, m.Name) == 0 && strings.Compare(a.Description, m.Description) == 0 {
							refreshedMods[i] = m
							break
						}
					}
					modfetch.UpdateMods(refreshedMods)
					manager.updateMods(listBox)
				})
			}(mod, syncThunderstoreAlert)
			rowBox.PackEnd(syncThunderstoreAlert, false, false, 2)
		} else {
			func(m modfetch.Mod) {
				if modfetch.ThunderstoreModHasUpdate(&m) {
					updateAvailable, _ := gtk.ButtonNewFromIconName("software-update-urgent-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
					updateAvailable.Connect("clicked", func() {
						updatedMod := modfetch.ThunderstoreUpdateMod(&m, manager.window)
						refreshedMods := modfetch.GetMods()
						modfetch.ThunderstoreLocalToOnline(&m)
						for i, a := range refreshedMods {
							if strings.Compare(a.Uuid4, updatedMod.Uuid4) == 0 {
								refreshedMods[i] = *updatedMod
								break
							}
						}
						modfetch.UpdateMods(refreshedMods)
						manager.updateMods(listBox)
					})
					rowBox.PackEnd(updateAvailable, false, false, 2)
				}
			}(mod)

		}

		// Dependency Icon
		for _, a := range mod.Dependencies {
			if !mod.DependencyExists(&a) {
				missingDependency, _ := gtk.ButtonNewFromIconName("sync-error-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
				rowBox.PackEnd(missingDependency, false, false, 2)
				fmt.Println("Missing mod:", a.Name)
				func(mod modfetch.Mod, dep modfetch.ModDependency) {
					missingDependency.Connect("clicked", func() {
						// Missing Dependency Dialog
						updatedMod := modfetch.ThunderstoreGetDependency(a.Name, manager.window)
						if updatedMod != nil {
							refreshedMods := modfetch.GetMods()
							refreshedMods = append(refreshedMods, *updatedMod)
							modfetch.UpdateMods(refreshedMods)
							manager.updateMods(listBox)
						}
					})
				}(mod, a)
				break
			}
		}

		// Events
		func(m modfetch.Mod) {
			delete.Connect("clicked", func() {
				manager.showSettingsDialog(&m)
				manager.updateMods(listBox)
			})
		}(mod)

		listBox.Add(row)

		// End of loop
	}
	manager.window.ShowAll()
}

func (manager *ManagerScreen) downloadThunderstoreList(listBox *gtk.ListBox, installedBox *gtk.ListBox) {
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
			image.SetFromPixbuf(modfetch.ThunderstoreGetPixbufFromUUID4(mod.Uuid4))
		}

		label, _ := gtk.LabelNew(mod.Name)
		box.PackStart(label, false, false, 2)

		download, _ := gtk.ButtonNewFromIconName("go-down-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
		download.SetProperty("relief", gtk.RELIEF_NONE)
		box.PackEnd(download, false, false, 2)

		func(mod modfetch.ThunderstoreJSON) {
			download.Connect("clicked", func() {
				newMod := modfetch.ThunderstoreDownloadMod(mod.Uuid4, manager.window)
				if newMod == nil {
					return
				}
				newMod.Uuid4 = mod.Uuid4
				refreshedMods := modfetch.GetMods()
				refreshedMods = append(refreshedMods, *newMod)
				modfetch.UpdateMods(refreshedMods)
				manager.updateMods(installedBox)
			})
		}(mod)
	}
	listBox.ShowAll()
}

func (manager *ManagerScreen) showSettingsDialog(mod *modfetch.Mod) {

	// Init
	dialog, _ := gtk.DialogNew()
	box, _ := dialog.GetContentArea()
	box.SetBorderWidth(10)

	// Display Mod Information
	name, _ := gtk.LabelNew("Mod: " + mod.Name)
	version, _ := gtk.LabelNew("Version: " + mod.Version.String())
	box.PackStart(name, false, false, 5)
	box.PackStart(version, false, false, 5)

	if len(mod.Uuid4) > 0 {
		thunder := modfetch.ThunderstoreGetModByUUID4(mod.Uuid4)
		if thunder != nil {
			owner, _ := gtk.LabelNew("Author: " + thunder.Owner)
			url, _ := gtk.LinkButtonNewWithLabel(thunder.Package_url, "View On Thunderstore")
			box.PackStart(owner, false, false, 5)
			box.PackStart(url, false, false, 5)
		}
	}

	dialog.AddButton("Delete", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Close", gtk.RESPONSE_CLOSE)

	dialog.SetTitle(mod.Name + " Settings")

	dialog.SetPosition(gtk.WIN_POS_CENTER)

	dialog.ShowAll()

	switch dialog.Run() {
	case gtk.RESPONSE_CLOSE:
		dialog.Destroy()
		return
	case gtk.RESPONSE_CANCEL:
		modfetch.RemoveMod(mod)
		dialog.Destroy()
		return
	default:
		dialog.Destroy()
		return
	}
}
