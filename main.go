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
	if len(os.Args) > 1 {
		if strings.HasPrefix(os.Args[1], "ror2mm://") {
			globals.ROR2MMProtocol = strings.TrimPrefix(os.Args[1], "ror2mm://")
		}
	} else {
		fmt.Println("No args")
	}
	exec, _ := os.Executable()
	globals.ExecutableLocation = exec
	globals.RootDirectory = filepath.Dir(exec)
	gtk.Init(&os.Args)
	splash := screens.SplashScreen{}
	splash.Show()
	gtk.Main()
}
