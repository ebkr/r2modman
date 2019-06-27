package screens

import (
	"fmt"
	"github.com/ebkr/r2modman/program/globals"
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

var searchInstalled *gtk.Entry
var searchAvailable *gtk.Entry
var globalListInstalled *gtk.ListBox

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

	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	mainBox.SetBorderWidth(10)

	register := globals.R2Registry{}
	registerProtocol, _ := gtk.InfoBarNew()
	registerLabel, _ := gtk.LabelNew("Associate ROR2MM:// links with r2modman?")
	contentArea, _ := registerProtocol.GetContentArea()
	contentArea.PackStart(registerLabel, false, false, 2)
	registerOk, _ := gtk.ButtonNewWithLabel("Ok")
	registerProtocol.AddActionWidget(registerOk, gtk.RESPONSE_ACCEPT)
	registerProtocol.SetShowCloseButton(true)
	if !register.IsAssociatedWithProtocol().Valid {
		mainBox.PackStart(registerProtocol, false, false, 2)
	}

	registerActivated := false
	_, _ = registerOk.Connect("clicked", func() {
		if registerActivated {
			mainBox.Remove(registerProtocol)
		}
		res := register.SetAssociatedProtocol()
		if !res.Valid {
			fmt.Println("Register Protocol Failed with Code:", res.Error.Code)
			fmt.Println("Reason:", res.Error.Reason)
			registerLabel.SetLabel("You need admin privileges to perform this action")
			registerActivated = true
			return
		} else {
			mainBox.Remove(registerProtocol)
		}
	})

	stack, _ := gtk.StackNew()
	stackSwitcher, _ := gtk.StackSwitcherNew()
	stackSwitcher.SetStack(stack)

	// Search bars
	searchInstalled, _ = gtk.EntryNew()
	searchAvailable, _ = gtk.EntryNew()

	searchInstalled.SetPlaceholderText("Search:")
	searchAvailable.SetPlaceholderText("Search:")

	// Main boxes
	boxInstalled, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	boxAvailable, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	stack.AddTitled(boxInstalled, "installed", "Installed Mods")
	stack.AddTitled(boxAvailable, "available", "Available Mods")

	boxInstalled.PackStart(searchInstalled, false, false, 2)
	boxAvailable.PackStart(searchAvailable, false, false, 2)

	// Scroll windows
	scrollWindowInstalled, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollWindowAvailable, _ := gtk.ScrolledWindowNew(nil, nil)
	boxInstalled.PackStart(scrollWindowInstalled, true, true, 2)
	boxAvailable.PackStart(scrollWindowAvailable, true, true, 2)

	// Lists
	listInstalled, _ := gtk.ListBoxNew()
	listAvailable, _ := gtk.ListBoxNew()
	scrollWindowInstalled.Add(listInstalled)
	scrollWindowAvailable.Add(listAvailable)

	// Reference global to allow update on download
	globalListInstalled = listInstalled

	// Buttons
	buttonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
	local, _ := gtk.ButtonNewFromIconName("list-add", gtk.ICON_SIZE_SMALL_TOOLBAR)
	play, _ := gtk.ButtonNewWithLabel("Play Risk of Rain 2")
	getUpdates, _ := gtk.ButtonNewFromIconName("view-refresh", gtk.ICON_SIZE_SMALL_TOOLBAR)

	buttonBox.PackStart(local, false, true, 2)
	buttonBox.PackStart(getUpdates, false, true, 2)
	buttonBox.PackStart(play, false, true, 2)

	_, _ = getUpdates.Connect("clicked", func() {
		manager.updateMods(listInstalled, "")
	})

	mainBox.PackEnd(buttonBox, false, true, 2)

	mainBox.PackStart(stackSwitcher, false, true, 0)
	mainBox.PackStart(stack, true, true, 0)
	manager.window.Add(mainBox)

	manager.updateMods(listInstalled, "")

	// Events
	_, _ = local.Connect("clicked", func() {
		fileChooser, _ := gtk.FileChooserNativeDialogNew("Select Mods", manager.window, gtk.FILE_CHOOSER_ACTION_OPEN, "Install Mod", "Cancel")
		fileFilter, _ := gtk.FileFilterNew()
		fileFilter.AddPattern("*.zip")
		fileChooser.SetFilter(fileFilter)
		fileChooser.Show()
		_, _ = fileChooser.Connect("response", func() {
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
				manager.updateMods(listInstalled, "")
			}
		})
	})

	_, exists := os.Open(globals.RootDirectory + "/program/path.txt")
	if os.IsNotExist(exists) {
		play.SetLabel("Locate Risk of Rain 2")
	}

	_, _ = play.Connect("clicked", func() {
		// Symlink
		file, exists := os.Open(globals.RootDirectory + "/program/path.txt")
		if os.IsNotExist(exists) {
			tempf, creationErr := os.Create(globals.RootDirectory + "/program/path.txt")
			if creationErr != nil {
				return
			}
			file = tempf
			selector, _ := gtk.FileChooserNativeDialogNew("Select Risk of Rain 2 Location", manager.window, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER, "Select Folder", "Cancel")
			selector.Show()
			_, _ = selector.Connect("response", func() {
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
				os.Remove(globals.RootDirectory + "/program/path.txt")
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

	// Search bar events
	_, _ = searchInstalled.Connect("notify::text", func() {
		scrollWindowInstalled.Remove(listInstalled)
		newInstalled, _ := gtk.ListBoxNew()
		listInstalled = newInstalled
		globalListInstalled = newInstalled
		scrollWindowInstalled.Add(newInstalled)
		filter, _ := searchInstalled.GetText()
		manager.updateMods(newInstalled, strings.ToLower(filter))
	})
	_, _ = searchAvailable.Connect("notify::text", func() {
		scrollWindowAvailable.Remove(listAvailable)
		newAvailable, _ := gtk.ListBoxNew()
		listAvailable = newAvailable
		scrollWindowAvailable.Add(newAvailable)
		filter, _ := searchAvailable.GetText()
		manager.downloadThunderstoreList(newAvailable, strings.ToLower(filter))
	})

	go manager.downloadThunderstoreList(listAvailable, "")

}

func (manager *ManagerScreen) updateMods(listBox *gtk.ListBox, filter string) {
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
			if !mod.Enabled {
				icon.SetSensitive(false)
			}
		}

		name, _ := gtk.LabelNew(mod.Name)
		rowBox.PackStart(name, false, false, 2)

		//delete, _ := gtk.ButtonNewFromIconName("window-close", gtk.ICON_SIZE_SMALL_TOOLBAR)
		settings, _ := gtk.ButtonNewFromIconName("emblem-system-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
		settings.SetProperty("relief", gtk.RELIEF_NONE)
		rowBox.PackEnd(settings, false, false, 2)

		if !mod.Enabled {
			name.SetSensitive(false)
		}

		// Thunderstore Update Integration
		if strings.Compare(mod.Uuid4, "") == 0 {
			syncThunderstoreAlert, _ := gtk.ButtonNewFromIconName("dialog-warning", gtk.ICON_SIZE_SMALL_TOOLBAR)
			func(m modfetch.Mod, button *gtk.Button) {
				button.Connect("clicked", func() {
					refreshedMods := modfetch.GetMods()
					modfetch.ThunderstoreLocalToOnline(&m)
					for i, a := range refreshedMods {
						if strings.Compare(a.FullName, m.FullName) == 0 && strings.Compare(a.Description, m.Description) == 0 {
							refreshedMods[i] = m
							break
						}
					}
					modfetch.UpdateMods(refreshedMods)
					manager.updateMods(listBox, filter)
				})
			}(mod, syncThunderstoreAlert)
			rowBox.PackEnd(syncThunderstoreAlert, false, false, 2)
		} else {
			func(m modfetch.Mod) {
				if modfetch.ThunderstoreModHasUpdate(&m) {
					fmt.Println(mod.Name, "has update")
					updateAvailable, _ := gtk.ButtonNewFromIconName("software-update-urgent", gtk.ICON_SIZE_SMALL_TOOLBAR)
					updateAvailable.Connect("clicked", func() {
						updatedMod := modfetch.ThunderstoreUpdateMod(&m, manager.window)
						refreshedMods := modfetch.GetMods()
						for i, a := range refreshedMods {
							if strings.Compare(a.Uuid4, updatedMod.Uuid4) == 0 {
								refreshedMods[i] = *updatedMod
								break
							}
						}
						modfetch.UpdateMods(refreshedMods)
						manager.updateMods(listBox, filter)
					})
					rowBox.PackEnd(updateAvailable, false, false, 2)
				}
			}(mod)

		}

		// Dependency Icon
		for _, a := range mod.Dependencies {
			if !mod.DependencyExists(&a) {
				missingDependency, _ := gtk.ButtonNewFromIconName("sync-error", gtk.ICON_SIZE_SMALL_TOOLBAR)
				rowBox.PackEnd(missingDependency, false, false, 2)
				func(mod modfetch.Mod, dep modfetch.ModDependency) {
					missingDependency.Connect("clicked", func() {
						// Missing Dependency Dialog
						updatedMod := modfetch.ThunderstoreGetDependency(a.Name, manager.window)
						if updatedMod != nil {
							updatedMod.Enabled = true
							refreshedMods := modfetch.GetMods()
							found := false
							for i, reMod := range refreshedMods {
								if strings.Compare(reMod.Uuid4, updatedMod.Uuid4) == 0 {
									refreshedMods[i] = *updatedMod
									found = true
									break
								}
							}
							if !found {
								refreshedMods = append(refreshedMods, *updatedMod)
							}
							modfetch.UpdateMods(refreshedMods)
							manager.updateMods(listBox, filter)
						}
					})
				}(mod, a)
				break
			}
		}

		// Events
		func(m modfetch.Mod) {
			settings.Connect("clicked", func() {
				manager.showSettingsDialog(&m)
				refreshedMods := modfetch.GetMods()
				for i, a := range refreshedMods {
					if strings.Compare(a.FullName, m.FullName) == 0 && strings.Compare(a.Description, m.Description) == 0 {
						refreshedMods[i] = m
						break
					}
				}
				modfetch.UpdateMods(refreshedMods)
				manager.updateMods(listBox, filter)
			})
		}(mod)

		if strings.Contains(strings.ToLower(mod.Name), filter) {
			listBox.Add(row)
		} else {
			row.Destroy()
		}

		// End of loop
	}
	manager.window.ShowAll()
}

func (manager *ManagerScreen) downloadThunderstoreList(listBox *gtk.ListBox, filter string) {
	listBox.GetChildren().Foreach(func(child interface{}) {
		listBox.Remove(child.(gtk.IWidget))
	})
	for _, mod := range modfetch.ThunderstoreGetAll() {
		row, _ := gtk.ListBoxRowNew()
		box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
		row.Add(box)

		image, _ := gtk.ImageNew()
		box.PackStart(image, false, false, 5)
		if len(mod.Versions[0].Icon) > 0 {
			image.SetFromPixbuf(modfetch.ThunderstoreGetPixbufFromUUID4(mod.Uuid4))
		}

		label, _ := gtk.LabelNew(mod.Name)
		box.PackStart(label, false, false, 2)

		download, _ := gtk.ButtonNewFromIconName("go-down", gtk.ICON_SIZE_SMALL_TOOLBAR)
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
				text, _ := searchInstalled.GetText()
				manager.updateMods(globalListInstalled, strings.ToLower(text))
			})
		}(mod)

		if strings.Contains(strings.ToLower(mod.Name), filter) {
			listBox.Add(row)
		} else {
			row.Destroy()
		}
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

	dialog.AddButton("Uninstall", gtk.RESPONSE_CANCEL)

	if mod.Enabled {
		dialog.AddButton("Disable", gtk.RESPONSE_APPLY)
	} else {
		dialog.AddButton("Enable", gtk.RESPONSE_APPLY)
	}

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
	case gtk.RESPONSE_APPLY:
		mod.Enabled = !mod.Enabled
		dialog.Destroy()
		return
	default:
		dialog.Destroy()
		return
	}
}
