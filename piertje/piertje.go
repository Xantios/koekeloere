// Watch a list of directorys recursively
package piertje

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

var paths []string
var watcher *fsnotify.Watcher
var channel chan string

var log logrus.Logger
var verbose bool = false

// @TODO: Move this to main and use a setter for consistency and configurability
var filter = []string{".git", "node_modules", "vendor"}

func SetLogger(instance *logrus.Logger) {
	log = *instance
}

func SetVerbose(verb *bool) {
	if *verb {
		verbose = true
	}
}

func SetChannel(_channel chan string) {
	channel = _channel
}

func SetPaths(input string) {

	tmpPaths := strings.Split(input, ",")

	for i := 0; i < len(tmpPaths); i++ {

		// Drop empty
		if tmpPaths[i] == "" {
			continue
		}

		// Resolve to absolute path
		path, err := filepath.Abs(tmpPaths[i])

		if err != nil {
			log.Errorf("Cant add %s, %s", tmpPaths[i], err.Error())
			continue
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Warnf("Cant add %s, because it does not exist", path)
			continue
		}

		if verbose {
			log.Infof("Adding %s\n", path)
		}

		paths = append(paths, path)
	}
}

func GetPaths() []string {
	return paths
}

func Run() {

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Panic("Cant create FS Watcher: %s", err.Error())
	}

	defer watcher.Close()

	if verbose {
		log.Info("Populating watcher")
	}

	for _, path := range paths {
		if err := filepath.Walk(path, filterDir); err != nil {
			log.Warnf("Cant add %s: %s", path, err.Error())
			continue
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:

			if !ok {
				log.Errorf("SynChannel closed %#v", event)
				channel <- "closed"
				return
			}

			// log.Printf("Event: %#v", event)

			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Infof("write: %s", event.Name)
				channel <- fmt.Sprintf("write:%s", event.Name)
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Infof("create: %s", event.Name)
				if isDirectory(event.Name) {
					watcher.Add(event.Name)
				}
				channel <- fmt.Sprintf("create:%s", event.Name)
				continue
			}

			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				log.Infof("chmod: %s", event.Name)
				channel <- fmt.Sprintf("chmod:%s", event.Name)
				continue
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Infof("remove: %s", event.Name)
				if isDirectory(event.Name) {
					watcher.Remove(event.Name)
				}
				channel <- fmt.Sprintf("remove:%s", event.Name)
				continue
			}

			if event.Op&fsnotify.Rename == fsnotify.Rename {
				log.Infof("rename: %s", event.Name)
				// @TODO: Figure out what to do when renamed
				channel <- fmt.Sprintf("rename:%s", event.Name)
				continue
			}

			log.Infof("unknown event: %#v", event)

		case err, ok := <-watcher.Errors:

			if !ok {
				log.Errorf("ErrChannel closed %s", err.Error())
				return
			}

			log.Errorf("%s\n", err)
		}
	}
}

func filterDir(path string, fi os.FileInfo, err error) error {

	// Filter
	if checkFilter(path) {
		return nil
	}

	// Default behaviour is to watch all files, so we only monitor dirs
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

// checkFilter Check if we should filter
func checkFilter(path string) bool {
	for _, item := range filter {
		if strings.Contains(path, "/"+item) {
			return true
		}
	}

	return false
}

// isDirectory Check if directory
func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
