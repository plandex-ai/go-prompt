package main

import (
	"fmt"
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
	char, _ := utf8.DecodeLastRuneInString(input)
	return 0, char == '!'
}

func executor(s string) {
	fmt.Println("Your input: " + s)
}
