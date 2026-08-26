package repl

import (
	"os"
)

type commands struct {
	description string
	callback func() error
}

func GetCommands() map[string]commands {
	return map[string]commands{
		"exit": {
			description: "Exit dupe-sleuth.",
			callback: func() error {
				os.Exit(0)
				return nil
			},
		},
	}
}
