package dto

// BookStockResponseDTO represents the struct that is returned by rest endpoints
type BookStockResponseDTO struct {

	// Unique ID of the bookInfo
	ID string `json:"uuid"`
	// CurrentStock is the current stock of the book
	CurrentStock int `json:"currentStock"`
}
