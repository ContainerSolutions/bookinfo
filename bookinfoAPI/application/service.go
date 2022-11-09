package application

import (
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/domain"
)

// BookInfoRepository is the interface that we expect to be fulfilled to be used as a backend for BookInfo Service
type BookInfoRepository interface {
	List() ([]domain.BookInfo, error)
	Get(string) (domain.BookInfo, error)
}

// BookInfoService represents the struct which contains a BookInfoRepository and exports methods to access the data
type BookInfoService struct {
	bookInfoRepo BookInfoRepository
}

// NewBookInfoService creates a new BookInfoService instance and sets its repository
func NewBookInfoService(pr BookInfoRepository) BookInfoService {
	if pr == nil {
		panic("missing productRepository")
	}
	return BookInfoService{
		bookInfoRepo: pr,
	}
}

// List loads all the data from the included repository and returns them
// Returns an error if the repository returns one
func (ps BookInfoService) List() ([]domain.BookInfo, error) {
	bookInfos, err := ps.bookInfoRepo.List()
	return bookInfos, err
}

// Get selects the bookInfo from the included repository with the given unique identifier, and returns it
// Returns an error if the repository returns one
func (ps BookInfoService) Get(id string) (domain.BookInfo, error) {
	bookInfo, err := ps.bookInfoRepo.Get(id)
	return bookInfo, err
}
