package main

import (
	"path/filepath"
	"log"
	"os"
)

type File struct {
	Path string
	Size int64
	Hash string
	FirstBytes string
}

func filterBySizes(dir string) (map[int64][]File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filesBySize := map[int64][]File{}

	for _, entry := range entries {
		if entry.IsDir() {
			entryFilesBySize, err := filterBySizes(filepath.Join(dir, entry.Name()))
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

func filterByFirstBytes(filesBySize map[int64][]File, n int64) map[string][]File {
	filesByFirstBytes := map[string][]File{}

	for _, sliceOfFile := range filesBySize {
		for _, file := range sliceOfFile {
			firstBytes, err := getFirstBytes(file, n)
			if err != nil {
				log.Printf("Encountered error %w while handling file %s.\n", err, file.Path)
				continue
			}

			firstBytesString := string(firstBytes)
			file.FirstBytes = firstBytesString
			filesByFirstBytes[firstBytesString] = append(filesByFirstBytes[firstBytesString], file)
		}
	}

	return filesByFirstBytes
}

func filterByHashes(filesByFirstBytesHash map[string][]File) map[string][]File {
	filesByHash := make(map[string][]File)

	for _, files := range filesByFirstBytesHash {
		for _, file := range files {
			hash, err := hashFile(file)
			if err != nil {
				log.Printf("Encountered error %w while hashing file %s.\n", err, file.Path)
				continue
			}

			file.Hash = hash
			filesByHash[hash] = append(filesByHash[hash], file)
		}
	}

	return filesByHash
}
