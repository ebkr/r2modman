package modfetch

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ebkr/r2modman/program/screens/dialogs"

	"github.com/gotk3/gotk3/gdk"

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
var modPixbufs map[string]*gdk.Pixbuf

// ThunderstoreGenerateList : Get latest version of thunderstore data
func ThunderstoreGenerateList(progression chan float64) {
	if onlineMods != nil {
		return
	}
	store := []ThunderstoreJSON{}
	response, err := http.Get("https://thunderstore.io/api/v1/package/")
	if err != nil {
		fmt.Print(err.Error())
		onlineMods = store
		modPixbufs = map[string]*gdk.Pixbuf{}
		progression <- 0
		return
	}
	text, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(text, &store)
	onlineMods = store
	response.Body.Close()
	bufs := map[string]*gdk.Pixbuf{}
	progression <- 0
	modLength := len(onlineMods)
	for i, a := range onlineMods {
		if len(a.Versions[0].Icon) > 0 {
			pbloader, _ := gdk.PixbufLoaderNew()
			pbloader.SetSize(32, 32)
			imageResponse, _ := http.Get(a.Versions[0].Icon)
			imageBytes, _ := ioutil.ReadAll(imageResponse.Body)
			imageResponse.Body.Close()
			pbloader.Write(imageBytes)
			pix, _ := pbloader.GetPixbuf()
			bufs[a.Uuid4] = pix
		} else {
			bufs[a.Uuid4] = nil
		}
		progression <- float64(i) / float64(modLength)
	}
	modPixbufs = bufs
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

// ThunderstoreReady : Check if values are initialised
func ThunderstoreReady() bool {
	return (onlineMods != nil && modPixbufs != nil)
}

// ThunderstoreGetPixbufFromUUID4 : Get pixbuf to be used for thumbnails
func ThunderstoreGetPixbufFromUUID4(uuid4 string) *gdk.Pixbuf {
	return modPixbufs[uuid4]
}

// ThunderstoreDownloadMod : Download a mod directly from the store.
func ThunderstoreDownloadMod(uuid string, window *gtk.Window) *Mod {
	listener := make(chan *gtk.Dialog)
	result := make(chan *Mod)
	go func() {
		dialog := <-listener
		for _, a := range onlineMods {
			if a.Uuid4 == uuid {
				stream, err := http.Get(a.Versions[0].Download_url)
				if err != nil {
					break
				}
				defer stream.Body.Close()
				created, creationErr := os.Create("./mods/" + a.Versions[0].Full_name + ".zip")
				if creationErr != nil {
					break
				}
				_, copyErr := io.Copy(created, stream.Body)
				if copyErr != nil {
					break
				}

				res := Unzip(a.Full_name, "./mods/"+a.Versions[0].Full_name+".zip")
				val, exists := res["manifest.json"]
				if exists {
					created.Close()
					mod := MakeModFromManifest(val, "")
					mod.FullName = a.Full_name
					deleteErr := os.RemoveAll("./mods/" + a.Versions[0].Full_name + ".zip")
					if deleteErr != nil {
						fmt.Println(deleteErr.Error())
					}
					dialog.Destroy()
					result <- &mod
					return
				}
				dialog.Destroy()
			}
		}
		result <- nil
		dialog.Destroy()
	}()
	dialogs.NewAwaitDialog(listener, "Downloading", "Downloading Mod")
	res := <-result
	window.Present()
	return res
}

// ThunderstoreUpdateMod : Update a mod
func ThunderstoreUpdateMod(mod *Mod, window *gtk.Window) *Mod {
	if len(mod.Uuid4) == 0 {
		return mod
	}
	newMod := ThunderstoreDownloadMod(mod.Uuid4, window)
	if newMod == nil {
		return mod
	}
	os.RemoveAll(mod.Path)
	newMod.Uuid4 = mod.Uuid4
	return newMod
}

// ThunderstoreGetModByUUID4 : Return the mod associated by the UUID4
func ThunderstoreGetModByUUID4(uuid4 string) *ThunderstoreJSON {
	for _, a := range onlineMods {
		if strings.Compare(a.Uuid4, uuid4) == 0 {
			return &a
		}
	}
	return nil
}

// ThunderstoreGetDependency : Used to download any missing dependencies
func ThunderstoreGetDependency(modName string, window *gtk.Window) *Mod {
	for _, a := range onlineMods {
		if strings.Compare(a.Full_name, modName) == 0 {
			dialog, _ := gtk.DialogNew()
			box, _ := dialog.GetContentArea()
			box.SetBorderWidth(10)

			question, _ := gtk.LabelNew("Do you want to install missing dependency:")
			tag, _ := gtk.LabelNew(a.Name)

			box.PackStart(question, false, false, 5)
			box.PackStart(tag, false, false, 10)
			dialog.SetBorderWidth(10)
			yes, _ := dialog.AddButton("Yes", gtk.RESPONSE_YES)
			no, _ := dialog.AddButton("No", gtk.RESPONSE_NO)
			yes.SetHAlign(gtk.ALIGN_START)
			no.SetHAlign(gtk.ALIGN_START)

			dialog.SetPosition(gtk.WIN_POS_CENTER)
			dialog.SetTitle("Install Dependency")

			dialog.ShowAll()

			switch dialog.Run() {
			case gtk.RESPONSE_YES:
				dialog.Destroy()
				mod := ThunderstoreDownloadMod(a.Uuid4, window)
				mod.Uuid4 = a.Uuid4
				return mod
			case gtk.RESPONSE_NO:
				dialog.Destroy()
			}
			break
		}
	}
	return nil
}
