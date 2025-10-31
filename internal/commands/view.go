package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/DeprecatedLuar/ghtask/internal"
	"github.com/DeprecatedLuar/ghtask/internal/github"
)

func ViewIssue(args []string) {
	issueNum, err := ParseIssueNumber(args, "view")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	repo, err := github.GetRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	cmd := exec.Command("gh", "issue", "view", issueNum,
		"--repo", repo,
		"--json", "number,title,body,labels")

	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error viewing issue: %v\n", err)
		os.Exit(1)
	}

	var viewData struct {
		Number int              `json:"number"`
		Title  string           `json:"title"`
		Body   string           `json:"body"`
		Labels []internal.Label `json:"labels"`
	}

	if err := json.Unmarshal(output, &viewData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing issue: %v\n", err)
		os.Exit(1)
	}

	issue := internal.Issue{
		Number: viewData.Number,
		Title:  viewData.Title,
		Labels: viewData.Labels,
	}

	priority := internal.ExtractPriority(issue)
	color := internal.GetPriorityColor(priority)
	reset := "\033[0m"

	fmt.Printf("%s#%d - %s%s\n\n", color, issue.Number, issue.Title, reset)
	if viewData.Body != "" {
		fmt.Println(viewData.Body)
	}
}
