package filewatcher

import (
	"fmt"
	"os"
	"time"
)

// FileWatcher : Watch for changes made to a file
type FileWatcher struct {
	Path     string
	OnChange func()
}

// Watch : Detect file change
func (watcher *FileWatcher) Watch() {
	file, err := os.Stat(watcher.Path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for true {
		newFile, statErr := os.Stat(watcher.Path)
		if statErr != nil {
			fmt.Println(statErr.Error())
			break
		}
		if newFile.Size() != file.Size() {
			file = newFile
			watcher.OnChange()
		}
		<-time.After(10 * time.Millisecond)
	}

}
