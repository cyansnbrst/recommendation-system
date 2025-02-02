package models

// Register user DTO struct
type RegisterUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Login user DTO struct
type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create profile DTO struct
type CreateProfileDTO struct {
	Name string `json:"name"`
}

// Error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// User's data response
type UserResponse struct {
	UserUID string `json:"user_uid"`
	IsAdmin bool   `json:"is_admin"`
}
