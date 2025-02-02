package models

// Recommendations model
type Recommendation struct {
	ID        int64  `json:"id,omitempty"`
	UserUID   string `json:"user_uid,omitempty"`
	ProductID int64  `json:"product_id"`
}

// User interests model
type User struct {
	UserUID   string   `json:"user_uid"`
	Interests []string `json:"interests"`
}

// Product rating model
type Product struct {
	ProductID  int64    `json:"product_id"`
	Tags       []string `json:"tags"`
	Popularity int64    `json:"popularity"`
}
