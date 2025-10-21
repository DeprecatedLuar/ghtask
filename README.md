# gt - GitHub Tasks

Lightweight CLI tool for managing GitHub Issues with a dstask-inspired workflow.

## Features

- **Fast task creation** via priority shortcuts (`g0`, `g1`, `g2`, `g3`)
- **Visual priority feedback** with color-coded output
- **Zero local state** - GitHub Issues is the source of truth
- **Automatic repo detection** - works in any git project with GitHub remote
- **Cross-platform** - Linux and Windows support
- **Collaborative** - both users see all tasks in real-time

## Installation

### Prerequisites

1. Install [GitHub CLI](https://cli.github.com/):
   ```bash
   # Check if already installed
   gh --version
   ```

2. Authenticate with GitHub:
   ```bash
   gh auth login
   ```

### Build & Install

```bash
# Clone the repository
cd ~/Workspace/projects/cli/ghtasks/

# Build the binary
go build -o gt main.go

# Install to PATH
cp gt ~/Workspace/tools/bin/

# Create symlinks for priority shortcuts
cd ~/Workspace/tools/bin
ln -sf gt g0
ln -sf gt g1
ln -sf gt g2
ln -sf gt g3
```

### Repository Setup

For each project you want to use `gt` with:

```bash
# Navigate to your project
cd ~/projects/your-project/

# Create required labels
gh label create "inbox" --color "d4c5f9"
gh label create "active" --color "0e8a16"
gh label create "P0" --color "d93f0b"
gh label create "P1" --color "ff9800"
gh label create "P2" --color "ffeb3b"
gh label create "P3" --color "cccccc"
```

## Usage

### Quick Task Creation

```bash
g0 production is down!              # Critical (P0)
g1 implement user authentication    # Important (P1)
g2 add unit tests for login         # Normal (P2)
g3 update documentation             # Low priority (P3)
```

### List Tasks

```bash
gt              # List all open issues
gt p0           # Filter by P0 priority
gt p1           # Filter by P1 priority
gt p2           # Filter by P2 priority
gt p3           # Filter by P3 priority
gt active       # Show only active tasks
```

### Manage Tasks

```bash
gt start 234    # Mark issue #234 as active
gt done 567     # Close issue #567
```

## How It Works

### Automatic Repository Detection

`gt` automatically detects your GitHub repository from the current directory's git remote:

```bash
cd ~/projects/project-one/
g1 fix bug           # Posts to project-one's GitHub repo
gt                   # Lists project-one's issues

cd ~/projects/project-two/
g1 add feature       # Posts to project-two's GitHub repo
gt                   # Lists project-two's issues
```

No configuration needed - just `cd` to your project and use `gt`!

### Visual Priority System

Tasks are color-coded for quick scanning:

- **P0 (Critical)**: Red - do this first!
- **P1 (Important)**: Orange
- **P2 (Normal)**: White
- **P3 (Low)**: Gray

Active tasks get a highlighted background to stand out.

### Label Behavior

- Tasks created with `g1 fix bug` get labels: `inbox`, `P1`
- Running `gt start 234` adds `active` label: `inbox`, `P1`, `active`
- Labels coexist - `inbox` remains when `active` is added
- `gt done 234` closes the issue

## Daily Workflow Example

**Morning - Check tasks:**
```bash
$ gt
#234  P0  Fix auth bug               # RED - do this first!
#235  P1  Implement feature X        # ORANGE - important
#567  P2  Add tests                  # WHITE (ACTIVE - highlighted)
#891  P3  Update docs                # GRAY - whenever
```

**Pick a task to work on:**
```bash
$ gt start 235
✓ Activated #235
```

**Quick task capture during the day:**
```bash
$ g1 investigate memory leak
Created issue #892

$ g2 refactor authentication module
Created issue #893
```

**Complete a task:**
```bash
$ gt done 567
✓ Closed #567: Add tests
```

## Commands

| Command | Description |
|---------|-------------|
| `gt` | List all open issues |
| `gt p0/p1/p2/p3` | Filter by priority |
| `gt active` | Show only active tasks |
| `gt start <number>` | Mark issue as active |
| `gt done <number>` | Close issue |
| `g0 <title>` | Create P0 (critical) issue |
| `g1 <title>` | Create P1 (important) issue |
| `g2 <title>` | Create P2 (normal) issue |
| `g3 <title>` | Create P3 (low) issue |

## Troubleshooting

**Not in a git repository:**
```
Error: not in a git repository or no origin remote
```
Solution: Make sure you're in a git repository with a GitHub remote (`git remote -v`)

**Remote is not GitHub:**
```
Error: could not parse GitHub repo from: git@gitlab.com:user/repo.git
```
Solution: `gt` only works with GitHub repositories

**gh not installed:**
```
Error: 'gh' command not found
```
Solution: Install from https://cli.github.com/

**Not authenticated:**
```
Error: Not authenticated with GitHub
```
Solution: Run `gh auth login`

## Environment Variables (Optional)

```bash
# Override git-detected repo (rarely needed)
export GT_REPO="owner/repo-name"

# Use different GitHub account
export GITHUB_TOKEN="ghp_..."
```

## Development

### Building from Source

```bash
go build -o gt main.go
```

### Running Tests

```bash
# Test in a git repository
cd ~/projects/some-github-project/
gt help
gt
```

## Philosophy

- **Simple**: Fast task capture, no complex workflows
- **Zero local state**: GitHub Issues is always the source of truth
- **Visual**: Color-coded priorities for quick scanning
- **Collaborative**: Both users see all tasks in real-time
- **Context-aware**: Automatically uses the right repo based on current directory

## License

MIT
