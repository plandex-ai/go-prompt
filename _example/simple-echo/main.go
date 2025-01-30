package main

import (
	"fmt"

	prompt "github.com/plandex-ai/go-prompt"
	pstrings "github.com/plandex-ai/go-prompt/strings"
)

func completer(in prompt.Document) ([]prompt.Suggest, pstrings.RuneNumber, pstrings.RuneNumber) {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
	}
	endIndex := in.CurrentRuneIndex()
	w := in.GetWordBeforeCursor()
	startIndex := endIndex - pstrings.RuneCountInString(w)

	return prompt.FilterHasPrefix(s, w, true), startIndex, endIndex
}

func main() {
	in := prompt.Input(
		prompt.WithPrefix(">>> "),
		prompt.WithTitle("sql-prompt"),
		prompt.WithHistory([]string{"SELECT * FROM users;"}),
		prompt.WithPrefixTextColor(prompt.Yellow),
		prompt.WithSelectedSuggestionBGColor(prompt.LightGray),
		prompt.WithSuggestionBGColor(prompt.DarkGray),
		prompt.WithCompleter(completer),
	)
	fmt.Println("Your input: " + in)
}
