package main

import (
	"os"
	"time"

	"github.com/ebkr/r2modman/program/screens"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(&os.Args)
	mainScreen := screens.ManagerScreen{}
	mainScreen.Show()
	time.Sleep(100)
	gtk.Main()
}
