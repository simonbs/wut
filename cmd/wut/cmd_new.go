package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/git"
	"github.com/simonbs/wut/src/worktree"
)

func cmdNew(args []string) {
	context.RequireWrapper("new")

	var branch string
	fromRef := "HEAD"

	// Parse args: first non-flag arg is the branch name
	positional := []string{}
	for i := 0; i < len(args); i++ {
		if args[i] == "--from" && i+1 < len(args) {
			fromRef = args[i+1]
			i++
		} else {
			positional = append(positional, args[i])
		}
	}

	ctx, err := context.Create()
	if err != nil {
		fail(err.Error())
	}

	entries, _ := worktree.ParseList(ctx.RepoRoot)

	if len(positional) > 0 {
		branch = positional[0]
		if existing := worktree.FindByBranch(entries, branch); existing != nil {
			fail(fmt.Sprintf("Branch '%s' already has a worktree at %s", branch, existing.Path))
		}
	} else {
		// Generate a random branch name
		branch = petname.Generate(2, "-") + "-" + time.Now().Format("2006-01-02")
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
}
