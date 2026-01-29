package main

import "fmt"

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
	fmt.Println("Ephemeral worktrees that stay out of your vibe zone.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wut new <branch> [--from ref] 🌱 Create a new worktree")
	fmt.Println("  wut list                      📋 List worktrees")
	fmt.Println("  wut go [branch]               🚀 Navigate to a worktree")
	fmt.Println("  wut path <branch>             📂 Print worktree path")
	fmt.Println("  wut rm <branch> [--force]     🗑  Remove a worktree")
	fmt.Println("  wut gc [--dry-run]            🧹 Remove orphaned worktrees")
	fmt.Println()
	fmt.Println("Add shell integration to ~/.zshrc or ~/.bashrc:")
	fmt.Println("  eval \"$(wut init)\"")
}
