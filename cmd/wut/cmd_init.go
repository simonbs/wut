package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/worktree"
)

func detectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "fish") {
		return "fish"
	}
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	return "bash"
}

func cmdInit(args []string) {
	shell := ""
	for _, arg := range args {
		switch arg {
		case "--fish":
			shell = "fish"
		case "--bash":
			shell = "bash"
		case "--zsh":
			shell = "zsh"
		}
	}
	if shell == "" {
		shell = detectShell()
	}

	binPath, _ := os.Executable()

	switch shell {
	case "fish":
		printFishWrapper(binPath)
	default:
		printBashZshWrapper(binPath)
	}
}

func printBashZshWrapper(binPath string) {
	wrapper := `wut() {
  local wut_bin="` + binPath + `"
  export WUT_WRAPPER_ACTIVE=1
  local output
  output=$("$wut_bin" "$@" 2>&1)
  local exit_code=$?
  local cd_marker
  cd_marker=$(echo "$output" | grep "^__WUT_CD__:" | head -1)
  if [ -n "$cd_marker" ]; then
    local target_dir="${cd_marker#__WUT_CD__:}"
    if [ -d "$target_dir" ]; then
      cd "$target_dir" || return 1
    fi
    local filtered
    filtered=$(printf "%s" "$output" | grep -v "^__WUT_CD__:")
    if [[ -n "${filtered//[[:space:]]/}" ]]; then
      printf "%s\\n" "$filtered"
    fi
  else
    if [[ -n "${output//[[:space:]]/}" ]]; then
      printf "%s\\n" "$output"
    fi
  fi
  return $exit_code
}

_wut_completions() {
  local cur="${COMP_WORDS[COMP_CWORD]}"
  local prev="${COMP_WORDS[COMP_CWORD-1]}"
  
  if [[ ${COMP_CWORD} -eq 1 ]]; then
    COMPREPLY=($(compgen -W "new list go path rm" -- "$cur"))
    return
  fi
  
  case "$prev" in
    go|path|rm)
      local branches
      branches=$("` + binPath + `" --completions branches 2>/dev/null)
      COMPREPLY=($(compgen -W "$branches" -- "$cur"))
      ;;
  esac
}

complete -F _wut_completions wut

# zsh completion
if [[ -n ${ZSH_VERSION-} ]]; then
  autoload -U +X bashcompinit && bashcompinit
fi`
	fmt.Println(wrapper)
}

func printFishWrapper(binPath string) {
	wrapper := `function wut
    set -l wut_bin "` + binPath + `"
    set -x WUT_WRAPPER_ACTIVE 1
    set -l output ($wut_bin $argv 2>&1)
    set -l exit_code $status
    set -l cd_marker (echo "$output" | grep "^__WUT_CD__:" | head -1)
    if test -n "$cd_marker"
        set -l target_dir (string replace "__WUT_CD__:" "" "$cd_marker")
        if test -d "$target_dir"
            cd "$target_dir"
        end
        set -l filtered (printf "%s" "$output" | grep -v "^__WUT_CD__:")
        if test -n "$filtered"
            printf "%s\n" "$filtered"
        end
    else
        if test -n "$output"
            printf "%s\n" "$output"
        end
    end
    return $exit_code
end

# Subcommands
complete -c wut -f -n "__fish_use_subcommand" -a "new" -d "Create new worktree"
complete -c wut -f -n "__fish_use_subcommand" -a "list" -d "List worktrees"
complete -c wut -f -n "__fish_use_subcommand" -a "go" -d "Go to worktree"
complete -c wut -f -n "__fish_use_subcommand" -a "path" -d "Print worktree path"
complete -c wut -f -n "__fish_use_subcommand" -a "rm" -d "Remove worktree"

# Branch completions for go/path/rm
complete -c wut -f -n "__fish_seen_subcommand_from go path rm" -a "(` + binPath + ` --completions branches 2>/dev/null)"`
	fmt.Println(wrapper)
}

func cmdCompletions(args []string) {
	if len(args) < 1 {
		return
	}

	switch args[0] {
	case "branches":
		ctx, err := context.Create()
		if err != nil {
			return
		}
		entries, _ := worktree.ParseList(ctx.RepoRoot)
		for _, e := range entries {
			if e.BranchName != "" {
				fmt.Println(e.BranchName)
			}
		}
	}
}
