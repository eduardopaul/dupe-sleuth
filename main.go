package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Path string
	Size int64
}

func gatherFiles(dir string, filesBySize map[int64][]File) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = gatherFiles(filepath.Join(dir, entry.Name()), filesBySize)
			if err != nil {
				return err
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			
			size := info.Size()

			filesBySize[size] = append(filesBySize[size], File{
				Path: filepath.Join(dir, entry.Name()),
				Size: size,
			})
		}
	}

	return nil
}

func prune(filesBySize map[int64][]File) {
	for key, files := range filesBySize {
		if len(files) < 2 {
			delete(filesBySize, key)
		}
	}
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("You must provide the root path the search should performed in.")
		os.Exit(1)
	} else if flag.NArg() > 1 {
		fmt.Println("Too many arguments: you should provide just the root path the search should be performed in.")
		os.Exit(1)
	}

	dir := filepath.Clean(flag.Arg(0))

	info, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Specified directory does not exist.")
			os.Exit(1)
		} 

		log.Fatal(err)
	}

	if !info.IsDir() {
		fmt.Println("Specified path is not a directory.")
		os.Exit(1)
	}

	filesBySize := map[int64][]File{}

	err = gatherFiles(dir, filesBySize)
	if err != nil {
		log.Fatal(err)
	}

	prune(filesBySize)

	// check hash

	fmt.Println(filesBySize)
}

