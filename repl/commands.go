package repl

import (
	"dupe-sleuth/app"
	"errors"
	"fmt"
	"os"
	"slices"
)

type Command struct {
	description string
	callback func([]string) error
}

var opt = app.Options

var DupeFiles = map[string][]app.File{}

var commands = map[string]Command{
	"flee": {
		description: "Exit dupe-sleuth, letting go of any changes.",
		callback: func(args []string) error {
			os.Exit(0)
			return nil
		},
	},
	"sleuth": {
		description: "Find duplicate files.",
		callback: func(args []string) error {
			var err error
			newDupeFiles, err := app.Sleuth(args[0], *opt.Logging, *opt.Concurrent)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf(`Value "%s" does not exist or is not a valid directory.`, args[0])
				}
				return err
			}

			for hash, newSliceOfFile := range newDupeFiles {
				oldSliceOfFile, exists := DupeFiles[hash]

				if exists {
					for _, file := range newSliceOfFile {
						if !slices.Contains(oldSliceOfFile, file) {
							DupeFiles[hash] = append(DupeFiles[hash], file)
						}
					}
				} else {
					DupeFiles[hash] = newSliceOfFile
				}
			}

			return err
		},
	},
	"unveil": {
		description: "Show the duplicate files that have already been found.",
		callback: func(args []string) error {
			app.PrintGroups(DupeFiles)
			return nil
		},
	},
	"stamp": {
		description: "Mark file to receive action.",
		callback: func(args []string) error {
			return nil
		},
	},
	"efface": {
		description: "Erase marked files.",
		callback: func(args []string) error {
			return nil
		},
	},
	"whereabouts": {
		description: "Print the current working directory.",
		callback: func(args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			fmt.Println(wd)

			return nil
		},
	},
	"wander": {
		description: "Change the current working directory.",
		callback: func(args []string) error {
			err := os.Chdir(args[0])
			if err != nil {
				return err
			}
			return nil
		},
	},
}

func Aid() error {
	fmt.Println("List of available commands:")
	fmt.Println()
	fmt.Println("aid - Print this message.")

	for name, cmd := range commands {
		fmt.Printf("%s - %s\n", name, cmd.description)
	}

	return nil
}

