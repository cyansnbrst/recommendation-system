package analytics

import "time"

// Recommendations repository interface
type Repository interface {
	Insert(action string, objectID string, actionTime time.Time) error
}
