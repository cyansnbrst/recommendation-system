package recommendations

import "net/http"

// Recommendations handlers interface
type Handlers interface {
	GetInfo() http.HandlerFunc
}
