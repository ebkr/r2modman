package screens

import (
	"time"

	"github.com/gotk3/gotk3/gtk"
)

// SplashScreen : Screen featuring a splash image
type SplashScreen struct {
	screen
	progBar *gtk.ProgressBar
}

func NewSplashScreen() *SplashScreen {
	splash := SplashScreen{}
	splash.init("Splash")
	splash.create()
	splash.window.SetDecorated(false)
	splash.window.SetPosition(gtk.WIN_POS_CENTER)
	splash.window.Connect("destroy", func() {
		gtk.MainQuit()
	})
	return &splash
}

// Show : Show screen
func (splash *SplashScreen) Show(completed chan bool) {
	splash.window.ShowAll()
	go splash.waitAndProgress(completed)
}

func (splash *SplashScreen) create() {

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	box.SetBorderWidth(2)

	image, _ := gtk.ImageNewFromFile("./program/assets/r2modman.png")
	progBar, _ := gtk.ProgressBarNew()
	splash.progBar = progBar

	// Add to interface
	box.PackStart(image, false, false, 0)
	box.PackEnd(progBar, true, true, 0)
	splash.window.Add(box)
}

func (splash *SplashScreen) waitAndProgress(completed chan bool) {
	frac := 0.0
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second / 100)
		splash.progBar.SetFraction(frac)
		frac += 0.01
	}
	time.Sleep(100)
	splash.progBar.SetFraction(1)
	time.Sleep(time.Second)
	splash.window.Destroy()
	completed <- true
}
