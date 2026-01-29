package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/simonbs/wut/src/git"
)

type Entry struct {
	Path       string
	Head       string
	BranchRef  string
	BranchName string
	Detached   bool
	Managed    bool
}

func ParseList(repoRoot string) ([]Entry, error) {
	output, err := git.Run([]string{"worktree", "list", "--porcelain"}, repoRoot)
	if err != nil {
		return nil, err
	}

	worktreesDir := git.GetWorktreesDir(repoRoot)
	var entries []Entry
	var current Entry

	for _, line := range strings.Split(output, "\n") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			current.Path = strings.TrimPrefix(line, "worktree ")
			current.Managed = strings.HasPrefix(current.Path, worktreesDir)
		case strings.HasPrefix(line, "HEAD "):
			current.Head = strings.TrimPrefix(line, "HEAD ")
		case strings.HasPrefix(line, "branch "):
			ref := strings.TrimPrefix(line, "branch ")
			current.BranchRef = ref
			current.BranchName = strings.TrimPrefix(ref, "refs/heads/")
			current.Detached = false
		case strings.HasPrefix(line, "detached"):
			current.Detached = true
		case line == "":
			if current.Path != "" {
				entries = append(entries, current)
			}
			current = Entry{}
		}
	}

	if current.Path != "" {
		entries = append(entries, current)
	}

	return entries, nil
}

func BranchToRelativePath(branch string) string {
	return strings.ReplaceAll(branch, "/", "-")
}

func UniquePath(basePath string) string {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return basePath
	}

	counter := 1
	for {
		candidate := fmt.Sprintf("%s-%d", basePath, counter)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
		counter++
	}
}

func FindByBranch(entries []Entry, branch string) *Entry {
	for i := range entries {
		if entries[i].BranchName == branch {
			return &entries[i]
		}
	}
	return nil
}

func FindByPath(entries []Entry, absPath string) *Entry {
	for i := range entries {
		entryAbs, _ := filepath.Abs(entries[i].Path)
		if entryAbs == absPath {
			return &entries[i]
		}
	}
	return nil
}

func GetStalePaths(repoRoot string) ([]string, error) {
	entries, err := ParseList(repoRoot)
	if err != nil {
		return nil, err
	}

	managedPaths := make(map[string]bool)
	for _, e := range entries {
		if e.Managed {
			abs, _ := filepath.Abs(e.Path)
			managedPaths[abs] = true
		}
	}

	worktreesDir := git.GetWorktreesDir(repoRoot)
	if _, err := os.Stat(worktreesDir); os.IsNotExist(err) {
		return nil, nil
	}

	var stale []string
	err = filepath.WalkDir(worktreesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == ".git" {
			return filepath.SkipDir
		}

		gitFile := filepath.Join(path, ".git")
		info, err := os.Stat(gitFile)
		if err == nil && !info.IsDir() {
			abs, _ := filepath.Abs(path)
			if !managedPaths[abs] {
				stale = append(stale, path)
			}
			return filepath.SkipDir
		}
		return nil
	})

	return stale, err
}
