package main

import (
	"fmt"
	"log"

	zfs "github.com/bicomsystems/go-libzfs"
)

func zfsPools() {
	// Lets open handles to all active pools on system
	pools, err := zfs.PoolOpenAll()
	if err != nil {
		println(err)
	}

	// Print each pool name and properties
	for _, p := range pools {
		// Print fancy header
		fmt.Printf("\n -----------------------------------------------------------\n")
		fmt.Printf("   POOL: %49s   \n", p.Properties[zfs.PoolPropName].Value)
		fmt.Printf("|-----------------------------------------------------------|\n")
		fmt.Printf("|  PROPERTY      |  VALUE                |  SOURCE          |\n")
		fmt.Printf("|-----------------------------------------------------------|\n")

		// Iterate pool properties and print name, value and source
		for key, prop := range p.Properties {
			pkey := zfs.Prop(key)
			if pkey == zfs.PoolPropName {
				continue // Skip name its already printed above
			}
			fmt.Printf("|%14s  | %20s  | %15s  |\n",
				zfs.PoolPropertyToName(pkey),
				prop.Value,
				prop.Source)
			println("")
		}
		println("")

		// Close pool handle and free memory, since it will not be used anymore
		p.Close()
	}

	datasets, err := zfs.DatasetOpenAll()
	if err != nil {
		panic(err.Error())
	}
	defer zfs.DatasetCloseAll(datasets)

	// Print out path and type of root datasets
	for _, d := range datasets {
		path, err := d.Path()
		if err != nil {
			panic(err.Error())
		}
		p, err := d.GetProperty(zfs.DatasetPropType)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%30s | %10s\n", path, p.Value)
	}
}

func getPoolState(poolName string) {
	p, err := zfs.PoolOpen(poolName)
	if err != nil {
		panic(err)
	}
	defer p.Close()
	pstate, err := p.State()
	if err != nil {
		panic(err)
	}
	fmt.Printf("POOL %s state: %s\n", poolName, zfs.PoolStateToName(pstate))
}

func printVDevTree(vt zfs.VDevTree, pref string) {
	first := pref + vt.Name
	fmt.Printf("%-30s | %-10s | %-10s | %s\n", first, vt.Type, vt.Stat.State.String(), vt.Path)

	for _, v := range vt.Devices {
		printVDevTree(v, "  "+pref)
	}

	if len(vt.Spares) > 0 {
		fmt.Println("spares:")
		for _, v := range vt.Spares {
			printVDevTree(v, "  "+pref)
		}
	}

	if len(vt.L2Cache) > 0 {
		fmt.Println("l2cache:")
		for _, v := range vt.L2Cache {
			printVDevTree(v, "  "+pref)
		}
	}
}

func getVdevTree(poolName string) {
	var vdevs zfs.VDevTree
	println("pool VDevTree")

	pool, err := zfs.PoolOpen(poolName)
	errorCheck(err)
	defer pool.Close()

	vdevs, err = pool.VDevTree()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("%-30s | %-10s | %-10s | %s\n", "NAME", "TYPE", "STATE", "PATH")
	printVDevTree(vdevs, "")
}

func zfsList() {
	// Lets open handles to all active pools on system
	pools, err := zfs.PoolOpenAll()
	if err != nil {
		println(err)
	}

	// Print each pool name and properties
	for _, p := range pools {
		// Print fancy header
		fmt.Printf("\n -----------------------------------------------------------\n")
		fmt.Printf("   POOL: %49s   \n", p.Properties[zfs.PoolPropName].Value)
		fmt.Printf("|-----------------------------------------------------------|\n")
		fmt.Printf("|  PROPERTY      |  VALUE                |  SOURCE          |\n")
		fmt.Printf("|-----------------------------------------------------------|\n")

		// Iterate pool properties and print name, value and source
		for key, prop := range p.Properties {
			pkey := zfs.Prop(key)
			if pkey == zfs.PoolPropName {
				continue // Skip name its already printed above
			}
			fmt.Printf(
				"|%14s  | %20s  | %15s  |\n",
				zfs.PoolPropertyToName(pkey),
				prop.Value,
				prop.Source)
		}
		getPoolState(p.Properties[zfs.PoolPropName].Value)
		getVdevTree(p.Properties[zfs.PoolPropName].Value)
		println("")

		// Iterate pool features
		for key, prop := range p.Features {
			fmt.Printf(
				"%s: %s\n",
				key,
				prop)
		}
		println("")

		// Close pool handle and free memory, since it will not be used anymore
		p.Close()
	}
}

func showChild(datasets []zfs.Dataset, depth int) {
	// for _, child := range datasets {
	// 	log.Println(child.Path())

	// 	snapshots, err := child.Snapshots()
	// 	errorCheck(err)
	// 	log.Println(snapshots)
	// 	for snapshot := range snapshots {
	// 		log.Println(snapshot)
	// 	}

	// 	showChild(child.Children)
	// }

	// Print out path and type of root datasets
	for _, d := range datasets {
		path, err := d.Path()
		errorCheck(err)

		p, err := d.GetProperty(zfs.DatasetPropName)
		errorCheck(err)

		fmt.Printf("*** Dataset: %d: %s: %s\n", depth, path, p.Value)
		// fmt.Println(d)

		// Iterate pool properties and print name, value and source
		for key, prop := range d.Properties {
			pkey := zfs.Prop(key)
			// fmt.Printf(
			// 	"|%14s  | %20s  | %15s  |\n",
			// 	zfs.DatasetPropertyToName(pkey),
			// 	prop.Value,
			// 	prop.Source)

			log.Println(fmt.Sprintf("%s: %s: %s", zfs.DatasetPropertyToName(pkey), prop.Value, prop.Source))
		}

		showChild(d.Children, depth+1)

		// snapshots, err := d.Snapshots()
		// errorCheck(err)
		// log.Println(snapshots)
		// for snapshot := range snapshots {
		// 	log.Println(snapshot)
		// }
	}
}

// func zfsGetDatasets() <-chan *zfs.Dataset {
// 	datasets, err := zfs.DatasetOpenAll()
// 	errorCheck(err)
// 	defer zfs.DatasetCloseAll(datasets)

// 	chnl := make(chan string)
// 	go func() {
// 		for scanner.Scan() {
// 			chnl <- scanner.Text()
// 		}
// 		close(chnl)
// 	}()

// 	return datasets
// }

func datasetList() {
	datasets, err := zfs.DatasetOpenAll()
	errorCheck(err)
	defer zfs.DatasetCloseAll(datasets)

	showChild(datasets, 0)
}
