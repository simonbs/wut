package context

import (
	"fmt"
	"os"

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
