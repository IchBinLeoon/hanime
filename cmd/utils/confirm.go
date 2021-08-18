package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskForConfirmation(question string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [Y/N] ", question)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		res := strings.ToLower(strings.TrimSpace(input))
		if res == "y" || res == "yes" {
			return true, nil
		} else if res == "n" || res == "no" {
			break
		} else {
			continue
		}
	}
	return false, nil
}
