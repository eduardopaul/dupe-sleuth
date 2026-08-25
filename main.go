package main

import (
	"log"
)


func main() {
	options := Parse()

	var logFlag *bool = options.Log
	var dir string = options.Dir
	var concurrent *bool = options.Concurrent

	filesBySize, err := filterBySizes(dir)
	if err != nil {
		log.Fatal(err)
	}

	prune(filesBySize)

	if *logFlag {
		logSize(filesBySize)
	}

	filesByFirstBytes := filterByFirstBytes(filesBySize, 8)
	prune(filesByFirstBytes)

	if *logFlag {
		logBytes(filesByFirstBytes)
	}

	filesByHash := filterByHashes(filesByFirstBytes, concurrent)

	prune(filesByHash)

	printGroups(filesByHash)
}

