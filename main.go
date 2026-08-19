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

func getFirstBytes(file File, n int64) ([]byte, error) {
	f, err := os.Open(file.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b := make([]byte, min(file.Size, n))
	_, err = io.ReadFull(f, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func hashFile(file File) (string, error) {
		h := md5.New()

		f, err := os.Open(file.Path)
		if err != nil {
			return "", err
		}
		defer f.Close()

		_, err = io.Copy(h, f)
		if err != nil {
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

func main() {
	var dir string = Parse()

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

	filesByHash := filterByHashes(filesByFirstBytes)

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

