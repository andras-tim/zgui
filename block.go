package main

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

// ShowBlockStorage Show all blocks storage
func ShowBlockStorage() {
	block, err := ghw.Block()
	errorCheck(err)

	fmt.Printf("%v\n", block)

	for _, disk := range block.Disks {
		fmt.Printf(" %v\n", disk)
		fmt.Printf(" %v\n", disk.PhysicalBlockSizeBytes)

		// var stat syscall.Stat_t
		// syscall.Stat("/dev/"+disk.Name, &stat)
		// fmt.Printf("syscall: %s\n", stat.Blocks)

		for _, part := range disk.Partitions {
			fmt.Printf("  %v\n", part)
		}
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
