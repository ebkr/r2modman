package modfetch

import (
	"bufio"
	"encoding/json"
	"os"
)

type Mod struct {
	Name string
	URL  string
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
	file, _ := os.Create("./mods/mods.json")
	b, _ := json.Marshal(mods)
	file.Write(b)
	file.Close()
}
