package main

import (
	"log"
	"os"
)

func main() {
	os.Exit(StartGTKApplication())
	// ShowBlockStorage()
	// zfsPools()
	// zfsList()
	// datasetList()
}

func errorCheck(e error) {
	if e != nil {
		log.Panic(e)
	}
}
