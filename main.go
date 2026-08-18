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
	FirstBytesHash string
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

func filterByFirstBytes(filesBySize map[int64][]File, n int64) map[string][]File {
	filesByFirstBytesHash := map[string][]File{}

	for _, sliceOfFile := range filesBySize {
		for _, file := range sliceOfFile {
			hash, err := processHash(file, n)
			if err != nil {
				log.Printf("Encountered error %w while handling file %s.\n", err, file.Path)
				continue
			}

			file.FirstBytesHash = hash
			filesByFirstBytesHash[hash] = append(filesByFirstBytesHash[hash], file)
		}
	}

	return filesByFirstBytesHash
}

func processHash(file File, n int64) (string, error) {
		h := md5.New()

		f, err := os.Open(file.Path)
		if err != nil {
			return "", err
		}
		defer f.Close()

		if n == 0 {
			_, err = io.Copy(h, f)
		} else {
			_, err = io.CopyN(h, f, n)
		}

		if err != nil && err != io.EOF {
			return "", err
		}

		hash := fmt.Sprintf("%x", h.Sum(nil))
		return hash, nil
}

func prune[K comparable](filesByProperty map[K][]File) {
	for key, files := range filesByProperty {
		if len(files) < 2 {
			delete(filesByProperty, key)
		}
	}
}

func getHashes(filesByFirstBytesHash map[string][]File) map[string][]File {
	filesByHash := make(map[string][]File)

	for _, files := range filesByFirstBytesHash {
		for _, file := range files {
			hash, err := processHash(file, 0)
			if err != nil {
				log.Printf("Encountered error %w while hashing file %s.\n", err, file.Path)
				break
			}

			file.Hash = hash
			filesByHash[hash] = append(filesByHash[hash], file)
		}
	}

	return filesByHash
}

func main() {
	var dir string = Parse()

	filesBySize, err := gatherFiles(dir)
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

	filesByFirstBytesHash := filterByFirstBytes(filesBySize, 8)
	prune(filesByFirstBytesHash)

	fmt.Println("Duplicate groups after first-bytes filtering:")
	for hash, sliceOfFile := range filesByFirstBytesHash {
		fmt.Printf("	%s\n", hash)
		for _, file := range sliceOfFile {
			fmt.Println("	-", file.Path)
		}
		fmt.Println()
	}

	filesByHash := getHashes(filesByFirstBytesHash)
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

