package application

import (
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/domain"
)

// BookStockRepository is the interface that we expect to be fulfilled to be used as a backend for BookStock Service
type BookStockRepository interface {
	Get(string) (domain.BookStock, error)
}

// BookStockService represents the struct which contains a BookStockRepository and exports methods to access the data
type BookStockService struct {
	bookInfoRepo BookStockRepository
}

// NewBookStockService creates a new BookStockService instance and sets its repository
func NewBookStockService(pr BookStockRepository) BookStockService {
	if pr == nil {
		panic("missing productRepository")
	}
	return BookStockService{
		bookInfoRepo: pr,
	}
}

// Get selects the bookInfo from the included repository with the given unique identifier, and returns it
// Returns an error if the repository returns one
func (ps BookStockService) Get(id string) (domain.BookStock, error) {
	bookInfo, err := ps.bookInfoRepo.Get(id)
	return bookInfo, err
}
