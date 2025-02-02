package profiles

import "net/http"

// Profiles handlers interface
type Handlers interface {
	GetInfo() http.HandlerFunc
	EditData() http.HandlerFunc
	CreateProfile() http.HandlerFunc
}
