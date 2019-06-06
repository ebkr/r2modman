package modfetch

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/artdarek/go-unzip"
	"github.com/ebkr/r2modman/program/globals"
)

type modManifest struct {
	Name           string
	Version_number string
	Website_url    string
	Description    string
	Dependencies   []string
}

// Unzip : Unzips a file to /mods/ with a given name. Returns list of file paths
func Unzip(name, zipSource string) map[string]string {
	modDirectory := "./mods/" + globals.SelectedProfile + "/"
	files := map[string]string{}
	uz := unzip.New(zipSource, modDirectory+name)
	extractErr := uz.Extract()
	if extractErr != nil {
		fmt.Println("Error extracting mod:", extractErr.Error())
		return files
	}
	lookthrough, lookErr := ioutil.ReadDir(modDirectory + name)
	if lookErr != nil {
		fmt.Println("Lookthrough Error:", lookErr.Error())
		return files
	}
	for _, a := range lookthrough {
		files[a.Name()] = modDirectory + name + "/" + a.Name()
	}
	return files
}

// MakeModFromManifest : Create a mod from the included manifest file
func MakeModFromManifest(manifestFile, uuid string) Mod {
	m, err := os.Open(manifestFile)
	if err != nil {
		return Mod{}
	}
	scanner := bufio.NewScanner(m)
	text := ""
	for scanner.Scan() {
		text += scanner.Text()
	}
	m.Close()
	manifest := modManifest{}
	byteText := bytes.TrimPrefix([]byte(text), []byte("\xef\xbb\xbf"))
	jsonErr := json.Unmarshal(byteText, &manifest)
	if jsonErr != nil {
		fmt.Println("JSON Err:", jsonErr.Error())
	}
	manifestPath := strings.Split(manifestFile, "/")
	folderDirectory := strings.Join(manifestPath[0:len(manifestPath)-1], "/")
	unzippedMod := createMod(
		manifest.Name,
		manifest.Description,
		manifest.Website_url,
		manifest.Version_number,
		folderDirectory,
		uuid,
		manifest.Name,
	)
	unzippedMod.AddDependencies(manifest.Dependencies)
	return unzippedMod
}
