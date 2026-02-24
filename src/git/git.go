package git

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const legacyWorktreesDirName = ".worktrees"

var repoNameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

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

	// Ask Git for the common .git directory so this also works when called
	// from nested paths inside linked worktrees.
	commonDir, err := Run([]string{"rev-parse", "--path-format=absolute", "--git-common-dir"}, cwd)
	if err == nil {
		commonDir = strings.TrimSpace(commonDir)
		if commonDir != "" && filepath.Base(commonDir) == ".git" {
			return filepath.Dir(commonDir), nil
		}
	}

	// Fallback for older Git versions or unusual environments.
	output, err := Run([]string{"rev-parse", "--show-toplevel"}, cwd)
	if err != nil {
		return "", fmt.Errorf("Not inside a Git repository.")
	}
	return output, nil
}

func GetWorktreesDir(repoRoot string) string {
	if dir, err := ResolveWorktreesDir(repoRoot); err == nil {
		return dir
	}
	return LegacyWorktreesDir(repoRoot)
}

func LegacyWorktreesDir(repoRoot string) string {
	return filepath.Join(repoRoot, legacyWorktreesDirName)
}

func ResolveWorktreesDir(repoRoot string) (string, error) {
	if value := strings.TrimSpace(os.Getenv("WUT_WORKTREES_DIR")); value != "" {
		return resolvePath(value, repoRoot)
	}

	if value, ok := getRepoConfigValue(repoRoot, "wut.worktreesDir"); ok {
		return resolvePath(value, repoRoot)
	}

	legacyDir := LegacyWorktreesDir(repoRoot)
	if info, err := os.Stat(legacyDir); err == nil && info.IsDir() {
		return legacyDir, nil
	}

	baseDir, err := resolveDefaultBaseDir(repoRoot)
	if err != nil {
		return "", err
	}

	repoAbs, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", err
	}

	return filepath.Join(baseDir, buildRepoKey(repoAbs)), nil
}

func getRepoConfigValue(repoRoot, key string) (string, bool) {
	value, err := Run([]string{"config", "--get", key}, repoRoot)
	if err != nil || strings.TrimSpace(value) == "" {
		return "", false
	}
	return strings.TrimSpace(value), true
}

func resolveDefaultBaseDir(repoRoot string) (string, error) {
	if value := strings.TrimSpace(os.Getenv("WUT_WORKTREES_BASE_DIR")); value != "" {
		return resolvePath(value, repoRoot)
	}

	if value, ok := getGlobalConfigValue("wut.worktreesBaseDir"); ok {
		return resolvePath(value, repoRoot)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".wut", "repos"), nil
}

func getGlobalConfigValue(key string) (string, bool) {
	value, err := Run([]string{"config", "--global", "--get", key}, ".")
	if err != nil || strings.TrimSpace(value) == "" {
		return "", false
	}
	return strings.TrimSpace(value), true
}

func getBoolValue(value string) (bool, bool) {
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return false, false
	}
	return parsed, true
}

func getRepoConfigBool(repoRoot, key string) (bool, bool) {
	value, ok := getRepoConfigValue(repoRoot, key)
	if !ok {
		return false, false
	}
	return getBoolValue(value)
}

func getGlobalConfigBool(key string) (bool, bool) {
	value, ok := getGlobalConfigValue(key)
	if !ok {
		return false, false
	}
	return getBoolValue(value)
}

func getEnvBool(key string) (bool, bool) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return false, false
	}
	return getBoolValue(value)
}

func resolvePath(pathValue, repoRoot string) (string, error) {
	trimmed := strings.TrimSpace(pathValue)
	if trimmed == "" {
		return "", fmt.Errorf("worktree path is empty")
	}

	expanded := os.ExpandEnv(trimmed)
	if strings.HasPrefix(expanded, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		suffix := strings.TrimPrefix(expanded, "~")
		suffix = strings.TrimPrefix(suffix, string(os.PathSeparator))
		expanded = filepath.Join(homeDir, suffix)
	}

	if !filepath.IsAbs(expanded) {
		expanded = filepath.Join(repoRoot, expanded)
	}

	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func buildRepoKey(repoRoot string) string {
	repoName := filepath.Base(repoRoot)
	repoName = strings.TrimSpace(repoName)
	if repoName == "" || repoName == "." || repoName == string(filepath.Separator) {
		repoName = "repo"
	}
	repoName = repoNameSanitizer.ReplaceAllString(repoName, "-")

	if !shouldIncludeRepoHash(repoRoot) {
		return repoName
	}

	sum := sha1.Sum([]byte(repoRoot))
	hash := hex.EncodeToString(sum[:])
	return fmt.Sprintf("%s-%s", repoName, hash[:10])
}

func shouldIncludeRepoHash(repoRoot string) bool {
	if value, ok := getEnvBool("WUT_WORKTREES_INCLUDE_REPO_HASH"); ok {
		return value
	}

	if value, ok := getRepoConfigBool(repoRoot, "wut.includeRepoHash"); ok {
		return value
	}

	if value, ok := getGlobalConfigBool("wut.includeRepoHash"); ok {
		return value
	}

	return false
}
