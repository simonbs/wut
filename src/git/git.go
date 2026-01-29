package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run(args []string, cwd string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = cwd
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("%s", errMsg)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func RefExists(repoRoot, ref string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", ref)
	cmd.Dir = repoRoot
	return cmd.Run() == nil
}

func GetRepoRoot(cwd string) (string, error) {
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	// Check if .git is a file (worktree)
	gitPath := filepath.Join(cwd, ".git")
	info, err := os.Stat(gitPath)
	if err == nil && !info.IsDir() {
		content, err := os.ReadFile(gitPath)
		if err == nil {
			line := strings.TrimSpace(string(content))
			if strings.HasPrefix(line, "gitdir: ") {
				gitdir := strings.TrimPrefix(line, "gitdir: ")
				// Extract main repo path from worktree gitdir
				if idx := strings.Index(gitdir, "/.git/worktrees"); idx != -1 {
					return gitdir[:idx], nil
				}
			}
		}
	}

	// Normal repo - use rev-parse
	output, err := Run([]string{"rev-parse", "--show-toplevel"}, cwd)
	if err != nil {
		return "", fmt.Errorf("Not inside a Git repository.")
	}
	return output, nil
}

func GetWorktreesDir(repoRoot string) string {
	return filepath.Join(repoRoot, ".worktrees")
}
