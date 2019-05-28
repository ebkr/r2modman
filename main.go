package main

import (
	"fmt"
	"os"

	"github.com/ebkr/r2modman/program/screens"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(&os.Args)
	awaitSplash := make(chan bool)
	splash := screens.NewSplashScreen()
	splash.Show(awaitSplash)
	gtk.Main()
	fmt.Println("Received")
	mainScreen := screens.ManagerScreen{}
	mainScreen.Show()
	gtk.Main()
}
