package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ask asks a question and returns the answer
// It removes the newline from the answer
// It returns an empty string if no question is asked
//
// Example:
//
//	answer := Ask("What is your name? ")
//	fmt.Println("Your name is", answer)
func Ask(question string) string {

	if question == "" {
		return "" // no question to ask
	}

	fmt.Print(question) // ask the question
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n') // read the answer, ending with a newline
	return strings.TrimSpace(text)     // remove the newline
}
