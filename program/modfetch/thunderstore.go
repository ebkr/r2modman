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
func ThunderstoreLocalToOnline(mod *mod) {
	for _, thunder := range onlineMods {
		if strings.Compare(thunder.Name, mod.Name) == 0 {
			waiter := make(chan bool)
			dialog, _ := gtk.DialogNew()
			confirm, _ := dialog.AddButton("Confirm", gtk.RESPONSE_ACCEPT)
			cancel, _ := dialog.AddButton("Deny", gtk.RESPONSE_CANCEL)
			dialog.SetTitle("Confirm Correct Mod")
			box, _ := dialog.GetContentArea()
			box.SetBorderWidth(10)
			modName, _ := gtk.LabelNew("Mod: " + mod.Name)
			authorName, _ := gtk.LabelNew("Author: " + thunder.Owner)
			box.PackStart(modName, false, false, 2)
			box.PackStart(authorName, false, false, 2)
			dialog.ShowAll()

			confirm.Connect("clicked", func() {
				mod.Uuid4 = thunder.Uuid4
				waiter <- true
			})
			cancel.Connect("clicked", func() {
				waiter <- false
			})
			res := <-waiter
			if res {
				break
			}
		}
	}
}

// ThunderstoreGetAll : Get all
func ThunderstoreGetAll() []ThunderstoreJSON {
	return onlineMods
}

func ThunderstoreUpdateMod(mod *mod) {}
