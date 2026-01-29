package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/simonbs/wut/src/config"
	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/worktree"
)

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func maybeAutoGc(ctx *context.Context) {
	enabled, interval := config.GetAutoGcSettings()
	if !enabled || interval == 0 {
		return
	}

	state := config.ReadState()
	if state.LastRunAt != nil && time.Since(*state.LastRunAt) < interval {
		return
	}

	stale, _ := worktree.GetStalePaths(ctx.RepoRoot)
	for _, path := range stale {
		os.RemoveAll(path)
	}

	now := time.Now()
	config.WriteState(config.State{LastRunAt: &now})
}

func tildify(path string) string {
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
