package main

/*
Parse arguments: directory
list files in directory
group them by size
check duplicates inside each group by hash
*/

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

func gatherFiles(dir string, files *[]File) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = gatherFiles(filepath.Join(dir, entry.Name()), files)
			if err != nil {
				return err
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				return err
			}

			*files = append(*files, File{
				Path: filepath.Join(dir, entry.Name()),
				Size: info.Size(),
			})
		}
	}

	return nil
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

	var files []File
	err = gatherFiles(dir, &files)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(files)
}

