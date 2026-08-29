package main

import (
	"dupe-sleuth/app"
	"dupe-sleuth/repl"
	"log"
)

func main() {
	opt := app.Options
	if *opt.Interactive {
		repl.Run()
	} else {
		dupeFiles, err := app.Sleuth(opt.Dir, *opt.Logging, *opt.Concurrent)
		if err != nil {
			log.Fatal(err)
		}

		app.PrintGroups(dupeFiles)
	}
}

