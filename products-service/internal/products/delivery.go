package products

import "net/http"

// Products handlers interface
type Handlers interface {
	Get() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	Create() http.HandlerFunc
}
