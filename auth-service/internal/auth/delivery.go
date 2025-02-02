package auth

import "net/http"

// Auth handlers interface
type Handlers interface {
	Register() http.HandlerFunc
	Login() http.HandlerFunc
	TokenValidation() http.HandlerFunc
}
