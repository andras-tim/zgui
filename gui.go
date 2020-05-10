package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	zfs "github.com/bicomsystems/go-libzfs"
	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jaypipes/ghw"
)

//go:generate go run generator.go Button Label ListStore ProgressBar TreeStore TreeView Window

const appID = "com.github.beteras.zgui"

// StartGTKApplication Start GTK application
func StartGTKApplication() int {
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	errorCheck(err)

	application.Connect("startup", func() {
		log.Println("application startup")
	})

	application.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	application.Connect("activate", func() {
		log.Println("application activate")
		var gui GUI
		gui.init(application)
	})

	// Launch the application
	return application.Run(os.Args)
}

// GUI GUI
type GUI struct {
	// GTK main
	application *gtk.Application
	builder     *gtk.Builder
	window      *gtk.Window

	// GTK datastore

	datasetPropertiesListStore *gtk.ListStore
	datasetTreeStore           *gtk.TreeStore
	poolFeaturesListStore      *gtk.ListStore
	poolListStore              *gtk.ListStore
	poolPropertiesListStore    *gtk.ListStore
	poolVDevsTreeStore         *gtk.TreeStore
	storageTreeStore           *gtk.TreeStore

	// GTK Widgets

	datasetPropertiesTreeView *gtk.TreeView
	datasetTreeView           *gtk.TreeView
	poolFeaturesTreeView      *gtk.TreeView
	poolPropertiesTreeView    *gtk.TreeView
	poolTreeView              *gtk.TreeView
	poolVdevsTreeView         *gtk.TreeView
	storageTreeView           *gtk.TreeView

	quitButton             *gtk.Button
	refreshButton          *gtk.Button
	poolStateLabel         *gtk.Label
	poolUsageProgressBar   *gtk.ProgressBar
	poolSizeLabel          *gtk.Label
	poolFreeLabel          *gtk.Label
	datasetAvailableLabel  *gtk.Label
	datasetCreationLabel   *gtk.Label
	datasetMountPointLabel *gtk.Label

	// GTK TreeView sort

	datasetModelSort *gtk.TreeModelSort
	poolModelSort    *gtk.TreeModelSort

	// GTK Icons

	iconDatasetClone      *gdk.Pixbuf
	iconDatasetFilesystem *gdk.Pixbuf
	iconDatasetSnapshot   *gdk.Pixbuf
	iconDatasetVolume     *gdk.Pixbuf
	iconStateBad          *gdk.Pixbuf
	iconStateOK           *gdk.Pixbuf
	iconStateWarning      *gdk.Pixbuf
	iconStorageHDD        *gdk.Pixbuf
	iconStorageNVMe       *gdk.Pixbuf
	iconStoragePartition  *gdk.Pixbuf
	iconStorageSSD        *gdk.Pixbuf
	iconStorageUSB        *gdk.Pixbuf
	iconZFSRaidZ          *gdk.Pixbuf
}

func (w *GUI) init(application *gtk.Application) {
	w.application = application

	// Get the Glade UI builder
	builder, err := gtk.BuilderNewFromFile(getPath("zgui.glade"))
	errorCheck(err)
	w.builder = builder

	w.initIcons()
	w.initGladeObjects()

	w.quitButton.Connect("clicked", application.Quit)
	w.refreshButton.Connect("clicked", w.refresh)

	w.refresh()

	w.window.Show()
	w.application.AddWindow(w.window)
}

func (w *GUI) initIcons() {
	w.iconDatasetClone = getImage("dataset-clone")
	w.iconDatasetFilesystem = getImage("dataset-filesystem")
	w.iconDatasetSnapshot = getImage("dataset-snapshot")
	w.iconDatasetVolume = getImage("dataset-volume")
	w.iconStateBad = getImage("state-bad")
	w.iconStateOK = getImage("state-ok")
	w.iconStateWarning = getImage("state-warning")
	w.iconStorageHDD = getImage("storage-hdd")
	w.iconStorageNVMe = getImage("storage-nvme")
	w.iconStoragePartition = getImage("storage-partition")
	w.iconStorageSSD = getImage("storage-ssd")
	w.iconStorageUSB = getImage("storage-usb")
	w.iconZFSRaidZ = getImage("zfs-raidz")
}

