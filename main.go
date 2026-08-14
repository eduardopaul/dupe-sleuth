package main

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Path string
	Size int64
	Hash string
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

func prune[K comparable](filesByProperty map[K][]File) {
	for key, files := range filesByProperty {
		if len(files) < 2 {
			delete(filesByProperty, key)
		}
	}
}

func getHashes(filesBySize map[int64][]File) (map[string][]File, error) {
	filesByHash := make(map[string][]File)

	for _, files := range filesBySize {
		for _, file := range files {
			h := md5.New()

			f, err := os.Open(file.Path)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			for {
				var n int64 = 1024
				_, err := io.CopyN(h, f, n)
				if err != nil {
					if err != io.EOF {
						log.Println("Encountered error %w when hashing file %s.", err, file)
					}

					break
				}
			}

			hash := fmt.Sprintf("%x", h.Sum(nil))

			filesByHash[hash] = append(filesByHash[hash], File{
				Path: file.Path,
				Size: file.Size,
				Hash: hash,
			})
		}
	}

	return filesByHash, nil
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

	filesByHash, err := getHashes(filesBySize)
	if err != nil {
		log.Fatal(err)
	}

	prune(filesByHash)

	fmt.Println(filesByHash)
}

