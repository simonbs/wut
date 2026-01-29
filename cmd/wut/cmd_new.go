package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/git"
	"github.com/simonbs/wut/src/worktree"
)

func cmdNew(args []string) {
	context.RequireWrapper("new")

	if len(args) < 1 {
		fail("Usage: wut new <branch> [--from <ref>]")
	}

	branch := args[0]
	fromRef := "HEAD"

	for i := 1; i < len(args); i++ {
		if args[i] == "--from" && i+1 < len(args) {
			fromRef = args[i+1]
			i++
		}
	}

	if err := context.EnsureGitignoreConfigured(); err != nil {
		fail(err.Error())
	}

	ctx, err := context.Create()
	if err != nil {
		fail(err.Error())
	}

	entries, _ := worktree.ParseList(ctx.RepoRoot)
	if existing := worktree.FindByBranch(entries, branch); existing != nil {
		fail(fmt.Sprintf("Branch '%s' already has a worktree at %s", branch, existing.Path))
	}

	worktreesDir := git.GetWorktreesDir(ctx.RepoRoot)
	relativePath := worktree.BranchToRelativePath(branch)
	basePath := filepath.Join(worktreesDir, relativePath)
	worktreePath := worktree.UniquePath(basePath)

	if err := os.MkdirAll(filepath.Dir(worktreePath), 0755); err != nil {
		fail(err.Error())
	}

	branchRef := "refs/heads/" + branch
	remoteRef := "refs/remotes/origin/" + branch

	var gitArgs []string
	if git.RefExists(ctx.RepoRoot, branchRef) {
		gitArgs = []string{"worktree", "add", worktreePath, branch}
	} else if git.RefExists(ctx.RepoRoot, remoteRef) {
		gitArgs = []string{"worktree", "add", "-b", branch, worktreePath, "origin/" + branch}
	} else {
		gitArgs = []string{"worktree", "add", "-b", branch, worktreePath, fromRef}
	}

	if _, err := git.Run(gitArgs, ctx.RepoRoot); err != nil {
		fail(err.Error())
	}

	fmt.Printf("__WUT_CD__:%s\n", worktreePath)
	maybeAutoGc(ctx)
}