func (w *GUI) initGladeObjects() {
	w.datasetPropertiesListStore = w.getGtkListStore("datasetPropertiesListStore")
	w.datasetTreeStore = w.getGtkTreeStore("datasetTreeStore")
	w.poolFeaturesListStore = w.getGtkListStore("poolFeaturesListStore")
	w.poolListStore = w.getGtkListStore("poolListStore")
	w.poolPropertiesListStore = w.getGtkListStore("poolPropertiesListStore")
	w.poolVDevsTreeStore = w.getGtkTreeStore("poolVDevsTreeStore")
	w.storageTreeStore = w.getGtkTreeStore("storageTreeStore")

	w.datasetPropertiesTreeView = w.getGtkTreeView("datasetPropertiesTreeView")
	w.datasetTreeView = w.getGtkTreeView("datasetTreeView")
	w.poolFeaturesTreeView = w.getGtkTreeView("poolFeaturesTreeView")
	w.poolPropertiesTreeView = w.getGtkTreeView("poolPropertiesTreeView")
	w.poolTreeView = w.getGtkTreeView("poolTreeView")
	w.poolVdevsTreeView = w.getGtkTreeView("poolVdevsTreeView")
	w.storageTreeView = w.getGtkTreeView("storageTreeView")

	w.quitButton = w.getGtkButton("quitButton")
	w.refreshButton = w.getGtkButton("refreshButton")

	w.poolStateLabel = w.getGtkLabel("poolStateLabel")
	w.poolUsageProgressBar = w.getGtkProgressBar("poolUsageProgressBar")
	w.poolSizeLabel = w.getGtkLabel("poolSizeLabel")
	w.poolFreeLabel = w.getGtkLabel("poolFreeLabel")
	w.datasetAvailableLabel = w.getGtkLabel("datasetAvailableLabel")
	w.datasetCreationLabel = w.getGtkLabel("datasetCreationLabel")
	w.datasetMountPointLabel = w.getGtkLabel("datasetMountPointLabel")

	w.window = w.getGtkWindow("mainWindow")

	w.addTreeViewSelectionChangedEvent(w.poolTreeView, w.onPoolSelectionChanged, 1)
	w.addTreeViewSelectionChangedEvent(w.datasetTreeView, w.onDatasetSelectionChanged, 9)

	// Treeview sorting

	w.datasetModelSort = w.setTreeViewSortColumn(w.datasetTreeView, &w.datasetTreeStore.TreeModel, 0)
	w.setTreeViewSortColumn(w.datasetPropertiesTreeView, &w.datasetPropertiesListStore.TreeModel, 0)
	w.setTreeViewSortColumn(w.poolFeaturesTreeView, &w.poolFeaturesListStore.TreeModel, 0)

	w.poolModelSort = w.setTreeViewSortColumn(w.poolTreeView, &w.poolListStore.TreeModel, 1)
	w.setTreeViewSortColumn(w.poolPropertiesTreeView, &w.poolPropertiesListStore.TreeModel, 0)
	w.setTreeViewSortColumn(w.poolVdevsTreeView, &w.poolVDevsTreeStore.TreeModel, 2)
	w.setTreeViewSortColumn(w.storageTreeView, &w.storageTreeStore.TreeModel, 1)
}

func (w *GUI) setTreeViewSortColumn(treeView *gtk.TreeView, treeModel *gtk.TreeModel, column int) *gtk.TreeModelSort {
	tms, err := gtk.TreeModelSortNew(treeModel)
	errorCheck(err)

	tms.SetSortColumnId(column, gtk.SORT_ASCENDING)
	treeView.SetModel(tms)

	return tms
}

func (w *GUI) onDatasetSelectionChanged(name string) {
	w.datasetPropertiesListStore.Clear()

	dataset, err := zfs.DatasetOpenSingle(name)
	defer dataset.Close()
	errorCheck(err)

	for key, prop := range dataset.Properties {
		iter := w.datasetPropertiesListStore.Append()

		pkey := zfs.Prop(key)

		err = w.datasetPropertiesListStore.SetValue(iter, 0, zfs.DatasetPropertyToName(pkey))
		errorCheck(err)

		// TODO: find why: uint vs int 64 overflow ?
		value := prop.Value
		if value == "18446744073709551615" {
			value = "-1"
		}
		err = w.datasetPropertiesListStore.SetValue(iter, 1, value)
		errorCheck(err)

		err = w.datasetPropertiesListStore.SetValue(iter, 2, prop.Source)
		errorCheck(err)

		err = w.datasetPropertiesListStore.SetValue(iter, 3, "this is my tooltip: TOOLTIP\nBold ? <b>w00t</b>")
		errorCheck(err)
	}

	w.datasetAvailableLabel.SetText(dataset.Properties[zfs.DatasetPropAvailable].Value)

	creation, err := strconv.ParseInt(dataset.Properties[zfs.DatasetPropCreation].Value, 0, 64)
	errorCheck(err)

	w.datasetCreationLabel.SetText(time.Unix(creation, 0).String())

	w.datasetMountPointLabel.SetText(dataset.Properties[zfs.DatasetPropMountpoint].Value)

	// Convert size to unsigned integer
	sizeAvailaible, err := strconv.ParseUint(dataset.Properties[zfs.DatasetPropAvailable].Value, 0, 64)
	if err == nil {
		// Set size/free label text
		w.datasetAvailableLabel.SetText(humanize.Bytes(sizeAvailaible))
	}
}

