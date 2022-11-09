package mappers

import (
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/data/mongodb/dao"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/domain"
	"github.com/google/uuid"
)

// MapBookInfoDAO2BookInfo maps dao BookInfo to domain BookInfo
func MapBookInfoDAO2BookInfo(pd dao.BookInfoDAO) domain.BookInfo {
	return domain.BookInfo{
		ID:     pd.ID,
		Name:   pd.Name,
		Author: pd.Author,
	}
}

// MapBookInfo2BookInfoDAO maps domain BookInfo to dao BookInfo
func MapBookInfo2BookInfoDAO(p domain.BookInfo) dao.BookInfoDAO {
	id := p.ID
	if id == "" {
		id = uuid.New().String()
	}
	return dao.BookInfoDAO{
		ID:     id,
		Name:   p.Name,
		Author: p.Author,
	}
}
