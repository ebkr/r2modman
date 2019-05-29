package modfetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type ThunderstoreVersion struct {
	Name           string
	Full_name      string
	Description    string
	Icon           string
	Version_number string
	Dependencies   []string
	Download_url   string
	Downloads      int
	Date_created   string
	Website_url    string
	Is_active      bool
	Uuid4          string
}

type ThunderstoreJSON struct {
	Name          string
	Full_name     string
	Owner         string
	Package_url   string
	Date_created  string
	Date_updated  string
	Uuid4         string
	Is_pinned     bool
	Is_deprecated bool
	Versions      []ThunderstoreVersion
}

var onlineMods []ThunderstoreJSON

// ThunderstoreGenerateList : Get latest version of thunderstore data
func ThunderstoreGenerateList() {
	if onlineMods != nil {
		return
	}
	store := []ThunderstoreJSON{}
	response, err := http.Get("https://thunderstore.io/api/v1/package/")
	if err != nil {
		fmt.Print(err.Error())
		onlineMods = store
	}
	text, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(text, &store)
	onlineMods = store
	response.Body.Close()
}

// ThunderstoreLocalToOnline : Get a local mod, and find the thunderstore equivalent
func ThunderstoreLocalToOnline(mod *Mod) {
	toTraverse := []ThunderstoreJSON{}
	for _, thunder := range onlineMods {
		if strings.Compare(thunder.Name, mod.Name) == 0 {
			toTraverse = append(toTraverse, thunder)
		}
	}
	thunderstoreNextLTODialog(mod, toTraverse)
}

// thunderstoreNextLTODialog : Move to the next LocalToOnline dialog if not accepted
func thunderstoreNextLTODialog(mod *Mod, traverse []ThunderstoreJSON) {
	if len(traverse) == 0 {
		return
	}
	thunder := traverse[0]
	dialog, _ := gtk.DialogNew()
	dialog.AddButton("Confirm", gtk.RESPONSE_ACCEPT)
	dialog.AddButton("Deny", gtk.RESPONSE_CANCEL)
	dialog.SetTitle("Link to Thunderstore")
	box, _ := dialog.GetContentArea()
	box.SetBorderWidth(10)
	modName, _ := gtk.LabelNew("Mod: " + mod.Name)
	authorName, _ := gtk.LabelNew("Author: " + thunder.Owner)
	box.PackStart(modName, false, false, 2)
	box.PackStart(authorName, false, false, 10)
	dialog.ShowAll()

	switch dialog.Run() {
	case gtk.RESPONSE_ACCEPT:
		mod.Uuid4 = thunder.Uuid4
		dialog.Destroy()
		return
	case gtk.RESPONSE_CANCEL:
		dialog.Destroy()
		thunderstoreNextLTODialog(mod, traverse[1:])
		return
	}

}

// ThunderstoreGetAll : Get all
func ThunderstoreGetAll() []ThunderstoreJSON {
	return onlineMods
}

func ThunderstoreModHasUpdate(mod *Mod) bool {
	for _, a := range onlineMods {
		if strings.Compare(a.Uuid4, mod.Uuid4) == 0 {
			versionNumbers := getVersion(a.Versions[0].Version_number)
			if versionNumbers.Major != mod.Version.Major || versionNumbers.Minor != mod.Version.Minor || versionNumbers.Patch != mod.Version.Patch {
				return true
			}
		}
	}
	return false
}

func ThunderstoreUpdateMod(mod *Mod) {}

// ThunderstoreReady : Check if values are initialised
func ThunderstoreReady() bool {
	return onlineMods == nil
}
