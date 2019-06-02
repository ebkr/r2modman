package modinstall

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"

	"github.com/ebkr/r2modman/program/modfetch"
)

type installedMod struct {
	Path string
}

// PrepareInstall : Used to remove previous mod installs.
func PrepareInstall() error {
	for _, a := range getInstalled() {
		removeErr := os.RemoveAll(a.Path)
		if removeErr != nil {
			return removeErr
		}
	}
	return updateInstalled([]installedMod{})
}

// InstallMod : Used to install mods in BepInEx/ directory. Compatibility layer for BepInEx sub-directories.
func InstallMod(mod *modfetch.Mod, bepPath string) error {

	files, err := ioutil.ReadDir(mod.Path)
	if err != nil {
		return err
	}

	installedMods := getInstalled()
	installedDirectories := map[string]bool{}

	for _, file := range files {
		path := ""
		var copyErr error
		if file.IsDir() {
			if isScopableFolder(file) {
				path = filepath.Join(filepath.Join(bepPath, file.Name(), mod.Name))
				copyErr = copy.Copy(filepath.Join(mod.Path, file.Name()), path)
			}
		} else {
			path = filepath.Join(bepPath, "plugins", mod.Name)
			copyErr = copy.Copy(filepath.Join(mod.Path, file.Name()), filepath.Join(path, file.Name()))
		}
		if copyErr != nil {
			return copyErr
		}
		installedDirectories[path] = true
	}

	for k := range installedDirectories {
		installedMods = append(installedMods, installedMod{k})
	}

	updateInstalled(installedMods)

	return nil
}

// Return true if folder is a directory of BepInEx
func isScopableFolder(file os.FileInfo) bool {
	scopable := []string{"config", "core", "monomod", "patchers", "plugins"}
	for _, a := range scopable {
		if strings.Compare(strings.ToLower(file.Name()), a) == 0 {
			return true
		}
	}
	return false
}

// Update the list of installed mods
func updateInstalled(installed []installedMod) error {
	data, parseErr := json.Marshal(installed)
	if parseErr != nil {
		return parseErr
	}
	file, createErr := os.Create("./mods/installed.json")
	if createErr != nil {
		return createErr
	}
	defer file.Close()
	file.Write(data)
	return nil
}

// Get the list of installed mod paths
func getInstalled() []installedMod {
	data := []installedMod{}
	file, fileErr := os.Open("./mods/installed.json")
	if fileErr != nil {
		return data
	}
	defer file.Close()
	read, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return data
	}
	json.Unmarshal(read, &data)
	return data
}
