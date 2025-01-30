package main

import (
	"fmt"
	"os"
	"os/exec"

	prompt "github.com/plandex-ai/go-prompt"
)

func executor(t string) {
	if t != "bash" {
		fmt.Println("Sorry, I don't understand.")
		return
	}

	cmd := exec.Command("bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	p := prompt.New(
		executor,
	)
	p.Run()
}
