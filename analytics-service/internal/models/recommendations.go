package models

import "time"

// Action model
type Action struct {
	ID       int64     `json:"id,omitempty"`
	Action   string    `json:"action"`
	ObjectID string    `json:"object_id"`
	Time     time.Time `json:"time"`
}
