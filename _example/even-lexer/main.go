package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/elk-language/go-prompt"
	"github.com/elk-language/go-prompt/strings"
)

func main() {
	p := prompt.New(
		executor,
		prompt.WithLexer(prompt.NewEagerLexer(lexer)),
	)

	p.Run()
}

func lexer(line string) []prompt.Token {
	var elements []prompt.Token

	for i, value := range line {
		var color prompt.Color
		// every even char must be green.
		if i%2 == 0 {
			color = prompt.Green
		} else {
			color = prompt.White
		}
		lastByteIndex := strings.ByteNumber(i + utf8.RuneLen(value) - 1)
		element := prompt.NewSimpleToken(color, lastByteIndex)

		elements = append(elements, element)
	}

	return elements
}

func executor(s string) {
	fmt.Println("You printed: " + s)
}
