package commands

import "fmt"

// ShowHelp displays the CLI help text
func ShowHelp() {
	help := `gt - GitHub Tasks

USAGE:
  gt [-v]                       List all open issues
  gt <number>                   View issue details
  gt <number> -e <field> [text] Edit issue (field: body or title)
  gt <title>                    Create P2 (normal) issue (default)
  gt p0/p1/p2/p3 [-v]           Filter by priority
  gt active [-v]                Show only active tasks
  gt start <number>             Mark issue as active
  gt pause <number>             Remove active label (alias: stop)
  gt done <number>              Close issue
  gt rm <number>                Delete issue (permanent)
  gt setup                      Create required labels in repo

  gt0 <title> [--body [text]]   Create P0 (critical) issue
  gt1 <title> [--body [text]]   Create P1 (important) issue
  gt2 <title> [--body [text]]   Create P2 (normal) issue
  gt3 <title> [--body [text]]   Create P3 (low) issue

FLAGS:
  -v, --verbose                 Show priority labels in output
  -b, --body [text]             Add issue body (inline, editor, or piped)
  -e, --edit <field> [text]     Edit issue field (inline, editor, or piped)

EXAMPLES:
  gt setup                              # Setup labels for this repo
  gt                                    # List all tasks (colors only)
  gt -v                                 # List all tasks (with priority text)
  gt 123                                # View issue #123 (colored title + body)

  # Creating issues
  gt refactor legacy code               # Create default P2 task (no body)
  gt1 fix authentication bug            # Create P1 task (no body)
  gt1 fix auth bug --body               # Create P1 task, open $EDITOR for body
  gt1 "task title" --body "body text"   # Create P1 task with inline body
  cat template.md | gt1 "title" --body  # Create P1 task with piped body

  # Editing issues
  gt 123 -e body                        # Edit issue #123 body in $EDITOR
  gt 123 -e title "New title"           # Update title with inline text
  gt 123 -e body "Updated description"  # Update body with inline text
  echo "New body" | gt 123 -e body      # Update body from stdin (Claude Code!)

  # Workflow
  gt p1 -v                              # List P1 tasks with priority labels
  gt start 234                          # Mark #234 as active
  gt pause 234                          # Remove active (keep open)
  gt done 567                           # Close #567
  gt rm 890                             # Delete #890 (permanent)

WORKFLOW:
  gt2 <title>   - Creates a P2 issue
  gt p2         - Lists existing P2 issues

SETUP:
  1. Run: gh auth login
  2. Navigate to a git repo with GitHub remote
  3. Run: gt setup (creates all required labels)
`
	fmt.Print(help)
}
