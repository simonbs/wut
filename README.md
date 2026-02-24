<div align="center">
  <img width="700" src="/screenshot.png" />
  <h3><strong>wut?</strong> — Worktrees Unexpectedly Tolerable</h3>
  <p>Ephemeral worktrees without the ceremony.</p>
</div>

<hr />

<div align="center">
  <pre>
                                ▄▄▄▄▄   
                       ██      █▀▀▀▀██  
██      ██ ██    ██  ███████       ▄█▀  
▀█  ██  █▀ ██    ██    ██        ▄██▀   
 ██▄██▄██  ██    ██    ██        ██     
 ▀██  ██▀  ██▄▄▄███    ██▄▄▄     ▄▄     
  ▀▀  ▀▀    ▀▀▀▀ ▀▀     ▀▀▀▀     ▀▀     
  </pre>
</div>

<div align="center">
  <a href="#-why">✨ Why?</a>&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="#-getting-started">🚀 Getting Started</a>&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="#-usage">🧭 Usage</a>
</div>

<hr />

## ✨ Why?
Git’s native worktree commands feel tedious and geared toward long-lived worktrees, but I just spin them up for short-lived sessions. **wut?** streamlines that.

**wut?** keeps worktrees in `~/.wut/repos/` and exposes commands like `wut new`, `wut go`, `wut list`, and `wut rm` to manage them.

It still builds directly on Git’s worktrees, so it plays nicely with any other Git CLI or UI. Very opinionated and very much designed for the ✨agentic era✨, unlike the built-in commands that are super tedious.

## 🚀 Getting Started

Install **wut?** using Homebrew as shown below.

```sh
brew tap simonbs/wut https://github.com/simonbs/wut.git
brew install wut
```

You'll need Git on your machine. After installation, add shell integration to your `~/.zshrc` or `~/.bashrc`:

```sh
eval "$(wut init)"
```

This enables automatic directory changing when you run `wut new` or `wut go`. Without it, these commands will prompt you to set up shell integration.

## 🧭 Usage
Run `wut` from inside the repo you want worktrees for.

```sh
$ wut new feature-login
# Creates worktree and switches to it

$ wut list
👉 feature-login  ~/.wut/repos/myapp/feature-login
🏠 main           ~/projects/myapp

$ wut go
# Switches to main worktree

$ wut rm feature-login
# Removes worktree and deletes branch
```

By default, `wut` stores managed worktrees outside the repo. You can override where worktrees are created by setting:

- `WUT_WORKTREES_DIR` (exact directory for this repo/session)
- `git config wut.worktreesDir <path>` (per-repo persistent setting)
- `WUT_WORKTREES_BASE_DIR` or `git config --global wut.worktreesBaseDir <path>` (base directory for all repos)
- `WUT_WORKTREES_INCLUDE_REPO_HASH=true` or `git config wut.includeRepoHash true` (append `<hash>` as `<repo>-<hash>` to reduce collisions for repos with the same name, like `api`)

If a repo already has a legacy `.worktrees/` directory, `wut` will keep using it for backwards compatibility.

```sh
wut new <branch> [--from ref] # Create a new worktree
wut list                      # List worktrees
wut go [branch]               # Navigate to a worktree
wut path <branch>             # Print worktree path
wut rm <branch> [--force]     # Remove a worktree
```

#### Autocompletion

Autocompletion is available for supported shells once you run `eval "$(wut init)"`, so you can tab-complete commands, branch names, and worktree names.

```sh
$ wut go feat<TAB>
# Completes to a matching worktree name
```

#### fzf

`wut` works great paired with [fzf](https://github.com/junegunn/fzf) for fuzzy matching worktrees and `wut list` can be paired with `fzf` as shown below (Thanks [@carlosefonesca](https://masto.pt/@carlosefonseca/115980682267766544)!)

```sh
alias wt='cd "$(wut list | fzf --height=~20 | awk "{print \$3}" | sed "s|^~|$HOME|")"'
```
