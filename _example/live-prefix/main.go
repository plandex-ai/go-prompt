package main

import (
	"fmt"

	prompt "github.com/elk-language/go-prompt"
	istrings "github.com/elk-language/go-prompt/strings"
)

var LivePrefix string = ">>> "

func executor(in string) {
	fmt.Println("Your input: " + in)
	if in == "" {
		LivePrefix = ">>> "
		return
	}
	LivePrefix = in + "> "
}

func completer(in prompt.Document) ([]prompt.Suggest, istrings.RuneNumber, istrings.RuneNumber) {
	endIndex := in.CurrentRuneIndex()
	w := in.GetWordBeforeCursor()
	startIndex := endIndex - istrings.RuneCountInString(w)

	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
	}
	return prompt.FilterHasPrefix(s, w, true), startIndex, endIndex
}

func changeLivePrefix() string {
	return LivePrefix
}

func main() {
	p := prompt.New(
		executor,
		prompt.WithPrefixCallback(changeLivePrefix),
		prompt.WithTitle("live-prefix-example"),
		prompt.WithCompleter(completer),
	)
	p.Run()
}
