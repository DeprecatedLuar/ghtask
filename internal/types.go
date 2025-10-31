// Package internal contains shared types and utilities used across gt commands.
package internal

type Label struct {
	Name string `json:"name"`
}

type Issue struct {
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	Labels    []Label `json:"labels"`
	CreatedAt string  `json:"createdAt"`
}
