package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

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


