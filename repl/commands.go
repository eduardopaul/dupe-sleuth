package repl

import (
	"fmt"
	"os"
)

type Command struct {
	description string
	callback func() error
}

var commands = map[string]Command{
	"flee": {
		description: "Exit dupe-sleuth, letting go of any changes.",
		callback: func() error {
			os.Exit(0)
			return nil
		},
	},
	"sleuth": {
		description: "Find duplicate files.",
		callback: func() error {
			return nil
		},
	},
	"unveil": {
		description: "Show the duplicate files that have already been found.",
		callback: func() error {
			return nil
		},
	},
	"stamp": {
		description: "Mark file to receive action.",
		callback: func() error {
			return nil
		},
	},
	"efface": {
		description: "Erase marked files.",
		callback: func() error {
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

