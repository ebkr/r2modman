package modfetch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ebkr/r2modman/program/globals"
)

// ModVersion : Track the mod version
type modVersion struct {
	Major int
	Minor int
	Patch int
}

// String : Return a string object of the mod version.
func (mv *modVersion) String() string {
	return strconv.Itoa(mv.Major) + "." + strconv.Itoa(mv.Minor) + "." + strconv.Itoa(mv.Patch)
}

// ModDependency : Keep track of dependencies.
type ModDependency struct {
	Name    string
	Version modVersion
}

// Mod : Struct to contain mod information
type Mod struct {
	FullName     string
	Name         string
	Description  string
	URL          string
	Version      modVersion
	Dependencies []ModDependency
	Path         string
	Uuid4        string
	Enabled      bool
}

// CreateMod : Used to create a mod object
func createMod(name, description, url, version, path, uuid, fullName string) Mod {
	return Mod{
		Name:         name,
		Description:  description,
		URL:          url,
		Version:      getVersion(version),
		Dependencies: []ModDependency{},
		Path:         path,
		Uuid4:        uuid,
		FullName:     fullName,
		Enabled:      true,
	}
}

// AddDependencies : Used to add dependencies from a list
func (mod *Mod) AddDependencies(dependencies []string) {
	list := []ModDependency{}
	for _, a := range dependencies {
		split := strings.Split(a, "-")
		modName := strings.Join(split[0:len(split)-1], "-")
		version := getVersion(split[len(split)-1])
		list = append(list, ModDependency{
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
	for i := len(versionNumbers); i < 3; i++ {
		versionNumbers = append(versionNumbers, 0)
	}
	return modVersion{
		Major: versionNumbers[0],
		Minor: versionNumbers[1],
		Patch: versionNumbers[2],
	}
}

// GetMods : Get an array of mods
func GetMods() []Mod {
	modDirectory := "./mods/" + globals.SelectedProfile + "/"
	file, fErr := os.Open(modDirectory + "mods.json")
	if os.IsNotExist(fErr) {
		file, err := os.Create(modDirectory + "mods.json")
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
	modDirectory := "./mods/" + globals.SelectedProfile + "/"
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
	file, _ := os.Create(modDirectory + "mods.json")
	b, _ := json.Marshal(newList)
	file.Write(b)
	file.Close()
}

// RemoveMod : Remove a mod from the mods directory, and json file.
func RemoveMod(mod *Mod) {
	refreshedMods := GetMods()
	index := -1
	for i, a := range refreshedMods {
		if strings.Compare(mod.Name, a.Name) == 0 {
			err := os.RemoveAll(mod.Path)
			if err != nil {
				fmt.Println(err.Error())
			}
			index = i
		}
	}
	if index >= 0 {
		refreshedMods = append(refreshedMods[:index], refreshedMods[index+1:]...)
		UpdateMods(refreshedMods)
	}
}

// DependencyExists : Check if a dependency is installed
func (mod *Mod) DependencyExists(dependency *ModDependency) bool {
	mods := GetMods()
	for _, a := range mods {
		if a.Enabled && strings.Compare(a.FullName, dependency.Name) == 0 {
			depVer := dependency.Version
			if a.Version.Major > depVer.Major {
				return true
			} else if a.Version.Major == depVer.Major && a.Version.Minor > depVer.Minor {
				return true
			} else if a.Version.Major == depVer.Major && a.Version.Minor == depVer.Minor && a.Version.Patch >= depVer.Patch {
				return true
			}
			return false
		}
	}
	return strings.Compare(strings.ToLower(dependency.Name), "bbepis-bepinexpack") == 0
}
