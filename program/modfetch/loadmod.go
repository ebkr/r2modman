package modfetch

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
)

type modManifest struct {
	Name           string
	Version_number string
	Website_url    string
	Description    string
	Dependencies   []string
}

// Unzip : Unzips a file to /mods/ with a given name. Returns list of file paths.
func Unzip(name, zipSource string) map[string]string {
	r, err := zip.OpenReader(zipSource)
	if err != nil {
		return nil
	}
	os.Mkdir("./mods/"+name, os.ModePerm)
	files := map[string]string{}
	for _, file := range r.File {
		read, _ := file.Open()
		lowerFileName := strings.ToLower(file.Name)
		outputFile, _ := os.OpenFile("./mods/"+name+"/"+lowerFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		io.Copy(outputFile, read)
		outputFile.Close()
		read.Close()
		files[lowerFileName] = "./mods/" + name + "/" + lowerFileName
	}
	return files
}

func MakeModFromManifest(manifestFile, uuid string) mod {
	m, err := os.Open(manifestFile)
	if err != nil {
		return mod{}
	}
	scanner := bufio.NewScanner(m)
	text := ""
	for scanner.Scan() {
		text += scanner.Text()
	}
	m.Close()
	manifest := modManifest{}
	json.Unmarshal([]byte(text), &manifest)
	manifestPath := strings.Split(manifestFile, "/")
	folderDirectory := strings.Join(manifestPath[0:len(manifestPath)-1], "/")
	unzippedMod := createMod(
		manifest.Name,
		manifest.Description,
		manifest.Website_url,
		manifest.Version_number,
		folderDirectory,
		uuid,
	)
	unzippedMod.AddDependencies(manifest.Dependencies)
	return unzippedMod
}