func (w *GUI) onPoolSelectionChanged(name string) {
	pool, err := zfs.PoolOpen(name)
	defer pool.Close()
	errorCheck(err)

	// Get pool state
	state, err := pool.State()
	errorCheck(err)

	// Get pool as name
	w.poolStateLabel.SetText(zfs.PoolStateToName(state))

	// Convert size to float
	size, err := strconv.ParseFloat(pool.Properties[zfs.PoolPropSize].Value, 32)
	errorCheck(err)
	// Convert free space to float
	free, err := strconv.ParseFloat(pool.Properties[zfs.PoolPropFree].Value, 32)
	errorCheck(err)
	// Set used progress bar percentage
	w.poolUsageProgressBar.SetFraction(1 - (free / size))

	// Convert size to unsigned integer
	sizeI, err := strconv.ParseUint(pool.Properties[zfs.PoolPropSize].Value, 0, 64)
	errorCheck(err)

	// Convert free to unsigned integer
	freeI, err := strconv.ParseUint(pool.Properties[zfs.PoolPropFree].Value, 0, 64)
	errorCheck(err)
	// Set size/free label text
	w.poolSizeLabel.SetText(humanize.Bytes(sizeI))
	w.poolFreeLabel.SetText(humanize.Bytes(freeI))

	w.poolPropertiesListStore.Clear()
	for key, prop := range pool.Properties {
		pkey := zfs.Prop(key)
		iter := w.poolPropertiesListStore.Append()

		err := w.poolPropertiesListStore.Set(iter,
			[]int{0, 1, 2},
			[]interface{}{
				zfs.PoolPropertyToName(pkey),
				prop.Value,
				prop.Source})

		errorCheck(err)
	}

	w.poolFeaturesListStore.Clear()
	for key, val := range pool.Features {
		iter := w.poolFeaturesListStore.Append()

		err := w.poolFeaturesListStore.Set(iter,
			[]int{0, 1},
			[]interface{}{key, val})

		errorCheck(err)
	}

	w.poolVDevsTreeStore.Clear()
	vdevs, err := pool.VDevTree()
	errorCheck(err)
	w.vDevsStoreAdd(vdevs, nil)
	w.poolVdevsTreeView.ExpandAll()
}

type treeViewSelectionChangedEvent func(string)

func (w *GUI) addTreeViewSelectionChangedEvent(treeView *gtk.TreeView, fct treeViewSelectionChangedEvent, index int) {
	selection, err := treeView.GetSelection()
	errorCheck(err)

	selection.Connect("changed", func(selection *gtk.TreeSelection) {
		model, iter, ok := selection.GetSelected()
		if !ok {
			return
		}

		treeModel, ok := model.(*gtk.TreeModel)
		if !ok {
			return
		}

		value, err := treeModel.GetValue(iter, index)
		errorCheck(err)
		name, err := value.GetString()
		errorCheck(err)

		fct(name)
	})
}

