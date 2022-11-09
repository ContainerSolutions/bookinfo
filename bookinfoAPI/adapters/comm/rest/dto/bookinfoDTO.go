package dto

// BookInfoResponseDTO represents the struct that is returned by rest endpoints
type BookInfoResponseDTO struct {

	// Unique ID of the bookInfo
	ID string `json:"uuid"`
	// Name of the bookInfo
	Name string `json:"name"`
	// Author of the book
	Author string `json:"author"`
}
