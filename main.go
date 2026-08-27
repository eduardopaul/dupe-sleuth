package main

import (
	"dupe-sleuth/repl"
	"log"
)

func main() {
	options := Parse()

	var logging *bool = options.Log
	var dir string = options.Dir
	var concurrent *bool = options.Concurrent
	var interactive *bool = options.Interactive

	if *interactive {
		repl.Run()
	} else {
		dupeFiles, err := sleuth(dir, *logging, *concurrent)
		if err != nil {
			log.Fatal(err)
		}

		printGroups(dupeFiles)
	}
}

