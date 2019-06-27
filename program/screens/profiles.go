package screens

import (
	"fmt"
	"github.com/ebkr/r2modman/program/modfetch"
	"io/ioutil"
	"os"

	"github.com/ebkr/r2modman/program/globals"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// ProfileScreen : Screen to allow profile selection
type ProfileScreen struct {
	screen
}

// Show : Show the profiles screen
func (profile *ProfileScreen) Show(isCmdLine bool) {
	if profile.window == nil {
		profile.init("r2modman : Profiles")
		profile.create(isCmdLine)
		profile.window.SetPosition(gtk.WIN_POS_CENTER)
		profile.window.Connect("destroy", func() {
			if len(globals.SelectedProfile) == 0 {
				gtk.MainQuit()
			}
		})
	}
	profile.window.ShowAll()
}

func (profile *ProfileScreen) create(isCmdLine bool) {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	box.SetBorderWidth(10)
	profile.window.Add(box)

	listBox, _ := gtk.ListBoxNew()
	box.PackStart(listBox, true, true, 5)

	buttonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)

	deleteProfile, _ := gtk.ButtonNewWithLabel("Delete Profile")
	newProfile, _ := gtk.ButtonNewWithLabel("Create Profile")
	selectButton, _ := gtk.ButtonNewWithLabel("Select Profile")

	buttonBox.PackStart(deleteProfile, false, false, 2)
	buttonBox.PackStart(newProfile, false, false, 2)
	buttonBox.PackStart(selectButton, false, false, 2)

	box.PackEnd(buttonBox, false, false, 2)

	profile.updateListBox(listBox)

	_, _ = deleteProfile.Connect("clicked", func() {
		selected := listBox.GetSelectedRow()
		if selected == nil {
			return
		}
		selected.GetChildren().Foreach(func(child interface{}) {
			label := &gtk.Label{gtk.Widget{glib.InitiallyUnowned{child.(*gtk.Widget).Object}}}
			text, _ := label.GetText()
			_ = os.RemoveAll(globals.RootDirectory + "/mods/" + text + "/")
		})
		profile.updateListBox(listBox)
	})

	_, _ = selectButton.Connect("clicked", func() {
		selected := listBox.GetSelectedRow()
		if selected == nil {
			return
		}
		selected.GetChildren().Foreach(func(child interface{}) {
			// Convert *gtk.Widget to *gtk.Label.
			// Solution found through: https://github.com/gotk3/gotk3/issues/67#issuecomment-177609570
			label := &gtk.Label{gtk.Widget{glib.InitiallyUnowned{child.(*gtk.Widget).Object}}}
			globals.SelectedProfile, _ = label.GetText()
		})
		if len(globals.SelectedProfile) == 0 {
			return
		}
		if !isCmdLine {
			profile.showMainScreen()
		} else {
			thunderstoreProgression := make(chan float64)
			go modfetch.ThunderstoreGenerateList(thunderstoreProgression)
			<-thunderstoreProgression
			mod := modfetch.ThunderstoreInstallFromProtocol(globals.ROR2MMProtocol, profile.window)
			if mod != nil {

				mods := modfetch.GetMods()
				found := false
				for i, a := range mods {
					if a.FullName == mod.FullName {
						mods[i] = *mod
						found = true
						break
					}
				}
				if !found {
					mods = append(mods, *mod)
				}
				modfetch.UpdateMods(mods)
			}
			fmt.Println("Mod:", mod)
			gtk.MainQuit()
		}
	})

	_, _ = newProfile.Connect("clicked", func() {
		profile.showNewProfileDialog(false)
		profile.updateListBox(listBox)
	})

	profile.window.ShowAll()
}

// Update the list box
func (profile *ProfileScreen) updateListBox(listBox *gtk.ListBox) {
	listBox.GetChildren().Foreach(func(child interface{}) {
		listBox.Remove(child.(gtk.IWidget))
	})
	files, err := ioutil.ReadDir(globals.RootDirectory + "/mods/")
	if err != nil {
		return
	}
	for _, a := range files {
		if a.IsDir() {
			row, _ := gtk.ListBoxRowNew()
			rowText, _ := gtk.LabelNew(a.Name())
			rowText.SetXAlign(0)
			row.Add(rowText)
			row.SetBorderWidth(10)
			listBox.Add(row)
		}
	}
	listBox.ShowAll()
}

// Show the text entry dialog to create a new profile
func (profile *ProfileScreen) showNewProfileDialog(nameExists bool) {
	dialog, _ := gtk.DialogNew()
	box, _ := dialog.GetContentArea()
	box.SetBorderWidth(10)

	label, _ := gtk.LabelNew("Enter a profile name:")
	if nameExists {
		label.SetText("Profile already exists, try again:")
	}
	textField, _ := gtk.EntryNew()

	box.PackStart(label, false, false, 5)
	box.PackStart(textField, false, false, 5)

	doneButton, _ := dialog.AddButton("Done", gtk.RESPONSE_APPLY)
	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)

	dialog.SetDefault(doneButton)

	dialog.ShowAll()

	switch dialog.Run() {
	case gtk.RESPONSE_APPLY:
		text, err := textField.GetText()
		if err != nil || len(text) == 0 {
			fmt.Println(err.Error())
			dialog.Destroy()
			return
		}
		_, exists := os.Stat(globals.RootDirectory + "/mods/" + text)
		if os.IsNotExist(exists) {
			os.MkdirAll(globals.RootDirectory+"/mods/"+text, 0777)
			dialog.Destroy()
		} else {
			dialog.Destroy()
			profile.showNewProfileDialog(true)
		}
	case gtk.RESPONSE_CANCEL:
		dialog.Destroy()
	}

}

func (profile *ProfileScreen) showMainScreen() {
	profile.window.GetChildren().Foreach(func(child interface{}) {
		profile.window.Remove(child.(gtk.IWidget))
	})
	main := ManagerScreen{}
	main.Show()
	profile.window.Destroy()
}
