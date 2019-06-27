package dialogs

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// AwaitDialog : Used to show that a task is in process
type AwaitDialog struct {
	Title       string
	Description string
}

// NewAwaitDialog : Creates a new dialog for awaiting a task
func NewAwaitDialog(listener chan *gtk.Dialog, title, description string) {
	ad := AwaitDialog{Title: title, Description: description}
	ad.Setup(listener)

}

// Setup : Handles the creation and timer of the dialog
func (await *AwaitDialog) Setup(channel chan *gtk.Dialog) {

	dialog, _ := gtk.DialogNew()
	dialog.SetTitle(await.Title)
	box, _ := dialog.GetContentArea()
	box.SetBorderWidth(10)
	dialog.SetPosition(gtk.WIN_POS_CENTER)

	action, _ := gtk.LabelNew(await.Description)
	pulse, _ := gtk.ProgressBarNew()

	box.PackStart(action, false, false, 5)
	box.PackStart(pulse, false, false, 5)

	dialog.ShowAll()

	dialog.AddTickCallback(func(w *gtk.Widget, frameClock *gdk.FrameClock, userData uintptr) bool {
		pulse.Pulse()
		return true
	}, 10000)

	channel <- dialog
	dialog.Run()
}
