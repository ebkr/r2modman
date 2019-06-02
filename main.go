package main

import (
	"os"

	"github.com/ebkr/r2modman/program/screens"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(&os.Args)
	splash := screens.SplashScreen{}
	splash.Show()
	gtk.Main()
}
