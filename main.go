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

func filterByFirstBytes(filesBySize map[int64][]File, n int64) (map[string][]File, error) {
	filesByFirstBytesHash := map[string][]File{}

	for _, sliceOfFile := range filesBySize {
		for _, file := range sliceOfFile {
			f, err := os.Open(file.Path)
			if err != nil {
				log.Printf("Encountered error %w while trying to open file %s.\n", err, file.Path)
				continue
			}

			h := md5.New()

			_, err = io.CopyN(h, f, n)
			if err != nil {
				if err != io.EOF {
					log.Printf("Encountered error %w when hashing file %s.\n", err, file)
					continue
				}
			}

			hash := fmt.Sprintf("%x", h.Sum(nil))
			file.FirstBytesHash = hash
			filesByFirstBytesHash[hash] = append(filesByFirstBytesHash[hash], file)

			f.Close()
		}
	}

	return filesByFirstBytesHash, nil
}

func prune[K comparable](filesByProperty map[K][]File) {
	for key, files := range filesByProperty {
		if len(files) < 2 {
			delete(filesByProperty, key)
		}
	}
}

func getHashes(filesByFirstBytesHash map[string][]File) (map[string][]File, error) {
	filesByHash := make(map[string][]File)

	for _, files := range filesByFirstBytesHash {
		for _, file := range files {
			h := md5.New()

			f, err := os.Open(file.Path)
			if err != nil {
				return nil, err
			}

			_, err = io.Copy(h, f)
			if err != nil {
				if err != io.EOF {
					log.Printf("Encountered error %w when hashing file %s.\n", err, file)
				}

				break
			}

			hash := fmt.Sprintf("%x", h.Sum(nil))
			file.Hash = hash
			filesByHash[hash] = append(filesByHash[hash], file)

			f.Close()
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

	fmt.Println("Duplicate groups after size filtering:")
	for size, sliceOfFile := range filesBySize {
		fmt.Printf("	%d bytes\n", size)
		for _, file := range sliceOfFile {
			fmt.Println("	-", file.Path)
		}
		fmt.Println()
	}

	// Filter by first bytes.
	filesByFirstBytesHash, err := filterByFirstBytes(filesBySize, 8)
	if err != nil {
		log.Fatal(err)
	}

	prune(filesByFirstBytesHash)

	fmt.Println("Duplicate groups after first-bytes filtering:")
	for hash, sliceOfFile := range filesByFirstBytesHash {
		fmt.Printf("	%s\n", hash)
		for _, file := range sliceOfFile {
			fmt.Println("	-", file.Path)
		}
		fmt.Println()
	}

	filesByHash, err := getHashes(filesByFirstBytesHash)
	if err != nil {
		log.Fatal(err)
	}

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

