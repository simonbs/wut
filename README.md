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
  <h3><strong>wut</strong> — Worktrees Unexpectedly Tolerable</h3>
  <p>Ephemeral worktrees that stay out of your vibe zone.</p>
</div>

<hr />

<div align="center">
  <a href="#-why">✨ Why</a>&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="#-getting-started">🚀 Getting Started</a>&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="#-usage">🧭 Usage</a>&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="#-configuration">⚙️ Configuration</a>&nbsp;&nbsp;&nbsp;&nbsp;
</div>

<hr />

## ✨ Why
If you love Git worktrees but hate the mess they leave behind, **wut** is for you. Worktrees are amazing for parallel tasks, but the default workflow tends to scatter folders in places you actually care about. wut moves all of that noise into a single hidden home, so your repo stays clean and your brain stays calmer.

It also keeps the workflow simple. You shouldn't have to remember where you put a temporary worktree last week or manually prune folders that Git no longer tracks. wut's job is to make worktrees feel lightweight again: create a branch, jump into it, move on.

## 🚀 Getting Started

Install wut using Homebrew as shown below.

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
Run wut from inside the repo you want worktrees for. wut uses your current repo to decide where to create and manage worktrees, and it won't run from outside to avoid surprises.

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
wut gc [--dry-run]            # Remove orphaned worktrees
```

## ⚙️ Configuration
By default, wut stores worktrees under `.worktrees` in your repo root.

The configuration file lives at `~/.wut/config.json`:

```json
{
  "autoGc": {
    "enabled": true,
    "intervalHours": 6
  }
}
```

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `autoGc.enabled` | bool | `true` | Enable automatic cleanup of orphaned worktrees |
| `autoGc.intervalHours` | int | `6` | Minimum hours between auto-cleanup runs |

You can override the base directory with `WUT_HOME`, which also moves where the config file lives:

```sh
export WUT_HOME="$HOME/.wut-custom"
```

Cleanup is explicit. wut **never** deletes active worktrees on its own. The `wut gc` command only removes orphaned directories that Git no longer knows about, and you can preview what it would remove with `--dry-run`.
