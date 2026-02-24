package context

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/simonbs/wut/src/git"
)

type Context struct {
	RepoRoot string
}

func Create() (*Context, error) {
	repoRoot, err := git.GetRepoRoot("")
	if err != nil {
		return nil, fmt.Errorf("Not inside a Git repository. Run wut from the repo you want worktrees for.")
	}
	return &Context{RepoRoot: repoRoot}, nil
}

func getGlobalGitignorePath() string {
	cmd := []string{"config", "--global", "core.excludesfile"}
	output, err := git.Run(cmd, ".")
	if err != nil || output == "" {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".gitignore_global")
	}

	if strings.HasPrefix(output, "~") {
		homeDir, _ := os.UserHomeDir()
		output = filepath.Join(homeDir, output[1:])
	}
	return output
}

func CheckGlobalGitignore() bool {
	path := getGlobalGitignorePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == ".worktrees" {
			return true
		}
	}
	return false
}

func AddToGlobalGitignore() error {
	path := getGlobalGitignorePath()

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\n# Added by wut (https://github.com/simonbs/wut)\n.worktrees\n"); err != nil {
		return err
	}

	// Ensure git config points to this file
	_, _ = git.Run([]string{"config", "--global", "core.excludesfile", path}, ".")
	return nil
}

func checkLocalGitignore(repoRoot string) bool {
	// Check if .worktrees is in the repo's .gitignore
	gitignorePath := filepath.Join(repoRoot, ".gitignore")
	data, err := os.ReadFile(gitignorePath)
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == ".worktrees" {
			return true
		}
	}
	return false
}

func EnsureGitignoreConfigured(repoRoot string) error {
	worktreesDir := git.GetWorktreesDir(repoRoot)
	legacyDir := git.LegacyWorktreesDir(repoRoot)

	worktreesAbs, _ := filepath.Abs(worktreesDir)
	legacyAbs, _ := filepath.Abs(legacyDir)

	// Only legacy in-repo worktrees require .worktrees ignore configuration.
	if worktreesAbs != legacyAbs {
		return nil
	}

	// Already in global gitignore
	if CheckGlobalGitignore() {
		return nil
	}

	// Already in local .gitignore, no need to add globally
	if checkLocalGitignore(repoRoot) {
		return nil
	}

	// If running through wrapper, auto-add without prompting
	// (prompts don't work well when output is captured)
	if IsWrapperActive() {
		AddToGlobalGitignore()
		return nil
	}

	fmt.Println("\n⚠️  .worktrees is not in your global gitignore")
	fmt.Println("   This directory stores your worktrees and should be ignored by Git.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Add .worktrees to global gitignore? (Y/n): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "n" {
		if err := AddToGlobalGitignore(); err != nil {
			fmt.Println("❌ Failed to add .worktrees to global gitignore")
			fmt.Println()
			fmt.Println("   You can manually add it by running:")
			fmt.Println("   echo '.worktrees' >> ~/.gitignore_global")
			fmt.Println("   git config --global core.excludesfile ~/.gitignore_global")
			fmt.Println()
		} else {
			fmt.Println("✅ Added .worktrees to global gitignore")
			fmt.Println()
		}
	} else {
		fmt.Println()
		fmt.Println("⚠️  Remember to add .worktrees to your global gitignore:")
		fmt.Println("   echo '.worktrees' >> ~/.gitignore_global")
		fmt.Println("   git config --global core.excludesfile ~/.gitignore_global")
		fmt.Println()
	}

	return nil
}

func IsWrapperActive() bool {
	return os.Getenv("WUT_WRAPPER_ACTIVE") == "1"
}

func RequireWrapper(commandName string) {
	if !IsWrapperActive() {
		fmt.Fprintf(os.Stderr, "⚠️  The '%s' command requires shell integration to change directories.\n", commandName)
		fmt.Fprintln(os.Stderr, "   Add this to your ~/.zshrc or ~/.bashrc:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "     eval \"$(wut init)\"")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "   Then reload your shell with: source ~/.zshrc")
		os.Exit(1)
	}
}
