package modfetch

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	r, err := zip.OpenReader(zipSource)
	defer r.Close()
	if err != nil {
		return nil
	}
	os.Mkdir(modDirectory+name, os.ModePerm)
	files := map[string]string{}
	for _, file := range r.File {
		if file.FileInfo().IsDir() {
			folderPath := filepath.Join(modDirectory+name+"/", file.Name)
			os.MkdirAll(folderPath, 0777)
		} else {
			read, _ := file.Open()
			defer read.Close()
			lowerFileName := strings.ToLower(file.Name)
			outputFile, _ := os.OpenFile(modDirectory+name+"/"+lowerFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			defer outputFile.Close()
			io.Copy(outputFile, read)
			files[lowerFileName] = modDirectory + name + "/" + lowerFileName
		}
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
