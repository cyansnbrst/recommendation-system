package models

// Edit profile DTO struct
type EditProfileDTO struct {
	Location  *string  `json:"location"`
	Interests []string `json:"interests"`
}

// Create profile DTO struct
type CreateProfileDTO struct {
	Name string `json:"name"`
}

// Profile response
type ProfileResponse struct {
	Profile Profile `json:"profile"`
}

// Error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Success response
type SuccessResponse struct {
	Message string `json:"message"`
}
