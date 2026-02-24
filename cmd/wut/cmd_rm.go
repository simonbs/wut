package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/git"
	"github.com/simonbs/wut/src/worktree"
)

func cmdRm(args []string) {
	if len(args) < 1 {
		fail("Usage: wut rm <branch-or-path> [--force]")
	}

	target := args[0]
	force := false
	for _, arg := range args[1:] {
		if arg == "--force" {
			force = true
		}
	}

	ctx, err := context.Create()
	if err != nil {
		fail(err.Error())
	}

	entries, _ := worktree.ParseList(ctx.RepoRoot)

	var resolvedPath string
	var branchName string

	if e := worktree.FindByBranch(entries, target); e != nil {
		resolvedPath = e.Path
		branchName = e.BranchName
	} else {
		absTarget, _ := filepath.Abs(target)
		if e := worktree.FindByPath(entries, absTarget); e != nil {
			resolvedPath = e.Path
			branchName = e.BranchName
		}
	}

	if resolvedPath == "" {
		fail(fmt.Sprintf("No worktree found for '%s'.", target))
	}

	cwd, _ := os.Getwd()
	resolvedAbs, _ := filepath.Abs(resolvedPath)
	isInWorktree := strings.HasPrefix(cwd, resolvedAbs)

	if isInWorktree && !force {
		fail("Cannot remove current worktree. Use --force to remove and return to repo root.")
	}

	if isInWorktree && !context.IsWrapperActive() {
		context.RequireWrapper("rm --force")
	}

	// Safety checks
	if !force {
		output, err := git.Run([]string{"status", "--porcelain"}, resolvedPath)
		if err == nil && strings.TrimSpace(output) != "" {
			fail("Worktree has uncommitted changes. Use --force to remove anyway.")
		}

		if branchName != "" {
			merged, _ := git.Run([]string{"branch", "--merged"}, ctx.RepoRoot)
			isMerged := false
			for _, line := range strings.Split(merged, "\n") {
				// Strip leading markers like "* ", "+ ", "  "
				line = strings.TrimSpace(line)
				line = strings.TrimPrefix(line, "* ")
				line = strings.TrimPrefix(line, "+ ")
				line = strings.TrimSpace(line)
				if line == branchName {
					isMerged = true
					break
				}
			}
			if !isMerged {
				fail(fmt.Sprintf("Branch '%s' is not fully merged. Use --force to remove anyway.", branchName))
			}
		}
	}

	gitArgs := []string{"worktree", "remove"}
	if force {
		gitArgs = append(gitArgs, "--force")
	}
	gitArgs = append(gitArgs, resolvedPath)

	if _, err := git.Run(gitArgs, ctx.RepoRoot); err != nil {
		fail(err.Error())
	}

	// Delete branch
	if branchName != "" {
		deleteFlag := "-d"
		if force {
			deleteFlag = "-D"
		}
		git.Run([]string{"branch", deleteFlag, branchName}, ctx.RepoRoot)
	}

	if isInWorktree && cwd != ctx.RepoRoot {
		fmt.Printf("__WUT_CD__:%s\n", ctx.RepoRoot)
	}

}
