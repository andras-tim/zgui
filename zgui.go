package main

import (
	"go/build"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	checkZFS()

	os.Exit(StartGTKApplication())
	// ShowBlockStorage()
	// zfsPools()
	// zfsList()
	// datasetList()
}

func errorCheck(e error) {
	if e != nil {
		log.Panicln(e)
	}
}

// Use to get path from dev or install env
func getPath(file string) string {
	// Try dev env
	ex, err := os.Executable()
	errorCheck(err)

	devPath := path.Join(filepath.Dir(ex), file)
	if _, err := os.Stat(devPath); err == nil {
		return file
	}

	// Try install env
	installPath := path.Join(build.Default.GOPATH, "src/gitlab.com/beteras/zgui", file)
	if _, err := os.Stat(installPath); err == nil {
		return file
	}

	log.Panic("Can't find file: ", devPath, installPath)

	return ""
}
