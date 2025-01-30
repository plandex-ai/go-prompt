package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/plandex-ai/go-prompt"
)

func main() {
	p := prompt.New(
		executor,
		prompt.WithPrefix(">>> "),
		prompt.WithExecuteOnEnterCallback(ExecuteOnEnter),
	)

	p.Run()
}

func ExecuteOnEnter(p *prompt.Prompt, indentSize int) (int, bool) {
	input := p.Buffer().Text()
	lines := strings.SplitAfter(input, "\n")
	var spaces int
	if len(lines) > 0 {
		lastLine := lines[len(lines)-1]
		for _, char := range lastLine {
			if char == '}' {
				spaces -= 2 * indentSize
				break
			}
			if char != ' ' {
				break
			}
			spaces++
		}
	}

	char, _ := utf8.DecodeLastRuneInString(input)
	return 1 + spaces/indentSize, char == '}' && strings.Count(input, "}") == strings.Count(input, "{")
}

func executor(s string) {
	fmt.Println("Your input: " + s)
}
