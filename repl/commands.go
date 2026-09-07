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
	callback func(app.AppType, []string) (app.AppType, error)
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
		callback: func(appStruct app.AppType, args []string) (app.AppType, error) {
			os.Exit(0)
			return appStruct, nil
		},
	},
	"sleuth": {
		description: "Find duplicate files.",
		callback: func(appStruct app.AppType, args []string) (app.AppType, error) {
			var err error
			newDupeFiles, err := app.Sleuth(args[0], *opt.Logging, *opt.Concurrent)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return appStruct, fmt.Errorf(`Value "%s" does not exist or is not a valid directory.`, args[0])
				}
				return appStruct, err
			}

			for hash, newSliceOfFile := range newDupeFiles {
				oldSliceOfFile, exists := appStruct.Duplicates[hash]

				if exists {
					for _, file := range newSliceOfFile {
						if !slices.Contains(oldSliceOfFile, file) {
							appStruct.Duplicates[hash] = append(appStruct.Duplicates[hash], file)
						}
					}
				} else {
					appStruct.Duplicates[hash] = newSliceOfFile
					appStruct.Order = append(appStruct.Order, hash)
				}
			}

			return appStruct, err
		},
	},
	"unveil": {
		description: "Show the duplicate files that have already been found.",
		callback: func(appStruct app.AppType, args []string) (app.AppType, error) {
			app.PrintGroups(appStruct.Duplicates)
			return appStruct, nil
		},
	},
	"stamp": {
		description: "Mark file to receive action.",
		callback: func(appStruct app.AppType, args []string) (app.AppType, error) {
			scanner := bufio.NewScanner(os.Stdin)
			for idx, hash := range appStruct.Order {
				fmt.Printf("Group %d\n", idx)

				choice := choose(scanner, appStruct.Duplicates[hash])
				
				fmt.Printf("Selected file %d\n", choice)

				fmt.Println()
			}

			return appStruct, nil
		},
	},
	"efface": {
		description: "Erase marked files.",
		callback: func(appStruct app.AppType, args []string) (app.AppType, error) {
			return appStruct, nil
		},
	},
	"whereabouts": {
		description: "Print the current working directory.",
		callback: func(appStruct app.AppType, args []string) (app.AppType, error) {
			wd, err := os.Getwd()
			if err != nil {
				return appStruct, err
			}

			fmt.Println(wd)

			return appStruct, nil
		},
	},
	"wander": {
		description: "Change the current working directory.",
		callback: func(appStruct app.AppType, args []string) (app.AppType, error) {
			err := os.Chdir(args[0])
			if err != nil {
				return appStruct, err
			}
			return appStruct, nil
		},
	},
	"catalog": {
		description: "List the current working directory content.",
		callback: func(appStruct app.AppType, args []string) (app.AppType, error) {
			wd, err := os.Getwd()
			if err != nil {
				return appStruct, err
			}

			list, err := os.ReadDir(wd)
			if err != nil {
				return appStruct, err
			}

			for _, item := range list {
				fmt.Print(item, "        ")
			}

			return appStruct, nil
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