func (w *GUI) vDevsStoreAdd(vt zfs.VDevTree, iter1 *gtk.TreeIter) {
	for _, v := range vt.Devices {
		// Get clean symlink path
		path, err := os.Readlink(v.Name)
		if err == nil {
			if !filepath.IsAbs(path) {
				path = filepath.Join(
					filepath.Dir(v.Name),
					path)
			}
		} else {
			path = v.Name
		}

		// TODO: What about other OS
		// $ go tool dist list
		if runtime.GOOS == "linux" {
			path = filepath.Base(path)
		}

		var icon *gdk.Pixbuf = nil
		var diskType string = ""

		if strings.HasPrefix(path, "raidz") {
			icon = w.iconZFSRaidZ
		} else {
			disk := GetDiskByPartition(path)
			if disk != nil {
				diskType = disk.DriveType.String()

				switch disk.DriveType {
				case ghw.DRIVE_TYPE_HDD:
					icon = w.iconStorageHDD
				case ghw.DRIVE_TYPE_SSD:
					icon = w.iconStorageSSD
					// default:
					// 	icon = w.iconUnknownStorage
				}
			}
		}

		var iconState *gdk.Pixbuf
		switch v.Stat.State {
		case zfs.VDevStateHealthy:
			iconState = w.iconStateOK
		case zfs.VDevStateDegraded:
			iconState = w.iconStateWarning
		default:
			iconState = w.iconStateBad
		}

		iter := w.poolVDevsTreeStore.Append(iter1)

		err = w.poolVDevsTreeStore.SetValue(iter, 0, iconState)
		errorCheck(err)

		if icon != nil {
			err = w.poolVDevsTreeStore.SetValue(iter, 1, icon)
			errorCheck(err)
		}

		err = w.poolVDevsTreeStore.SetValue(iter, 2, path)
		errorCheck(err)

		err = w.poolVDevsTreeStore.SetValue(iter, 3, string(v.Type))
		errorCheck(err)

		err = w.poolVDevsTreeStore.SetValue(iter, 4, diskType)
		errorCheck(err)

		err = w.poolVDevsTreeStore.SetValue(iter, 5, v.Stat.State.String())
		errorCheck(err)

		w.vDevsStoreAdd(v, iter)
	}
}

func (w *GUI) refresh() {
	log.Println("Refreshing all tabs...")

	w.refreshDatasetTab()
	w.refreshPoolTab()
	w.refreshStorageTab()
}

func (w *GUI) refreshDatasetTab() {
	w.datasetTreeStore.Clear()

	datasets, err := zfs.DatasetOpenAll()
	errorCheck(err)
	defer zfs.DatasetCloseAll(datasets)

	w.datasetStoreAdd(datasets, nil)

	w.datasetTreeView.ExpandAll()

	// Select first dataset
	if iter, ok := w.datasetModelSort.GetIterFirst(); ok {
		sel, err := w.datasetTreeView.GetSelection()
		errorCheck(err)

		sel.SelectIter(iter)
	}
}

func (w *GUI) datasetStoreAdd(datasets []zfs.Dataset, parentIter *gtk.TreeIter) {
	for _, dataset := range datasets {
		iter := w.datasetTreeStore.Append(parentIter)

		dpath, err := dataset.Path()
		errorCheck(err)
		dtype, err := dataset.GetProperty(zfs.DatasetPropType)
		errorCheck(err)
		dencryption, err := dataset.GetProperty(zfs.DatasetPropEncryption)
		errorCheck(err)
		dcompressratio, err := dataset.GetProperty(zfs.DatasetPropCompressratio)
		errorCheck(err)
		dused, err := dataset.GetProperty(zfs.DatasetPropUsed)
		errorCheck(err)

		var name string
		if strings.Contains(dpath, "@") {
			name = strings.Split(dpath, "@")[1]
		} else {
			name = path.Base(dpath)
		}
		err = w.datasetTreeStore.SetValue(iter, 0, name)
		errorCheck(err)

		{
			// TODO: Find why dataset.Type is always equal to zfs.DatasetTypeFilesystem
			var icon *gdk.Pixbuf = nil
			switch dtype.Value {
			case "filesystem":
				icon = w.iconDatasetFilesystem
			case "snapshot":
				icon = w.iconDatasetSnapshot
			case "volume":
				icon = w.iconDatasetVolume
			case "bookmark":
				icon = w.iconDatasetClone
			}

			err = w.datasetTreeStore.SetValue(iter, 1, icon)
			errorCheck(err)
		}

		err = w.datasetTreeStore.SetValue(iter, 2, dtype.Value)
		errorCheck(err)

		dmounted, err := dataset.GetProperty(zfs.DatasetPropMounted)
		if err == nil {
			err = w.datasetTreeStore.SetValue(iter, 3, dmounted.Value)
			errorCheck(err)
		}

		err = w.datasetTreeStore.SetValue(iter, 4, dencryption.Value)
		errorCheck(err)

		ddedup, err := dataset.GetProperty(zfs.DatasetPropDedup)
		if err == nil {
			err = w.datasetTreeStore.SetValue(iter, 5, ddedup.Value)
			errorCheck(err)
		}

		dcompression, err := dataset.GetProperty(zfs.DatasetPropCompression)
		if err == nil {
			err = w.datasetTreeStore.SetValue(iter, 6, dcompression.Value)
			errorCheck(err)
		}

		err = w.datasetTreeStore.SetValue(iter, 7, dcompressratio.Value)
		errorCheck(err)

		s, err := strconv.ParseUint(dused.Value, 10, 64)
		if (err == nil) && (s > 0) {
			err := w.datasetTreeStore.SetValue(iter, 8, humanize.Bytes(s))
			errorCheck(err)
		}

		err = w.datasetTreeStore.SetValue(iter, 9, dpath)
		errorCheck(err)

		w.datasetStoreAdd(dataset.Children, iter)
	}
}

