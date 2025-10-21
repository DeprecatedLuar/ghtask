// Package commands implements all gt command operations including issue creation,
// listing, filtering, lifecycle management, and repository setup. All operations
// delegate to the gh CLI and use GitHub Issues as the single source of truth.
package commands

type Label struct {
	Name string `json:"name"`
}

type Issue struct {
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	Labels    []Label `json:"labels"`
	CreatedAt string  `json:"createdAt"`
}
