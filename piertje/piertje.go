// Watch a list of directorys recursively 
package piertje

import (
	"os"
	"strings"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

var paths []string
var log logrus.Logger
var verbose bool = false

func SetLogger(instance *logrus.Logger) {
	log = *instance
}

func SetVerbose(verb *bool) {
	if *verb {
		verbose = true
	}
}

func SetPaths(input string) {

	tmpPaths := strings.Split(input,",")
	
	for i:=0;i<len(tmpPaths);i++ {

		// Drop empty
		if tmpPaths[i] == "" {
			continue
		}

		// Resolve to absolute path
		path,err := filepath.Abs(tmpPaths[i])

		if err != nil {
			log.Errorf("Cant add %s, %s",tmpPaths[i],err.Error())
			continue
		}

		if _,err := os.Stat(path); os.IsNotExist(err) {
			log.Warnf("Cant add %s, because it does not exist",path)
			continue
		}

		if verbose {
			log.Infof("Adding %s\n",path)
		}

		paths = append(paths,path)
	}
}

func GetPaths() []string {
	return paths
}
