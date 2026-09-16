package repl

import (
	"bufio"
	"dupe-sleuth/app"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Command struct {
	description string
	callback func(app.State, []string) (app.State, error)
}

var opt = app.Options

func choose(scanner *bufio.Scanner, files []app.File) int {
	for idxFile, file := range files {
		fmt.Println("  ", idxFile, file.Path)
	}

	var (
		choice int
		err error
	)

	for {
		fmt.Printf("Choose which file to keep [0-%d]: ", len(files)-1)
		scanner.Scan()
		textChoice := scanner.Text()
		choice, err = strconv.Atoi(textChoice)
		if err != nil || choice < 0 || choice >= len(files) {
			fmt.Print("Invalid choice. ")
			continue
		}
		break
	}

	return choice
}

var commands = map[string]Command{
	"flee": {
		description: "Exit dupe-sleuth, letting go of any changes.",
		callback: func(appState app.State, args []string) (app.State, error) {
			os.Exit(0)
			return appState, nil
		},
	},
	"sleuth": {
		description: "Find duplicate files.",
		callback: func(appState app.State, args []string) (app.State, error) {
			var err error
			newDupeFiles, err := app.Sleuth(args[0], *opt.Logging, *opt.Concurrent)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return appState, fmt.Errorf(`Value "%s" does not exist or is not a valid directory.`, args[0])
				}
				return appState, err
			}

			for hash, newSliceOfFile := range newDupeFiles {
				oldSliceOfFile, exists := appState.Duplicates[hash]

				if exists {
					for _, file := range newSliceOfFile {
						if !slices.Contains(oldSliceOfFile, file) {
							appState.Duplicates[hash] = append(appState.Duplicates[hash], file)
						}
					}
				} else {
					appState.Duplicates[hash] = newSliceOfFile
					appState.Order = append(appState.Order, hash)
				}
			}

			return appState, err
		},
	},
	"unveil": {
		description: "Show the duplicate files that have already been found.",
		callback: func(appState app.State, args []string) (app.State, error) {
			app.PrintGroups(appState.Duplicates)
			return appState, nil
		},
	},
	"stamp": {
		description: "Mark file to keep.",
		callback: func(appState app.State, args []string) (app.State, error) {
			scanner := bufio.NewScanner(os.Stdin)

			for idx, hash := range appState.Order {
				fmt.Printf("Group %d\n", idx)

				choice := choose(scanner, appState.Duplicates[hash])
				
				appState.Marked[hash] = appState.Duplicates[hash][choice]

				fmt.Printf("Selected file %d\n", choice)

				fmt.Println()
			}

			return appState, nil
		},
	},
	"efface": {
		description: "Erase non-marked files.",
		callback: func(appState app.State, args []string) (app.State, error) {
			for hash, markedFile := range appState.Marked {
				fmt.Println(hash, markedFile)
				for idx, file := range appState.Duplicates[hash] {
					fmt.Println("---", idx, file, markedFile == file)
					if markedFile != file {
						err := os.Remove(file.Path)
						if err != nil {
							return appState, err
						}
					}
				}
			}

			return appState, nil
		},
	},
	"whereabouts": {
		description: "Print the current working directory.",
		callback: func(appState app.State, args []string) (app.State, error) {
			wd, err := os.Getwd()
			if err != nil {
				return appState, err
			}

			fmt.Println(wd)

			return appState, nil
		},
	},
	"wander": {
		description: "Change the current working directory.",
		callback: func(appState app.State, args []string) (app.State, error) {
			err := os.Chdir(args[0])
			if err != nil {
				return appState, err
			}
			return appState, nil
		},
	},
	"catalog": {
		description: "List the current working directory content.",
		callback: func(appState app.State, args []string) (app.State, error) {
			wd, err := os.Getwd()
			if err != nil {
				return appState, err
			}

			list, err := os.ReadDir(wd)
			if err != nil {
				return appState, err
			}

			for _, item := range list {
				fmt.Print(item, "        ")
			}

			return appState, nil
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

