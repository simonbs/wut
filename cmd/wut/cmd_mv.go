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

func cmdMv(args []string) {
	context.RequireWrapper("mv")

	if len(args) < 1 || len(args) > 2 {
		fail("Usage: wut mv [old-name] <new-name>")
	}

	ctx, err := context.Create()
	if err != nil {
		fail(err.Error())
	}

	entries, err := worktree.ParseList(ctx.RepoRoot)
	if err != nil {
		fail(err.Error())
	}

	var entry *worktree.Entry
	var newName string

	if len(args) == 1 {
		// wut mv <new-name>: rename the current worktree's branch
		newName = args[0]
		cwd, err := os.Getwd()
		if err != nil {
			fail(err.Error())
		}
		entry = findWorktreeForCwd(entries, cwd)
		if entry == nil {
			fail("Not inside a managed worktree. Use: wut mv <old-name> <new-name>")
		}
	} else {
		// wut mv <old-name> <new-name>
		oldName := args[0]
		newName = args[1]
		entry = worktree.FindByBranch(entries, oldName)
		if entry == nil {
			fail(fmt.Sprintf("No worktree found for branch '%s'.", oldName))
		}
	}

	if !entry.Managed {
		fail("Cannot rename the main worktree.")
	}

	if entry.BranchName == "" {
		fail("Worktree has a detached HEAD; cannot rename.")
	}

	oldName := entry.BranchName
	if oldName == newName {
		fail("New name is the same as the current name.")
	}

	// Check that the target branch name doesn't already exist
	if git.RefExists(ctx.RepoRoot, "refs/heads/"+newName) {
		fail(fmt.Sprintf("Branch '%s' already exists.", newName))
	}

	// Rename the branch
	if _, err := git.Run([]string{"branch", "-m", oldName, newName}, ctx.RepoRoot); err != nil {
		fail(fmt.Sprintf("Failed to rename branch: %s", err.Error()))
	}

	// Move the worktree directory to match the new branch name
	worktreesDir := git.GetWorktreesDir(ctx.RepoRoot)
	newRelativePath := worktree.BranchToRelativePath(newName)
	newPath := filepath.Join(worktreesDir, newRelativePath)
	newPath = worktree.UniquePath(newPath)

	if _, err := git.Run([]string{"worktree", "move", entry.Path, newPath}, ctx.RepoRoot); err != nil {
		// Rollback the branch rename
		git.Run([]string{"branch", "-m", newName, oldName}, ctx.RepoRoot)
		fail(fmt.Sprintf("Failed to move worktree: %s", err.Error()))
	}

	// If the user is inside the renamed worktree, update their shell cwd
	cwd, _ := os.Getwd()
	oldAbs, _ := filepath.Abs(entry.Path)
	if strings.HasPrefix(cwd, oldAbs) {
		// Preserve any subdirectory the user was in
		rel, _ := filepath.Rel(oldAbs, cwd)
		fmt.Printf("__WUT_CD__:%s\n", filepath.Join(newPath, rel))
	}
}

// findWorktreeForCwd finds the worktree entry that contains the given cwd.
func findWorktreeForCwd(entries []worktree.Entry, cwd string) *worktree.Entry {
	for i := range entries {
		abs, _ := filepath.Abs(entries[i].Path)
		if cwd == abs || strings.HasPrefix(cwd, abs+string(filepath.Separator)) {
			return &entries[i]
		}
	}
	return nil
}
