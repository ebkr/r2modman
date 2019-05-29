package modfetch

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// ModVersion : Track the mod version
type modVersion struct {
	Major int
	Minor int
	Patch int
}

type modDependency struct {
	Name    string
	Version modVersion
}

// Mod : Struct to contain mod information
type Mod struct {
	Name         string
	Description  string
	URL          string
	Version      modVersion
	Dependencies []modDependency
	Path         string
	Uuid4        string
}

// CreateMod : Used to create a mod object
func createMod(name, description, url, version, path, uuid string) Mod {
	return Mod{
		Name:         name,
		Description:  description,
		URL:          url,
		Version:      getVersion(version),
		Dependencies: []modDependency{},
		Path:         path,
		Uuid4:        uuid,
	}
}

// AddDependencies : Used to add dependencies from a list
func (mod *Mod) AddDependencies(dependencies []string) {
	list := []modDependency{}
	for _, a := range dependencies {
		split := strings.Split(a, "-")
		modName := strings.Join(split[0:len(split)-1], "-")
		version := getVersion(split[len(split)-1])
		list = append(list, modDependency{
			Name:    modName,
			Version: version,
		})
	}
	mod.Dependencies = list
}

func getVersion(v string) modVersion {
	versionNumbers := []int{}
	vn := strings.Split(v, ".")
	for _, num := range vn {
		conv, _ := strconv.Atoi(num)
		versionNumbers = append(versionNumbers, conv)
	}
	return modVersion{
		Major: versionNumbers[0],
		Minor: versionNumbers[1],
		Patch: versionNumbers[2],
	}
}

// GetMods : Get an array of mods
func GetMods() []Mod {
	file, fErr := os.Open("./mods/mods.json")
	if os.IsNotExist(fErr) {
		file, err := os.Create("./mods/mods.json")
		if err != nil {
			return []Mod{}
		}
		file.Write([]byte("{}"))
		file.Close()
		return GetMods()
	} else {
		scanner := bufio.NewScanner(file)
		text := ""
		for scanner.Scan() {
			text += scanner.Text()
		}
		data := []Mod{}
		json.Unmarshal([]byte(text), &data)
		file.Close()
		return data
	}
}

// UpdateMods : Update the list of mods
func UpdateMods(mods []Mod) {
	newList := []Mod{}
	for _, moda := range mods {
		found := false
		for _, modb := range newList {
			if strings.Compare(moda.Name, modb.Name) == 0 {
				found = true
				break
			}
		}
		if !found {
			newList = append(newList, moda)
		}
	}
	file, _ := os.Create("./mods/mods.json")
	b, _ := json.Marshal(newList)
	file.Write(b)
	file.Close()
}
