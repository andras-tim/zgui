// Auto generated file

package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func (w *GUI) getGtkButton(name string) *gtk.Button {
	obj, err := w.builder.GetObject(name)
	errorCheck(err)

	if win, ok := obj.(*gtk.Button); ok {
		return win
	}

	log.Panic()
	return nil
}

func (w *GUI) getGtkLabel(name string) *gtk.Label {
	obj, err := w.builder.GetObject(name)
	errorCheck(err)

	if win, ok := obj.(*gtk.Label); ok {
		return win
	}

	log.Panic()
	return nil
}

func (w *GUI) getGtkListStore(name string) *gtk.ListStore {
	obj, err := w.builder.GetObject(name)
	errorCheck(err)

	if win, ok := obj.(*gtk.ListStore); ok {
		return win
	}

	log.Panic()
	return nil
}

func (w *GUI) getGtkProgressBar(name string) *gtk.ProgressBar {
	obj, err := w.builder.GetObject(name)
	errorCheck(err)

	if win, ok := obj.(*gtk.ProgressBar); ok {
		return win
	}

	log.Panic()
	return nil
}

func (w *GUI) getGtkTreeStore(name string) *gtk.TreeStore {
	obj, err := w.builder.GetObject(name)
	errorCheck(err)

	if win, ok := obj.(*gtk.TreeStore); ok {
		return win
	}

	log.Panic()
	return nil
}

func (w *GUI) getGtkTreeView(name string) *gtk.TreeView {
	obj, err := w.builder.GetObject(name)
	errorCheck(err)

	if win, ok := obj.(*gtk.TreeView); ok {
		return win
	}

	log.Panic()
	return nil
}

func (w *GUI) getGtkWindow(name string) *gtk.Window {
	obj, err := w.builder.GetObject(name)
	errorCheck(err)

	if win, ok := obj.(*gtk.Window); ok {
		return win
	}

	log.Panic()
	return nil
}

