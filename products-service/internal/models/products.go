package models

// Product model
type Product struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags,omitempty"`
	Version int64    `json:"version"`
}
