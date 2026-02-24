package main

import (
	"fmt"
	"os"
)

const version = "0.2.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "init":
		cmdInit()
	case "new":
		cmdNew(args)
	case "list":
		cmdList(args)
	case "go":
		cmdGo(args)
	case "path":
		cmdPath(args)
	case "rm":
		cmdRm(args)
	case "--completions":
		cmdCompletions(args)
	case "--help", "-h":
		printUsage()
	case "--version", "-V":
		fmt.Println("wut", version)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
