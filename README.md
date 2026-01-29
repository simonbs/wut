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
  <h3><strong>wut?</strong> — Worktrees Unexpectedly Tolerable</h3>
  <p>Ephemeral worktrees that stay out of your vibe zone.</p>
</div>

<hr />

<div align="center">
  <a href="#-why">✨ Why</a>&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="#-getting-started">🚀 Getting Started</a>&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="#-usage">🧭 Usage</a>
</div>

<hr />

## ✨ Why
Git’s native worktree commands feel tedious and geared toward long-lived worktrees, but I just spin them up for short-lived sessions. **wut?* streamlines that.

**wut?** keeps worktrees in the `.worktrees` folder (auto-ignored by Git) inside the repo and exposes commands like `wut new`, `wut go`, `wut list`, and `wut rm` to manage them.

It still builds directly on Git’s worktrees, so it plays nicely with any other Git CLI or UI. Very opinionated. Very syntactic sugar. Very much designed for the agentic era, unlike the built-in commands that are super tedious.

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
Run `wut` from inside the repo you want worktrees for. `wut` uses your current repo to decide where to create and manage worktrees, and it won't run from outside to avoid surprises.

Autocompletion is available for supported shells once you run `eval "$(wut init)"`, so you can tab-complete commands, branch names, and worktree names.

Example tab completion:

```sh
$ wut go feat<TAB>
# Completes to a matching worktree name
```

```sh
$ wut new feature-login
# Creates worktree and switches to it

$ wut list
👉 feature-login  ~/projects/myapp/.worktrees/feature-login
🏠 main           ~/projects/myapp

$ wut go
# Switches to main worktree

$ wut rm feature-login
# Removes worktree and deletes branch
```

Here's the full command list:

```sh
wut new <branch> [--from ref] # Create a new worktree
wut list                      # List worktrees
wut go [branch]               # Navigate to a worktree
wut path <branch>             # Print worktree path
wut rm <branch> [--force]     # Remove a worktree
```
