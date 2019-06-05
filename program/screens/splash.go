package screens

import (
	"time"

	"github.com/ebkr/r2modman/program/modfetch"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type SplashScreen struct {
	screen
	completed bool
}

var progressBar *gtk.ProgressBar

func (splash *SplashScreen) Show() {
	if splash.window == nil {
		splash.init("r2modman")
		splash.create()
		splash.window.SetPosition(gtk.WIN_POS_CENTER)
		splash.window.Connect("destroy", func() {
			if !splash.completed {
				gtk.MainQuit()
			}
		})
		splash.window.SetDecorated(false)
	}
	splash.window.ShowAll()

	thunderstoreProgression := make(chan float64)
	receivedProgression := false

	var drawHandle glib.SignalHandle
	drawHandle, _ = splash.window.Connect("draw", func() {
		splash.window.HandlerDisconnect(drawHandle)
		go modfetch.ThunderstoreGenerateList(thunderstoreProgression)
	})
	var id int
	id = splash.window.AddTickCallback(func(w *gtk.Widget, frameClock *gdk.FrameClock, userData uintptr) bool {
		if !receivedProgression {
			select {
			case <-thunderstoreProgression:
				receivedProgression = true
				break
			case <-time.After(0):
				progressBar.Pulse()
			}
		} else if modfetch.ThunderstoreReady() {
			splash.window.RemoveTickCallback(id)
			splash.completed = true
			close(thunderstoreProgression)
			splash.showProfileScreen()
		} else {
			select {
			case frac := <-thunderstoreProgression:
				progressBar.SetFraction(frac)
				break
			case <-time.After(0):
				break
			}
		}
		return true
	}, 1000)
}

func (splash *SplashScreen) create() {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)

	image, _ := gtk.ImageNewFromFile("./program/assets/r2modman.png")
	progressBar, _ = gtk.ProgressBarNew()

	box.PackStart(image, false, false, 0)
	box.PackStart(progressBar, false, false, 0)
	splash.window.Add(box)
}

func (splash *SplashScreen) showProfileScreen() {
	splash.window.Destroy()
	profiles := ProfileScreen{}
	profiles.Show()
}
