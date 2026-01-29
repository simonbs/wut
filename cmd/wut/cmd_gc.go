package main

import (
	"fmt"
	"os"

	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/worktree"
)

func cmdGc(args []string) {
	dryRun := false
	for _, arg := range args {
		if arg == "--dry-run" {
			dryRun = true
		}
	}

	if err := context.EnsureGitignoreConfigured(); err != nil {
		fail(err.Error())
	}

	ctx, err := context.Create()
	if err != nil {
		fail(err.Error())
	}

	stale, err := worktree.GetStalePaths(ctx.RepoRoot)
	if err != nil {
		fail(err.Error())
	}

	if len(stale) == 0 {
		fmt.Println("No stale worktrees found.")
		return
	}

	for _, path := range stale {
		if dryRun {
			fmt.Println(path)
		} else {
			os.RemoveAll(path)
			fmt.Println(path)
		}
	}
}
