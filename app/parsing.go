package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	flag "github.com/spf13/pflag"
)

var concurrentFlag = flag.BoolP("concurrent", "c", false, "Process hashing concurrently.")
var logFlag = flag.Bool("log", false, "Enable logging.")
var interactiveFlag = flag.BoolP("interactive", "i", false, "Enter interactive REPL.")

type cliOptions struct {
	Dir string
	Concurrent *bool
	Logging *bool
	Interactive *bool
}

func Parse() cliOptions {
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

	return cliOptions{
		Dir: dir,
		Concurrent: concurrentFlag,
		Logging: logFlag,
		Interactive: interactiveFlag,
	}
}

var Options = Parse()

