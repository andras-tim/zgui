// +build generator

package main

import (
	"log"
	"os"
	"text/template"
)

func generateGtkWidgetsType(types []string) {
	path := "generator_gtk.tmpl"
	tmpl := template.Must(template.New(path).ParseFiles(path))

	out, err := os.Create("generator_gtk.go")
	errorCheck(err)

	err = tmpl.Execute(out, types)
	errorCheck(err)
}

func main() {
	log.Println("Generate gtk Widgets type")

	generateGtkWidgetsType(os.Args[1:])
}

func errorCheck(e error) {
	if e != nil {
		log.Panic(e)
	}
}
