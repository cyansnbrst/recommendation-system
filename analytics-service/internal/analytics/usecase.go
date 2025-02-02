package analytics

import "time"

// Analytics use case interface
type UseCase interface {
	Insert(action string, objectID string, actionTime time.Time) error
}
