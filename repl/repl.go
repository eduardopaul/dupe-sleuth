package repl

import (
	"dupe-sleuth/app"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(input string) []string {
	lowered := strings.ToLower(input)
	split := strings.Fields(lowered)
	return split
}

func Run() {
	scanner := bufio.NewScanner(os.Stdin)

	appStruct := app.AppType{
		Duplicates: map[string][]app.File{},
	}

	for {
		fmt.Print("\ndupe-sleuth > ")

		scanner.Scan()
		text := scanner.Text()
		tokens := cleanInput(text)

		if len(tokens) == 0 {
			continue
		}
		
		cmdName := tokens[0]
		if cmdName == "aid" {
			Aid()
			continue
		}

		cmd, ok := commands[cmdName]
		if !ok {
			fmt.Printf("Invalid command: %s\n", cmdName)
			Aid()
			continue
		}

		args := tokens[1:]

		var err error
		appStruct, err = cmd.callback(appStruct, args)
		if err != nil{
			fmt.Println(err)
		}
	}
}

