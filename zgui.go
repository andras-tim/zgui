package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/jaypipes/ghw"
)

func main() {
	checkZFS()

	os.Exit(StartGTKApplication())
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
		return devPath
	}

	vagrantPath := path.Join("/vagrant", file)
	if _, err := os.Stat(vagrantPath); err == nil {
		return vagrantPath
	}

	// Try install env
	installPath := path.Join(build.Default.GOPATH, "src/gitlab.com/beteras/zgui", file)
	if _, err := os.Stat(installPath); err == nil {
		return installPath
	}

	log.Panic("Can't find file: ", devPath, " ", vagrantPath, " ", installPath)

	return ""
}

func getImage(path string) *gdk.Pixbuf {
	path = getPath(fmt.Sprintf("icons/%s.png", path))
	image, err := gdk.PixbufNewFromFileAtScale(path, 20, 20, true)
	errorCheck(err)

	return image
}

// This looks for unix only...
func checkZFS() {
	path := "/dev/zfs"

	info, err := os.Stat(path)
	if err != nil {
		log.Panicln("ZFS not accessible or not loaded")
	}

	// 0b000000110 = 6
	if (info.Mode().Perm() & 6) != 6 {
		log.Panicln("Bad ZFS permission: current user allowed ?")
	}
}

// GetDiskByPartition Get disk by partition name
func GetDiskByPartition(name string) *ghw.Disk {
	block, err := ghw.Block()
	errorCheck(err)

	for _, disk := range block.Disks {
		for _, part := range disk.Partitions {
			if part.Name == name {
				return disk
			}
		}
	}

	return nil
}
