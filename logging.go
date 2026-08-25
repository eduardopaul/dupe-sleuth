package main

import "log"

func logSize(filesBySize map[int64][]File) {
	log.Println()
	log.Println("Duplicate groups after size filtering:")
	for size, sliceOfFile := range filesBySize {
		log.Printf("	%d bytes\n", size)
		for _, file := range sliceOfFile {
			log.Println("	-", file.Path)
		}
		log.Println()
	}
}

func logBytes(filesByFirstBytes map[string][]File) {
	log.Println("Duplicate groups after first-bytes filtering:")
	for firstBytes, sliceOfFile := range filesByFirstBytes {
		log.Printf("	\"%s\"\n", firstBytes)
		for _, file := range sliceOfFile {
			log.Println("	-", file.Path)
		}
		log.Println()
	}
}
