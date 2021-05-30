package main

import (
	"fmt"

	term "github.com/oxide-one/systemd.go/internal/terminal"
)

func main() {
	terminal := term.NewTerminal(term.DefaultSettings())
	_ = terminal
	fmt.Println(terminal.Passwords)
	term.Display(terminal)
}
