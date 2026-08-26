package repl

import (
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

	for {
		fmt.Print("\ndupe-sleuth > ")

		scanner.Scan()
		text := scanner.Text()
		tokens := cleanInput(text)

		if len(tokens) == 0 {
			continue
		}
		
		cmdName := tokens[0]

		cmd, ok := GetCommands()[cmdName]
		if !ok {
			fmt.Printf("Invalid command: %s\n", cmdName)
			continue
		}

		cmd.callback()
	}
}

