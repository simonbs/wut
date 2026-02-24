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
		oldName := args[0]
		newName = args[1]
		entry = worktree.FindByBranch(entries, oldName)
		if entry == nil {
			fail(fmt.Sprintf("No worktree found for branch '%s'.", oldName))
		}
	}

	entryAbs, _ := filepath.Abs(entry.Path)
	repoRootAbs, _ := filepath.Abs(ctx.RepoRoot)
	if entryAbs == repoRootAbs {
		fail("Cannot rename the main worktree.")
	}

	if entry.BranchName == "" {
		fail("Worktree has a detached HEAD; cannot rename.")
	}

	oldName := entry.BranchName
	if oldName == newName {
		fail("New name is the same as the current name.")
	}

	preMoveCwd, _ := os.Getwd()
	oldAbs, _ := filepath.Abs(entry.Path)
	wasInsideRenamedWorktree := isPathInside(preMoveCwd, oldAbs)

	if git.RefExists(ctx.RepoRoot, "refs/heads/"+newName) {
		fail(fmt.Sprintf("Branch '%s' already exists.", newName))
	}

	if _, err := git.Run([]string{"branch", "-m", oldName, newName}, ctx.RepoRoot); err != nil {
		fail(fmt.Sprintf("Failed to rename branch: %s", err.Error()))
	}

	worktreesDir := git.GetWorktreesDir(ctx.RepoRoot)
	newRelativePath := worktree.BranchToRelativePath(newName)
	newPath := filepath.Join(worktreesDir, newRelativePath)
	newPath = worktree.UniquePath(newPath)

	if _, err := git.Run([]string{"worktree", "move", entry.Path, newPath}, ctx.RepoRoot); err != nil {
		git.Run([]string{"branch", "-m", newName, oldName}, ctx.RepoRoot)
		fail(fmt.Sprintf("Failed to move worktree: %s", err.Error()))
	}

	if wasInsideRenamedWorktree {
		rel, _ := filepath.Rel(oldAbs, preMoveCwd)
		fmt.Printf("__WUT_CD__:%s\n", filepath.Join(newPath, rel))
	}
}

func findWorktreeForCwd(entries []worktree.Entry, cwd string) *worktree.Entry {
	cwdAbs, _ := filepath.Abs(cwd)
	var best *worktree.Entry
	bestLen := -1

	for i := range entries {
		abs, _ := filepath.Abs(entries[i].Path)
		if isPathInside(cwdAbs, abs) {
			if len(abs) > bestLen {
				best = &entries[i]
				bestLen = len(abs)
			}
		}
	}
	return best
}

func isPathInside(path, parent string) bool {
	pathAbs, _ := filepath.Abs(path)
	parentAbs, _ := filepath.Abs(parent)

	if pathAbs == parentAbs || strings.HasPrefix(pathAbs, parentAbs+string(filepath.Separator)) {
		return true
	}

	pathEval, errPath := filepath.EvalSymlinks(pathAbs)
	parentEval, errParent := filepath.EvalSymlinks(parentAbs)
	if errPath == nil && errParent == nil {
		return pathEval == parentEval || strings.HasPrefix(pathEval, parentEval+string(filepath.Separator))
	}

	return false
}
