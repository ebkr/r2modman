package main

import (
	"fmt"
	"github.com/ebkr/r2modman/program/globals"
	"os"
	"path/filepath"
	"strings"

	"github.com/ebkr/r2modman/program/screens"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	isCmdLine := false
	if len(os.Args) > 1 {
		if strings.HasPrefix(os.Args[1], "ror2mm://") {
			isCmdLine = true
			globals.ROR2MMProtocol = strings.TrimPrefix(os.Args[1], "ror2mm://")
		}
	} else {
		fmt.Println("No args")
	}
	exec, _ := os.Executable()
	globals.ExecutableLocation = exec
	globals.RootDirectory = filepath.Dir(exec)
	gtk.Init(&os.Args)
	if !isCmdLine {
		splash := screens.SplashScreen{}
		splash.Show()
	} else {
		profiles := screens.ProfileScreen{}
		profiles.Show(true)
	}
	gtk.Main()
}
