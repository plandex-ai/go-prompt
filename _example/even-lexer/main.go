package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/go-prompt"
	"github.com/elk-language/go-prompt/strings"
)

func main() {
	p := prompt.New(
		executor,
		prompt.WithLexer(prompt.NewEagerLexer(wordLexer)),
		prompt.WithLexer(prompt.NewEagerLexer(charLexer)), // the last one overrides the other
	)

	p.Run()
}

// colors every other character green
func charLexer(line string) []prompt.Token {
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

// colors every other word green
func wordLexer(line string) []prompt.Token {
	if len(line) == 0 {
		return nil
	}

	var elements []prompt.Token
	var currentByte strings.ByteNumber
	var wordIndex int
	var lastChar rune

	var color prompt.Color
	for i, char := range line {
		currentByte = strings.ByteNumber(i)
		if unicode.IsSpace(char) {
			if wordIndex%2 == 0 {
				color = prompt.Green
			} else {
				color = prompt.White
			}

			element := prompt.NewSimpleToken(color, currentByte)
			elements = append(elements, element)
			wordIndex++
			continue
		}
		lastChar = char
	}
	if !unicode.IsSpace(lastChar) {
		if wordIndex%2 == 0 {
			color = prompt.Green
		} else {
			color = prompt.White
		}
		element := prompt.NewSimpleToken(color, currentByte)
		elements = append(elements, element)
	}

	return elements
}

func executor(s string) {
	fmt.Println("Your input: " + s)
}