func (w *GUI) refreshPoolTab() {
	w.poolListStore.Clear()

	// Lets open handles to all active pools on system
	pools, err := zfs.PoolOpenAll()
	errorCheck(err)

	for _, pool := range pools {
		defer pool.Close()

		name := pool.Properties[zfs.PoolPropName].Value
		state, err := pool.State()
		errorCheck(err)

		var icon *gdk.Pixbuf
		switch state {
		case zfs.PoolStateActive:
			icon = w.iconStateOK
		default:
			icon = w.iconStateBad
		}

		iter := w.poolListStore.Append()
		err = w.poolListStore.Set(iter,
			[]int{0, 1},
			[]interface{}{icon, name})

		errorCheck(err)
	}

	// Select first pool
	if iter, ok := w.poolModelSort.GetIterFirst(); ok {
		sel, err := w.poolTreeView.GetSelection()
		errorCheck(err)

		sel.SelectIter(iter)
	}
}

func (w *GUI) refreshStorageTab() {
	w.storageTreeStore.Clear()

	block, err := ghw.Block()
	errorCheck(err)

	for _, disk := range block.Disks {
		// Discard LVM disk
		if strings.HasPrefix(disk.Name, "dm-") {
			continue
		}

		iterDisk := w.storageTreeStore.Append(nil)

		// Get the correct disk icon
		var icon *gdk.Pixbuf = nil
		if disk.BusType == ghw.BUS_TYPE_NVME {
			icon = w.iconStorageNVMe
		} else if strings.Contains(disk.BusPath, "usb") {
			icon = w.iconStorageUSB
		} else {
			switch disk.DriveType {
			case ghw.DRIVE_TYPE_HDD:
				icon = w.iconStorageHDD
			case ghw.DRIVE_TYPE_SSD:
				icon = w.iconStorageSSD
				// default:
				// 	icon = w.iconUnknownStorage
			}
		}

		if icon != nil {
			err = w.storageTreeStore.SetValue(iterDisk, 0, icon)
			errorCheck(err)
		}

		err = w.storageTreeStore.SetValue(iterDisk, 1, disk.Name)
		errorCheck(err)

		err = w.storageTreeStore.SetValue(iterDisk, 2, disk.DriveType.String())
		errorCheck(err)

		err = w.storageTreeStore.SetValue(iterDisk, 3, humanize.Bytes(disk.SizeBytes))
		errorCheck(err)

		err = w.storageTreeStore.SetValue(iterDisk, 4, disk.StorageController.String())
		errorCheck(err)

		err = w.storageTreeStore.SetValue(iterDisk, 5, disk.PhysicalBlockSizeBytes)
		errorCheck(err)

		err = w.storageTreeStore.SetValue(iterDisk, 6, disk.Vendor)
		errorCheck(err)

		err = w.storageTreeStore.SetValue(iterDisk, 7, disk.Model)
		errorCheck(err)

		err = w.storageTreeStore.SetValue(iterDisk, 8, disk.SerialNumber)
		errorCheck(err)

		for _, part := range disk.Partitions {
			iterPart := w.storageTreeStore.Append(iterDisk)

			err = w.storageTreeStore.SetValue(iterPart, 0, w.iconStoragePartition)
			errorCheck(err)

			err = w.storageTreeStore.SetValue(iterPart, 1, part.Name)
			errorCheck(err)

			err = w.storageTreeStore.SetValue(iterPart, 2, part.Type)
			errorCheck(err)

			err = w.storageTreeStore.SetValue(iterPart, 3, humanize.Bytes(part.SizeBytes))
			errorCheck(err)

			if part.MountPoint != "" {
				mountStr := fmt.Sprintf(" mounted@%s", part.MountPoint)
				err = w.storageTreeStore.SetValue(iterPart, 7, mountStr)
				errorCheck(err)
			}
		}
	}
}
