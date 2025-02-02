package models

// User's profile model
type Profile struct {
	UserUID   string   `json:"user_uid"`
	Name      string   `json:"name"`
	Location  string   `json:"location"`
	Interests []string `json:"interests,omitempty"`
}
