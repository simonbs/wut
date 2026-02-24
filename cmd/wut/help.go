package main

import (
	"fmt"

	"github.com/simonbs/wut/src/context"
)

func printUsage() {
	// ANSI color codes
	purple := "\033[35m"
	reset := "\033[0m"

	fmt.Printf(`%s
                                ▄▄▄▄▄   
                       ██      █▀▀▀▀██  
██      ██ ██    ██  ███████       ▄█▀  
▀█  ██  █▀ ██    ██    ██        ▄██▀   
 ██▄██▄██  ██    ██    ██        ██     
 ▀██  ██▀  ██▄▄▄███    ██▄▄▄     ▄▄     
  ▀▀  ▀▀    ▀▀▀▀ ▀▀     ▀▀▀▀     ▀▀     
%s`, purple, reset)
	fmt.Println()
	fmt.Println("Ephemeral worktrees without the ceremony.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wut new [branch] [--from ref] 🌱 Create a new worktree")
	fmt.Println("  wut mv [old] <new-name>       ✏️ Rename a worktree's branch")
	fmt.Println("  wut list                      📋 List worktrees")
	fmt.Println("  wut go [branch]               🚀 Navigate to a worktree")
	fmt.Println("  wut path <branch>             📂 Print worktree path")
	fmt.Println("  wut rm <branch> [--force]     🗑  Remove a worktree")

	if !context.IsWrapperActive() {
		fmt.Println()
		fmt.Println("⚠️  Add shell integration to ~/.zshrc or ~/.bashrc:")
		fmt.Println()
		fmt.Println("     eval \"$(wut init)\"")
	}
}
