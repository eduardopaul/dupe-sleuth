package app

import (
	"path/filepath"
	"log"
	"os"
	"sync"
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

func filterByHashes(filesByFirstBytesHash map[string][]File, concurrent bool) map[string][]File {
	filesByHash := make(map[string][]File)

	files := []File{}
	for _, sliceOfFiles := range filesByFirstBytesHash {
		for _, file := range sliceOfFiles {
			files = append(files, file)
		}
	}

	if concurrent {
		var wg sync.WaitGroup

		for _, file := range files{
			wg.Go(func() {
				hash, err := hashFile(file)
				if err != nil {
					log.Printf("Encountered error %w while hashing file %s.\n", err, file.Path)
					return
				}

				file.Hash = hash
				filesByHash[hash] = append(filesByHash[hash], file)
			})
		}

		wg.Wait()

	} else {
		for _, file := range files{
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

func Sleuth(dir string, logging bool, concurrent bool) (map[string][]File, error) {
	filesBySize, err := filterBySizes(dir)
	if err != nil {
		return nil, err
	}

	prune(filesBySize)

	if logging {
		logSize(filesBySize)
	}

	filesByFirstBytes := filterByFirstBytes(filesBySize, 8)
	prune(filesByFirstBytes)

	if logging {
		logBytes(filesByFirstBytes)
	}

	filesByHash := filterByHashes(filesByFirstBytes, concurrent)

	prune(filesByHash)

	return filesByHash, nil
}

