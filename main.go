package main

import (
	"fmt"
	"log"
)

func main() {
	options := Parse()

	var dir string = options.Dir
	var concurrent *bool = options.Concurrent

	filesBySize, err := filterBySizes(dir)
	if err != nil {
		log.Fatal(err)
	}

	prune(filesBySize)

	fmt.Println("Duplicate groups after size filtering:")
	for size, sliceOfFile := range filesBySize {
		fmt.Printf("	%d bytes\n", size)
		for _, file := range sliceOfFile {
			fmt.Println("	-", file.Path)
		}
		fmt.Println()
	}

	filesByFirstBytes := filterByFirstBytes(filesBySize, 8)
	prune(filesByFirstBytes)

	fmt.Println("Duplicate groups after first-bytes filtering:")
	for firstBytes, sliceOfFile := range filesByFirstBytes {
		fmt.Printf("	\"%s\"\n", firstBytes)
		for _, file := range sliceOfFile {
			fmt.Println("	-", file.Path)
		}
		fmt.Println()
	}

	filesByHash := filterByHashes(filesByFirstBytes, concurrent)

	prune(filesByHash)

	fmt.Println("Duplicate groups after hash filtering:")
	for hash, sliceOfFile := range filesByHash {
		fmt.Printf("	%s\n", hash)
		for _, file := range sliceOfFile {
			fmt.Println("	-", file.Path)
		}
		fmt.Println()
	}
}

