package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	flag "github.com/spf13/pflag"
)

var concurrentFlag = flag.BoolP("concurrent", "c", false, "Process hashing concurrently.")

type Options struct {
	Dir string
	Concurrent *bool
}

func Parse() Options {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("You must provide the root path the search should performed in.")
		os.Exit(1)
	} else if flag.NArg() > 1 {
		fmt.Println("Too many arguments: you should provide just the root path the search should be performed in.")
		os.Exit(1)
	}

	var dir string = filepath.Clean(flag.Arg(0))

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

	return Options{
		Dir: dir,
		Concurrent: concurrentFlag,
	}
}

