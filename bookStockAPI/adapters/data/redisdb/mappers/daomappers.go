package mappers

import (
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/adapters/data/redisdb/dao"
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/domain"
	"github.com/google/uuid"
)

// MapBookStockDAO2BookStock maps dao BookStock to domain BookStock
func MapBookStockDAO2BookStock(pd dao.BookStockDAO) domain.BookStock {
	return domain.BookStock{
		ID:           pd.ID,
		CurrentStock: pd.CurrentStock,
	}
}

// MapBookStock2BookStockDAO maps domain BookStock to dao BookStock
func MapBookStock2BookStockDAO(p domain.BookStock) dao.BookStockDAO {
	id := p.ID
	if id == "" {
		id = uuid.New().String()
	}
	return dao.BookStockDAO{
		ID:           id,
		CurrentStock: p.CurrentStock,
	}
}
