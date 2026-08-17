package main

import (
	"crypto/md5"
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

func gatherFiles(dir string) (map[int64][]File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filesBySize := map[int64][]File{}

	for _, entry := range entries {
		if entry.IsDir() {
			entryFilesBySize, err := gatherFiles(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}

			for size, sliceOfFile := range entryFilesBySize {
				filesBySize[size] = append(filesBySize[size], sliceOfFile...)
			}

		} else {
			info, err := entry.Info()
			if err != nil {
				return nil, err
			}
			
			size := info.Size()

			filesBySize[size] = append(filesBySize[size], File{
				Path: filepath.Join(dir, entry.Name()),
				Size: size,
			})
		}
	}

	return filesBySize, nil
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
	var dir string = Parse()

	filesBySize, err := gatherFiles(dir)
	if err != nil {
		log.Fatal(err)
	}

	prune(filesBySize)

	filesByHash, err := getHashes(filesBySize)
	if err != nil {
		log.Fatal(err)
	}

	prune(filesByHash)

	for _, sliceOfFile := range filesBySize {
		fmt.Println("Duplicate group:")
		for _, file := range sliceOfFile {
			fmt.Println("-", file.Path)
		}
	}
}

