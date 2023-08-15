package main

import (
	"fmt"

	prompt "github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
)

func executor(in string) {
	fmt.Println("Your input: " + in)
}

func completer(in prompt.Document) ([]prompt.Suggest, pstrings.RuneNumber, pstrings.RuneNumber) {
	s := []prompt.Suggest{
		{Text: "こんにちは", Description: "'こんにちは' means 'Hello' in Japanese"},
		{Text: "감사합니다", Description: "'안녕하세요' means 'Hello' in Korean."},
		{Text: "您好", Description: "'您好' means 'Hello' in Chinese."},
		{Text: "Добрый день", Description: "'Добрый день' means 'Hello' in Russian."},
	}
	endIndex := in.CurrentRuneIndex()
	w := in.GetWordBeforeCursor()
	startIndex := endIndex - pstrings.RuneCountInString(w)

	return prompt.FilterHasPrefix(s, w, true), startIndex, endIndex
}

func main() {
	p := prompt.New(
		executor,
		prompt.WithCompleter(completer),
		prompt.WithPrefix(">>> "),
		prompt.WithTitle("sql-prompt for multi width characters"),
	)
	p.Run()
}
