package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/worktree"
)

func cmdList(args []string) {
	ctx, err := context.Create()
	if err != nil {
		fail(err.Error())
	}

	entries, err := worktree.ParseList(ctx.RepoRoot)
	if err != nil {
		fail(err.Error())
	}

	if len(entries) == 0 {
		fmt.Println("No worktrees.")
		maybeAutoGc(ctx)
		return
	}

	maxLen := 0
	for _, e := range entries {
		name := e.BranchName
		if name == "" {
			name = "(detached)"
		}
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	cwd, _ := os.Getwd()

	// Find which worktree we're currently in (most specific match)
	currentWorktree := ""
	for _, e := range entries {
		if e.Path == cwd || strings.HasPrefix(cwd, e.Path+"/") {
			if len(e.Path) > len(currentWorktree) {
				currentWorktree = e.Path
			}
		}
	}

	for _, e := range entries {
		name := e.BranchName
		if name == "" {
			name = "(detached)"
		}

		// Icon: current worktree gets a pointer, others get a worktree icon
		icon := "🌿"
		if e.Path == currentWorktree {
			icon = "👉"
		} else if e.Path == ctx.RepoRoot {
			icon = "🏠"
		}

		fmt.Printf("%s %-*s  %s\n", icon, maxLen, name, tildify(e.Path))
	}

	maybeAutoGc(ctx)
}
