package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/worktree"
)

func cmdGo(args []string) {
	var target string

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			target = arg
		}
	}

	context.RequireWrapper("go")

	if err := context.EnsureGitignoreConfigured(); err != nil {
		fail(err.Error())
	}

	ctx, err := context.Create()
	if err != nil {
		fail(err.Error())
	}

	entries, _ := worktree.ParseList(ctx.RepoRoot)

	var resolvedPath string
	if target == "" {
		// Go to repo root
		for _, e := range entries {
			if e.Path == ctx.RepoRoot {
				resolvedPath = e.Path
				break
			}
		}
		if resolvedPath == "" && len(entries) > 0 {
			resolvedPath = ctx.RepoRoot
		}
	} else {
		if e := worktree.FindByBranch(entries, target); e != nil {
			resolvedPath = e.Path
		} else {
			absTarget, _ := filepath.Abs(target)
			if e := worktree.FindByPath(entries, absTarget); e != nil {
				resolvedPath = e.Path
			}
		}
	}

	if resolvedPath == "" {
		fail(fmt.Sprintf("No worktree found for '%s'.", target))
	}

	fmt.Printf("__WUT_CD__:%s\n", resolvedPath)
	maybeAutoGc(ctx)
}

func cmdPath(args []string) {
	if len(args) < 1 {
		fail("Usage: wut path <branch>")
	}

	branch := args[0]

	ctx, err := context.Create()
	if err != nil {
		fail(err.Error())
	}

	entries, _ := worktree.ParseList(ctx.RepoRoot)
	entry := worktree.FindByBranch(entries, branch)
	if entry == nil {
		fail(fmt.Sprintf("No worktree found for branch '%s'.", branch))
	}

	fmt.Println(entry.Path)
}
