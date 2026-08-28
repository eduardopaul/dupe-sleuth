package main

import (
	"dupe-sleuth/repl"
	"log"
)

func main() {
	if *Options.Interactive {
		repl.Run()
	} else {
		dupeFiles, err := sleuth(Options.Dir, *Options.Logging, *Options.Concurrent)
		if err != nil {
			log.Fatal(err)
		}

		printGroups(dupeFiles)
	}
}

