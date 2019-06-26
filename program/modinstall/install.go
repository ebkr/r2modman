package modinstall

import (
	"encoding/json"
	"fmt"
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

type modfileReference struct {
	file os.FileInfo
	path string
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

	if !mod.Enabled {
		return nil
	}

	/*
		files, err := ioutil.ReadDir(mod.Path)
		if err != nil {
			return err
		}
	*/

	installedMods := getInstalled()
	installedDirectories := map[string]bool{}

	locationsToInstall := findLocationToInstallFiles(mod.Path)
	for location, files := range locationsToInstall {
		installPath := filepath.Join(bepPath, location, mod.Name, "/")
		fmt.Println("Installing mod:", mod.Name, "to:", installPath)
		for _, file := range files {
			copyErr := copy.Copy(file.path, installPath+"/"+file.file.Name())
			if copyErr != nil {
				fmt.Println("Path:", file.path)
				fmt.Println(copyErr.Error())
			}
		}
		fmt.Println("Install complete")
		installedDirectories[installPath] = true
	}

	/*
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
	*/

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

func findLocationToInstallFiles(basePath string) map[string][]*modfileReference {
	store := map[string][]*modfileReference{}
	dirs, err := ioutil.ReadDir(basePath)
	if err != nil {
		return store
	}
	setLocations := []string{"config", "core", "monomod", "patchers", "plugins"}
	for _, a := range dirs {
		if a.IsDir() {
			found := false
			for _, scope := range setLocations {
				if strings.Compare(a.Name(), scope) == 0 {
					found = true
					files, _ := ioutil.ReadDir(filepath.Join(basePath, "/"+scope+"/"))
					array := []*modfileReference{}
					for _, genRef := range files {
						array = append(array, &modfileReference{
							file: genRef,
							path: filepath.Join(basePath, scope, genRef.Name()),
						})
					}
					store[scope] = array
				}
			}
			if !found {
				for key, val := range findLocationToInstallFiles(filepath.Join(basePath, "/"+a.Name()+"/")) {
					exVal, exists := store[key]
					if !exists {
						exVal = []*modfileReference{}
					}
					store[key] = append(exVal, val...)
				}
			}
		} else {
			array, exists := store["plugins"]
			if !exists {
				array = []*modfileReference{}
			}
			store["plugins"] = append(array, &modfileReference{
				file: a,
				path: filepath.Join(basePath, a.Name()),
			})
		}
	}
	return store
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
