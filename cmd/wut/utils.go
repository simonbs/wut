package main

import (
	"fmt"
	"os"
	"strings"
)

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func tildify(path string) string {
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "fish") {
		return "fish"
	}
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	return "bash"
}
