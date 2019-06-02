package screens

import (
	"github.com/gotk3/gotk3/gtk"
)

type screen struct {
	window   *gtk.Window
	Previous func()
}

func (sc *screen) init(title string) {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle(title)
	sc.window = win
}
