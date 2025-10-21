<h1 align="center">GhTask</h1>

<p align="center">Multiplayer task manager on github issues for co-developed repos</p>

<p align="center">
  <a href="https://github.com/DeprecatedLuar/ghtask/stargazers">
    <img src="https://img.shields.io/github/stars/DeprecatedLuar/ghtask?style=for-the-badge&logo=github&color=1f6feb&logoColor=white&labelColor=black"/>
  </a>
  <a href="https://github.com/DeprecatedLuar/ghtask/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/DeprecatedLuar/ghtask?style=for-the-badge&color=green&labelColor=black"/>
  </a>
</p>

---

## What makes it cool

- **Instant task capture** - `g0`, `g1`, `g2`, `g3` shortcuts for priority-based task creation
- **Multiplayer by default** - Both devs see all tasks in real-time, no sync needed
- **Visual task tracking** - Color-coded priorities (Red = urgent, Gray = later)
- **Context-aware** - Auto-detects repo, just `cd` and go
- **Zero local state** - GitHub Issues is the truth, no sync conflicts

---

## Installation

**Prerequisites:** [GitHub CLI](https://cli.github.com/) + `gh auth login`

### Linux

```bash
# Download latest binary
wget https://github.com/DeprecatedLuar/ghtask/releases/latest/download/gt-linux-amd64

# Make executable
chmod +x gt-linux-amd64

# Move to PATH
mv gt-linux-amd64 ~/Workspace/tools/bin/gt  # Or wherever you keep bins

# Setup labels in your project
cd ~/your-project/
gt setup
```

### Windows

```powershell
# Download latest binary (PowerShell)
Invoke-WebRequest -Uri "https://github.com/DeprecatedLuar/ghtask/releases/latest/download/gt-windows-amd64.exe" -OutFile "gt.exe"

# Move to a directory in your PATH (adjust as needed)
Move-Item gt.exe C:\Users\YourName\bin\gt.exe

# Setup labels in your project
cd C:\your-project
gt setup
```

<details>
<summary>Build from source</summary>

<br>

```bash
git clone https://github.com/DeprecatedLuar/ghtask.git
cd ghtask
go build -o gt
cp gt ~/Workspace/tools/bin/  # Linux/Mac

# Windows
go build -o gt.exe
move gt.exe C:\Users\YourName\bin\
```

</details>

---

## Commands

| Command | Description |
|---------|-------------|
| `gt` | List all open issues |
| `gt -v` | List all issues with priority labels (verbose) |
| `gt p0/p1/p2/p3` | Filter by priority |
| `gt active` | Show only active tasks |
| `gt start <number>` | Mark issue as active |
| `gt done <number>` | Close issue |
| `gt rm <number>` | Delete issue (permanent) |
| `gt setup` | Create required labels in repo |
| `g0 <title>` | Create P0 (critical) issue |
| `g1 <title>` | Create P1 (important) issue |
| `g2 <title>` | Create P2 (normal) issue |
| `g3 <title>` | Create P3 (low) issue |
| `gt <title>` | Create P2 (normal) issue (default) |

<details>
<summary>Quick Start</summary>

<br>

```bash
# Create tasks with priority shortcuts
g0 production is down!              # Critical (P0)
g1 implement user authentication    # Important (P1)
g2 add unit tests for login         # Normal (P2)
g3 update documentation             # Low priority (P3)

# List and filter
gt              # All open issues
gt p0           # Only P0 tasks
gt active       # Only active tasks

# Manage lifecycle
gt start 234    # Mark #234 as active
gt done 567     # Close #567
```

</details>

<details>
<summary>Workflow example</summary>

<br>

**Morning - Check what's on fire:**
```bash
$ gt
#234  P0  Fix auth bug               # RED - do this NOW
#235  P1  Implement feature X        # ORANGE - important
#567  P2  Add tests                  # WHITE (ACTIVE - highlighted)
#891  P3  Update docs                # GRAY - meh
```

**Start working on something:**
```bash
$ gt start 235
✓ Activated #235
```

**Capture tasks as they come up:**
```bash
$ g1 investigate memory leak
Created issue #892

$ g2 refactor auth module
Created issue #893
```

**Wrap up:**
```bash
$ gt done 567
✓ Closed #567: Add tests
```

</details>

<details>
<summary>Troubleshooting</summary>

<br>

**Not in a git repository:**
```
Error: not in a git repository or no origin remote
```
→ Make sure you're in a git repo with a GitHub remote (`git remote -v`)

**Remote is not GitHub:**
```
Error: could not parse GitHub repo from: git@gitlab.com:user/repo.git
```
→ `gt` only works with GitHub repos (sorry GitLab fans)

**gh not installed:**
```
Error: 'gh' command not found
```
→ Install from https://cli.github.com/

**Not authenticated:**
```
Error: Not authenticated with GitHub
```
→ Run `gh auth login`

**Environment variables (optional):**
```bash
export GT_REPO="owner/repo"        # Override auto-detected repo
export GITHUB_TOKEN="ghp_..."      # Use different GitHub account
```

</details>

---

<p align="center">
  <a href="https://github.com/DeprecatedLuar/ghtask/issues">
    <img src="https://img.shields.io/badge/Found%20a%20bug%3F-Report%20it!-red?style=for-the-badge&logo=github&logoColor=white&labelColor=black"/>
  </a>
</p>
