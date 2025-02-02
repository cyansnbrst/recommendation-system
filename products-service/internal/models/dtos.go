package models

// Create product DTO struct
type CreateProductDTO struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// Update product DTO struct
type UpdateProductDTO struct {
	Name *string  `json:"name"`
	Tags []string `json:"tags"`
}

// Product response
type ProductResponse struct {
	Product Product `json:"product"`
}

// Error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Success response
type SuccessResponse struct {
	Message string `json:"message"`
}
