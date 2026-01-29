package main

import (
	"fmt"
	"os"

	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/worktree"
)

func cmdInit() {
	binPath, _ := os.Executable()
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
